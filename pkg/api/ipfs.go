// A lot of this code is monkey-patched. Aside from CreateIpfsShell, PutToInfura and uploadToInfura, these types and functions were taken from go-ipfs-api and edited to fit our purposes.

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/utils"
	"github.com/go-chi/render"
	ipfsapi "github.com/ipfs/go-ipfs-api"
)

var req *ipfsapi.Request

// DagPutAndPin puts to infura, returns hash, and puts & pins to infura in a go routine
func DagPutAndPin(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	hash, err := utils.DagPut(body)

	if err != nil {
		lg.Errorln("Error putting dag to Infura:", err)
		render.Render(w, r, Error400(err))
		return
	}

	utils.AddToInfuraDagPinQueue(body, hash, 0)
	render.JSON(w, r, hash)
}

// RecursiveDagPutAndPin adds the request body's top level object to infura and
// recursively puts nested objects, adds the respective cids to each object,
// returns a response, and finally puts & pins in the background.
func RecursiveDagPutAndPin(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	var wg sync.WaitGroup
	var err error

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(body), &data)

	utils.AddCids(data, &wg, &err)
	wg.Wait()

	if err != nil {
		lg.Errorln("Error recursively putting dags:", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.JSON(w, r, data)
}
