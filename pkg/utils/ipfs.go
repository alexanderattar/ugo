package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/ipfs/go-ipfs-api"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs-cmdkit/files"
)

var s *ipfsapi.Shell
var req *ipfsapi.Request

type IPFSError ipfsapi.Error
type IPFSResponse struct {
	Output io.ReadCloser
	Error  *IPFSError
}

func (e *IPFSError) Error() string {
	var out string
	if e.Command != "" {
		out = e.Command + ": "
	}
	if e.Code != 0 {
		out = fmt.Sprintf("%s%d: ", out, e.Code)
	}
	return out + e.Message
}

func (r *IPFSResponse) Close() error {
	if r.Output != nil {
		ioutil.ReadAll(r.Output)
		return r.Output.Close()
	}
	return nil
}

// CreateIpfsShell creates an Infura shell
func CreateIpfsShell() {
	log.Println("Connecting to Infura IPFS shell")
	s = shell.NewShell("https://ipfs.infura.io:5001")
}

// DagPut is a helper function for other packages to call DagPut to Infura
func DagPut(data []byte) (string, error) {
	hash, err := dagPut(data, "json", "cbor", "false")
	return hash, err
}

// AddCids takes json and recursively adds cids to each object,
// returns the json with cids, and puts & pins to Infura nodes in the background.
func AddCids(data interface{}, wg *sync.WaitGroup, groupErr *error) {
	_, isMap := data.(map[string]interface{})
	_, isArray := data.([]interface{})

	if isMap {
		if data.(map[string]interface{})["cid"] != nil {
			return
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			// Protect against invalid request data
			if data.(map[string]interface{}) == nil {
				return
			}

			json, err := json.Marshal(data)
			if err != nil {
				*groupErr = err
				return
			}

			hash, err := dagPut(json, "json", "cbor", "false")
			if err != nil {
				*groupErr = err
				return
			}

			data.(map[string]interface{})["cid"] = hash
			AddToInfuraDagPinQueue(json, hash, 0)
		}()

		for _, v := range data.(map[string]interface{}) {
			switch v := v.(type) {
			case string:
			case []interface{}:
				AddCids(v, wg, groupErr)
			case interface{}:
				AddCids(v, wg, groupErr)
			}
		}
	} else if isArray {
		for _, v := range data.([]interface{}) {
			AddCids(v, wg, groupErr)
		}
	}
}

func uploadToInfura(data []byte, expectedHash string, attempts int) {
	start := time.Now()
	infuraHash, err := dagPut(data, "json", "cbor", "true")
	t := time.Now()
	elapsed := t.Sub(start)

	if infuraHash == expectedHash {
		lg.Println("Time to put to infura:", elapsed)
		return
	}

	if err != nil {
		lg.Errorln(err)
	} else {
		lg.Errorln("Got this hash from Infura:", infuraHash, " but expected:", expectedHash)
	}

	if attempts < 10 {
		attempts++
		AddToInfuraDagPinQueue(data, expectedHash, attempts)
	} else {
		lg.Errorln("Here's the data we weren't able to pin:", string(data))
	}
}

func newRequest(ctx context.Context, command string, args ...string) *shell.Request {
	return shell.NewRequest(ctx, "https://ipfs.infura.io:5001", command, args...)
}

func dagPut(data interface{}, ienc, kind string, pin string) (string, error) {
	req := newRequest(context.Background(), "dag/put")
	req.Opts = map[string]string{
		"input-enc": ienc,
		"format":    kind,
	}

	if pin == "true" {
		req.Opts["pin"] = "true"
	}

	var r io.Reader
	switch data := data.(type) {
	case string:
		r = strings.NewReader(data)
	case []byte:
		r = bytes.NewReader(data)
	case io.Reader:
		r = data
	default:
		return "", fmt.Errorf("cannot current handle putting values of type %T", data)
	}
	rc := ioutil.NopCloser(r)
	fr := files.NewReaderFile("", "", rc, nil)
	slf := files.NewSliceFile("", "", []files.File{fr})
	fileReader := files.NewMultiFileReader(slf, true)
	req.Body = fileReader
	c := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}

	resp, err := sendWithOrigin(c, req)
	if err != nil {
		return "", err
	}
	defer resp.Close()

	if resp.Error != nil {
		return "", resp.Error
	}

	var out struct {
		Cid struct {
			Target string `json:"/"`
		}
	}
	err = json.NewDecoder(resp.Output).Decode(&out)
	if err != nil {
		return "", err
	}

	return out.Cid.Target, nil
}

// SendWithOrigin monkey-patches the Send method from go-ipfs-api's request.go file.
// The only change is that we set the origin header as https://ujomusic.com which is
// on Infura's allow list.
func sendWithOrigin(c *http.Client, r *shell.Request) (*IPFSResponse, error) {
	url := getURL(r)
	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		return nil, err
	}

	if fr, ok := r.Body.(*files.MultiFileReader); ok {
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+fr.Boundary())
		req.Header.Set("Content-Disposition", "form-data: name=\"files\"")
		req.Header.Set("Origin", "https://ujomusic.com")
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, ";")
	contentType = parts[0]

	nresp := new(IPFSResponse)

	nresp.Output = resp.Body
	if resp.StatusCode >= http.StatusBadRequest {
		e := &IPFSError{
			Command: r.Command,
		}
		switch {
		case resp.StatusCode == http.StatusNotFound:
			e.Message = "command not found"
		case contentType == "text/plain":
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response read error: %s\n", err)
			}
			e.Message = string(out)
		case contentType == "application/json":
			if err = json.NewDecoder(resp.Body).Decode(e); err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response unmarshall error: %s\n", err)
			}
		default:
			fmt.Fprintf(os.Stderr, "ipfs-shell: warning! unhandled response encoding: %s", contentType)
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: response read error: %s\n", err)
			}
			e.Message = fmt.Sprintf("unknown ipfs-shell error encoding: %q - %q", contentType, out)
		}
		nresp.Error = e
		nresp.Output = nil

		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	return nresp, nil
}

func getURL(r *shell.Request) string {
	values := make(url.Values)
	for _, arg := range r.Args {
		values.Add("arg", arg)
	}
	for k, v := range r.Opts {
		values.Add(k, v)
	}

	return fmt.Sprintf("%s/%s?%s", r.ApiBase, r.Command, values.Encode())
}
