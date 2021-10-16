package api

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/consensys/ugo/pkg/db"
	"github.com/consensys/ugo/pkg/lg"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	_ "github.com/lib/pq" // postgres driver
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

var tokenString string
var adminTokenString string
var admin bool

func init() {
	secret := os.Getenv("UJO_API_SECRET")
	if secret == "" {
		panic("API secret has not been set")
	}

	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)
	_, tokenString, _ = tokenAuth.Encode(
		jwtauth.Claims{
			"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049", "exp": time.Now().Add(time.Hour * 365).Unix(),
		},
	)

	adminTokenAuth := jwtauth.New("HS256", []byte(secret), nil)
	_, adminTokenString, _ = adminTokenAuth.Encode(
		jwtauth.Claims{
			"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049", "admin": true, "exp": time.Now().Add(time.Hour * 365).Unix(),
		},
	)

	os.Setenv("TEST_ENV", "Yes")
}

func setup() {
	connectURI := os.Getenv("TEST_DATABASE_URL")
	if connectURI == "" {
		connectURI = "postgres://ubuntu@localhost/ujo-test?sslmode=disable"
	}

	db, err := sql.Open("postgres", connectURI)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "../../db/migrations",
	}

	// Make sure database is clear
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		panic(err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}

	log.Printf("Database Initialized. Applied %d migrations!\n", n)
}

func teardown() {
	migrations := &migrate.FileMigrationSource{
		Dir: "../../db/migrations",
	}

	connectURI := os.Getenv("TEST_DATABASE_URL")
	if connectURI == "" {
		connectURI = "postgres://ubuntu@localhost/ujo-test?sslmode=disable"
	}

	db, err := sql.Open("postgres", connectURI)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	_, err = migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		panic(err)
	}
	log.Printf("Database Cleared")
}

func server() *httptest.Server {
	logger := logrus.New()
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger

	serverCtx := context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Starting API")

	connectURI := os.Getenv("TEST_DATABASE_URL")
	if connectURI == "" {
		connectURI = "postgres://ubuntu@localhost/ujo-test?sslmode=disable"
	}
	db, err := db.NewDB(connectURI)
	if err != nil {
		log.Panic(err)
	}

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(lg.RequestLogger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("你好, Ujo!"))
	})

	r.Mount("/api/persons", PersonResource{DB: db}.Routes())
	r.Mount("/api/musicgroups", MusicGroupResource{DB: db}.Routes())
	r.Mount("/api/musicreleases", MusicReleaseResource{DB: db}.Routes())
	r.Mount("/api/purchases", PurchaseResource{DB: db}.Routes())
	r.Mount("/api/reports", ReportResource{DB: db}.Routes())
	r.Mount("/api/musicrecordings", MusicRecordingResource{DB: db}.Routes())
	r.Mount("/api/musicplaylists", MusicPlaylistResource{DB: db}.Routes())
	r.Mount("/api/signedmessages", SignedMessageResource{DB: db}.Routes())
	r.Mount("/api/playevents", PlayEventResource{DB: db}.Routes())
	r.Mount("/api/payevents", PayEventResource{DB: db}.Routes())

	ts := httptest.NewServer(r)
	return ts
}

// Handles setup and tear down of test environment
func TestMain(m *testing.M) {
	code := m.Run()
	// teardown()
	log.Print("Exiting api tests...")
	os.Exit(code)
}

func TestRootResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	resp := testRequest(
		t, ts, "GET", "/", nil,
	)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if string(respBody) != "你好, Ujo!" {
		t.Fatalf(string(respBody))
	}
}

func TestMusicGroupResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close() // Testing POST

	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)

	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201.")
	}

	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(PersonPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing POST")
	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201. %s", resp.Status)
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing GET")
	resp = testRequest(
		t, ts, "GET", "/api/musicgroups", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	}
	defer resp.Body.Close()

	log.Println("Testing GET by ID")
	resp = testRequest(
		t, ts, "GET", "/api/musicgroups/1", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	}
	defer resp.Body.Close()

	log.Println("Testing PUT")
	resp = testRequest(
		t, ts, "PUT", "/api/musicgroups/1", bytes.NewBuffer(MusicGroupPut),
	)
	if resp.StatusCode != 200 {
		t.Fatalf("PUT failed, should be 200")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPut), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}

	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("PUT failed, request and response content bodies don't match")
	}

	log.Println("Testing PUT with member added")
	resp = testRequest(
		t, ts, "PUT", "/api/musicgroups/1", bytes.NewBuffer(MusicGroupPutWithMemberAdded),
	)
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPutWithMemberAdded), resp.Body)

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}

	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("PUT failed, request and response content bodies don't match")
	}

	defer resp.Body.Close()

	log.Println("Testing PUT with member image added")
	resp = testRequest(
		t, ts, "PUT", "/api/musicgroups/1", bytes.NewBuffer(MusicGroupPutWithMemberImageAdded),
	)
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPutWithMemberImageAdded), resp.Body)

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}

	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("PUT failed, request and response content bodies don't match")
	}

	defer resp.Body.Close()

	log.Println("Testing PUT with member removed")
	resp = testRequest(
		t, ts, "PUT", "/api/musicgroups/1", bytes.NewBuffer(MusicGroupPutWithMemberRemoved),
	)
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPutWithMemberRemoved), resp.Body)

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}

	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("PUT failed, request and response content bodies don't match")
	}

	defer resp.Body.Close()

	log.Println("Testing DELETE")
	resp = testRequest(
		t, ts, "DELETE", "/api/musicgroups/1", nil,
	)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 || string(respBody) != "DELETE" {
		t.Fatalf("DELETE failed")
	}
}

func TestMusicRecordingResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(PersonPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	// MusicRecording objects require a MusicGroup so this creates one
	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	resp = testRequest(
		t, ts, "POST", "/api/musicrecordings", bytes.NewBuffer(MusicRecordingPost),
	)

	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicRecordingPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing GET")
	resp = testRequest(
		t, ts, "GET", "/api/musicrecordings", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Testing GET by id")
	resp = testRequest(
		t, ts, "GET", "/api/musicrecordings/1", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET by id failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Testing PUT")
	resp = testRequest(
		t, ts, "PUT", "/api/musicrecordings/1", bytes.NewBuffer(MusicRecordingPut),
	)
	if resp.StatusCode != 200 {
		t.Fatalf("PUT failed, should be 200")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicRecordingPut), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("PUT failed, request and response content bodies don't match")
	}

	log.Println("Testing DELETE")
	resp = testRequest(
		t, ts, "DELETE", "/api/musicrecordings/1", nil,
	)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 || string(respBody) != "DELETE" {
		t.Fatalf("DELETE failed")
	}
}

func TestPurchasesResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(PersonPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	// Purchase objects require MusicRelease and MusicRelease require MusicGroup
	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	// create MusicRelease
	resp = testRequest(
		t, ts, "POST", "/api/musicreleases", bytes.NewBuffer(MusicReleasePost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicReleasePost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing GET")
	resp = testRequest(
		t, ts, "GET", "/api/purchases", nil,
	)
	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Testing POST")
	resp = testRequest(
		t, ts, "POST", "/api/purchases", bytes.NewBuffer(PurchasePost),
	)

	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(PurchasePost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}

	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing GET by ID")
	resp = testRequest(
		t, ts, "GET", "/api/purchases/1", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	// PUT is not implemented for Purchases

	log.Println("Testing DELETE")
	resp = testRequest(
		t, ts, "DELETE", "/api/purchases/1", nil,
	)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || string(respBody) != "DELETE" {
		t.Fatalf("DELETE failed")
	}
}

func TestMusicReleaseResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(PersonPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	// MusicRelease objects require a MusicGroup so this creates one
	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing GET")
	resp = testRequest(
		t, ts, "GET", "/api/musicreleases", nil,
	)
	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Testing POST")
	resp = testRequest(
		t, ts, "POST", "/api/musicreleases", bytes.NewBuffer(MusicReleasePost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicReleasePost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing GET by ID")
	resp = testRequest(
		t, ts, "GET", "/api/musicreleases/1", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Testing PUT")
	resp = testRequest(
		t, ts, "PUT", "/api/musicreleases/1", bytes.NewBuffer(MusicReleasePut),
	)
	if resp.StatusCode != 200 {
		t.Fatalf("PUT failed, should be 200")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicReleasePut), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("PUT failed, request and response content bodies don't match")
	}

	log.Println("Testing DELETE")
	resp = testRequest(
		t, ts, "DELETE", "/api/musicreleases/1", nil,
	)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || string(respBody) != "DELETE" {
		t.Fatalf("DELETE failed")
	}
}

func TestMusicPlaylistResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	defer resp.Body.Close()

	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	defer resp.Body.Close()

	log.Println("Testing GET")
	resp = testRequest(
		t, ts, "GET", "/api/musicplaylists", nil,
	)
	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Creating recording")
	resp = testRequest(
		t, ts, "POST", "/api/musicrecordings", bytes.NewBuffer(MusicRecordingPost),
	)

	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	defer resp.Body.Close()

	log.Println("Testing POST")
	resp = testRequest(
		t, ts, "POST", "/api/musicplaylists", bytes.NewBuffer(MusicPlaylistPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	defer resp.Body.Close()

	log.Println("Testing GET by ID")
	resp = testRequest(
		t, ts, "GET", "/api/musicplaylists/1", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Testing PUT")
	resp = testRequest(
		t, ts, "PUT", "/api/musicplaylists/1", bytes.NewBuffer(MusicPlaylistPut),
	)
	if resp.StatusCode != 200 {
		t.Fatalf("PUT failed, should be 200")
	}

	defer resp.Body.Close()

	log.Println("Testing DELETE")
	resp = testRequest(
		t, ts, "DELETE", "/api/musicplaylists/1", nil,
	)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || string(respBody) != "DELETE" {
		t.Fatalf("DELETE failed")
	}
}

func TestReportResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(PersonPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	// MusicRelease objects require a MusicGroup so this creates one
	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicGroupPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing POST")
	resp = testRequest(
		t, ts, "POST", "/api/musicreleases", bytes.NewBuffer(MusicReleasePost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(MusicReleasePost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)
	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Creating Report")
	resp = testRequest(
		t, ts, "POST", "/api/reports", bytes.NewBuffer(ReportPost),
	)

	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(ReportPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq = reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	mResp, errResp = respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mResp)

	if testValues(mReq, mResp) != true {
		t.Fatalf("POST failed, request and response content bodies don't match")
	}

	log.Println("Testing unauthorized GET")
	resp = testRequest(
		t, ts, "GET", "/api/reports", nil,
	)
	if resp.StatusCode != 401 {
		t.Fatalf("Unauthorized GET failed, should be 401")
	}

	admin = true
	log.Println("Testing GET")
	resp = testRequest(
		t, ts, "GET", "/api/reports", nil,
	)
	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	}
	defer resp.Body.Close()

	log.Println("Testing GET by ID")
	resp = testRequest(
		t, ts, "GET", "/api/reports/1", nil,
	)
	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}

	log.Println("Testing PUT")
	resp = testRequest(
		t, ts, "PUT", "/api/reports/1", bytes.NewBuffer(ReportPut),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("PUT failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(ReportPut), resp.Body)
	defer resp.Body.Close()

	log.Println("Testing PUT to resolve")
	resp = testRequest(
		t, ts, "PUT", "/api/reports/1/resolve", bytes.NewBuffer(ReportPut),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("PUT failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(ReportPut), resp.Body)
	if respJSON.Get("state").MustString() != "resolved" {
		t.Fatalf("PUT to deactivate failed, state should be \"resolved\"")
	}
	defer resp.Body.Close()

	log.Println("Testing PUT to deactivate")
	resp = testRequest(
		t, ts, "PUT", "/api/reports/1/deactivate", bytes.NewBuffer(ReportPut),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("PUT failed, should be 201")
	}
	reqJSON, respJSON = simplejsonGenerator(bytes.NewBuffer(ReportPut), resp.Body)
	if respJSON.Get("state").MustString() != "deactivated" {
		t.Fatalf("PUT to deactivate failed, state should be \"deactivated\"")
	}

	log.Println("Testing MusicRelease's GET inactive")
	resp = testRequest(
		t, ts, "GET", "/api/musicreleases/inactive", nil,
	)

	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	} else {
		defer resp.Body.Close()
	}
	// TODO: Check if active flag on the music release got set to false
	defer resp.Body.Close()
}

func TestSignedMessageResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	log.Println("Testing POST")
	resp := testRequest(
		t, ts, "POST", "/api/signedmessages", bytes.NewBuffer(SignedMessagePost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(SignedMessagePost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	_, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}

	log.Println("Testing GET")
	resp = testRequest(
		t, ts, "GET", "/api/signedmessages", nil,
	)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	}

	log.Println("Testing GET by ID")
	resp = testRequest(
		t, ts, "GET", "/api/signedmessages/1", nil,
	)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("GET failed, should be 200")
	}
}

func TestPlayEventResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	// PlayEvent requires a person
	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	// PlavEvent objects require a MusicGroup so this creates one
	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	// PlayEvent requires a musicrecording
	resp = testRequest(
		t, ts, "POST", "/api/musicrecordings", bytes.NewBuffer(MusicRecordingPost),
	)

	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	// Start testing PlayEvent
	log.Println("Testing POST")
	resp = testRequest(
		t, ts, "POST", "/api/playevents", bytes.NewBuffer(PlayEventPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(PlayEventPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	_, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
}

func TestPayEventResource(t *testing.T) {
	setup()
	ts := server()
	defer ts.Close()

	// PayEvent requires a person
	log.Println("Creating Person")
	resp := testRequest(
		t, ts, "POST", "/api/persons", bytes.NewBuffer(PersonPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	// PavEvent objects require a MusicGroup so this creates one
	resp = testRequest(
		t, ts, "POST", "/api/musicgroups", bytes.NewBuffer(MusicGroupPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	// PayEvent requires a musicrecording
	resp = testRequest(
		t, ts, "POST", "/api/musicrecordings", bytes.NewBuffer(MusicRecordingPost),
	)

	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}

	// Start testing PayEvent
	log.Println("Testing POST")
	resp = testRequest(
		t, ts, "POST", "/api/payevents", bytes.NewBuffer(PayEventPost),
	)
	if resp.StatusCode != 201 {
		t.Fatalf("POST failed, should be 201")
	}
	reqJSON, respJSON := simplejsonGenerator(bytes.NewBuffer(PayEventPost), resp.Body)
	defer resp.Body.Close()

	mReq, errReq := reqJSON.Map()
	if errReq != nil {
		t.Fatalf("ERROR")
	}
	removeUntestableKeys(mReq)

	_, errResp := respJSON.Map()
	if errResp != nil {
		t.Fatalf("ERROR")
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) *http.Response {
	// TODO - Set AUTH Token in modular way
	req, err := http.NewRequest(method, ts.URL+path, body)

	if admin {
		adminAuthHeader := fmt.Sprintf("BEARER %s", adminTokenString)
		req.Header.Set("Authorization", adminAuthHeader)
	} else {
		authHeader := fmt.Sprintf("BEARER %s", tokenString)
		req.Header.Set("Authorization", authHeader)
	}

	if err != nil {
		t.Fatal(err)
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return resp
}

// converts the []byte HTTP request and response into simplejson.Json
func simplejsonGenerator(reqBody, respBody io.Reader) (*simplejson.Json, *simplejson.Json) {
	reqJSON, reqErr := simplejson.NewFromReader(reqBody)
	if reqErr != nil {
		fmt.Println("Error when converting HTTP request to simplejson: ", reqErr)
		return nil, nil
	}

	respJSON, respErr := simplejson.NewFromReader(respBody)
	if respErr != nil {
		fmt.Println("Error when converting HTTP response to simplejson: ", respErr)
		return nil, nil
	}

	return reqJSON, respJSON
}

// remove unoverlapping entries in request and response
func removeUntestableKeys(m map[string]interface{}) {
	for k, v := range m {
		if v == nil {
			continue
		}
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice: // value of 'members' is slice of interfaces
			for _, i := range v.([]interface{}) { // type assertion required
				j, ok := i.(map[string]interface{})
				if ok == true {
					removeUntestableKeys(j)
				}
			}
		case reflect.Map: // value of 'image' is map
			j, ok := v.(map[string]interface{})
			if ok == true {
				removeUntestableKeys(j)
			}
		case reflect.String:
			if k == "createdAt" || k == "updatedAt" || k == "id" || k == "email" {
				delete(m, k)
			}
		}
	}
}

// loop through each key in request and check whether the value in response aligns
func testValues(reqM, respM map[string]interface{}) bool {
	for k, v := range reqM { // loop over request but not response because we don't remove terms with nil as value
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			for y, i := range v.([]interface{}) {
				j, ok := i.(map[string]interface{})
				if ok == true {
					x, ok := respM[k].([]interface{})[y].(map[string]interface{}) // find the corresponding term in response with type assertion
					if ok == true {
						// Only return if testValues returns false so all array items get tested
						if !testValues(j, x) {
							return false
						}
					} else {
						return false
					}
				}
			}
		case reflect.Map:
			j, ok := v.(map[string]interface{})
			if ok == true {
				x, ok := respM[k].(map[string]interface{})
				if ok == true {
					return testValues(j, x)
				}
				return false
			}
		case reflect.String:
			if reqM[k] != respM[k] {
				fmt.Println("expected: ", reqM[k])
				fmt.Println("got: ", respM[k])
				return false
			}
		}
	}
	return true
}

var PersonPost = []byte(
	`{
		"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
		"@context":"http://schema.org",
		"@type":"Person",
		"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
		"email":"michelle.wu@consensys.net",
		"familyName":"Shibata",
		"givenName":"パウダー"
	}`,
)

var MusicGroupPost = []byte(
	`{
		"cid": "Qm",
		"@type":"MusicGroup",
		"@context":"http://schema.org/",
		"name":"Powder",
		"description":"Techno",
		"email":"powder@gmail.com",
		"members":[
			{
				"id": 1,
				"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
				"email":"michelle.wu@consensys.net",
				"familyName":"Shibata",
				"givenName":"パウダー",
				"description":"DJ",
				"percentageShares": 100,
				"musicgroupAdmin": true
			}
		],
		"image":{
			"cid":"ImB",
			"@type":"ImageObject",
			"@context":"http://schema.org/",
			"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
			"encodingFormat":"image/jpeg"
		}
	}`,
)

var MusicGroupPut = []byte(
	`{
		"id": 1,
		"cid":"vvv",
		"@type":"MusicGroup",
		"@context":"http://schema.org/",
		"name":"Powder",
		"description":"Techno",
		"email":"powder@gmail.com",
		"members":[
			{
				"id": 1,
				"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
				"email":"email@me.com",
				"familyName":"Shibata updated",
				"givenName":"パウダー"
			}
		],
		"image":{
			"id": 1,
			"cid":"ImB",
			"@type":"ImageObject",
			"@context":"http://schema.org/",
			"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
			"encodingFormat":"image/jpeg"
		}
	}`,
)

var MusicGroupPutWithMemberAdded = []byte(
	`{
		"id": 1,
		"cid":"vvv",
		"@type":"MusicGroup",
		"@context":"http://schema.org/",
		"name":"Powder",
		"description":"Techno",
		"email":"powder@gmail.com",
		"members":[
			{
				"id": 1,
				"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
				"email":"email@me.com",
				"familyName":"Shibatas",
				"givenName":"パウダー"
			},
			{
				"cid":"newmember",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f12345",
				"email":"email2@me.com",
				"familyName":"Member",
				"givenName":"New"
			}
		],
		"image":{
			"id": 1,
			"cid":"ImB",
			"@type":"ImageObject",
			"@context":"http://schema.org/",
			"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
			"encodingFormat":"image/jpeg"
		}
	}`,
)

var MusicGroupPutWithMemberImageAdded = []byte(
	`{
		"id": 1,
		"cid":"vvv",
		"@type":"MusicGroup",
		"@context":"http://schema.org/",
		"name":"Powder",
		"description":"Techno",
		"email":"powder@gmail.com",
		"members":[
			{
				"id": 1,
				"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
				"email":"email@me.com",
				"familyName":"Shibatas",
				"givenName":"パウダー"
			},
			{
				"id": 2,
				"cid":"newmember",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f12345",
				"email":"email2@me.com",
				"familyName":"Member",
				"givenName":"New",
				"image":{
					"cid":"ImBnewmembersssss",
					"@type":"ImageObject",
					"@context":"http://schema.org/",
					"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
					"encodingFormat":"image/jpeg"
				}
			}
		],
		"image":{
			"id": 1,
			"cid":"ImB",
			"@type":"ImageObject",
			"@context":"http://schema.org/",
			"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
			"encodingFormat":"image/jpeg"
		}
	}`,
)

var MusicGroupPutWithMemberRemoved = []byte(
	`{
		"id": 1,
		"cid":"vvv",
		"@type":"MusicGroup",
		"@context":"http://schema.org/",
		"name":"Powder",
		"description":"Techno",
		"email":"powder@gmail.com",
		"members":[
			{
				"id": 1,
				"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
				"email":"email@me.com",
				"familyName":"Shibatas",
				"givenName":"パウダー"
			}
		],
		"image":{
			"id": 1,
			"cid":"ImB",
			"@type":"ImageObject",
			"@context":"http://schema.org/",
			"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
			"encodingFormat":"image/jpeg"
		}
	}`,
)

var MusicRecordingPost = []byte(
	`{
		"cid": "Qm",
		"@type": "MusicRecording",
		"@context": "http://schema.org/",
		"image":{
			"id": 1,
			"cid":"ImB",
			"@type":"ImageObject",
			"@context":"http://schema.org/",
			"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
			"encodingFormat":"image/jpeg"
		},
		"name": "The Dopeness",
		"duration": "3:12",
		"isrc": "",
		"byArtist": {
			"id": 1,
			"cid": "Qm",
			"@type":"MusicGroup",
			"@context":"http://schema.org/",
			"name": "Powder",
			"description":"Techno",
			"email":"powder@gmail.com",
			"members":[
				{
					"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
					"@context":"http://schema.org",
					"@type":"Person",
					"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
					"email":"email@me.com",
					"familyName":"Shibata",
					"givenName":"パウダー"
				}
			],
			"image":{
				"cid":"ImB",
				"@type":"ImageObject",
				"@context":"http://schema.org/",
				"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
				"encodingFormat":"image/jpeg"
			}
		},
		"audio": {
			"cid": "Qm",
			"@type": "AudioObject",
			"@context": "http://schema.org/",
			"contentURL": "http://ipfs.io/ipfs/QmZ",
			"encodingFormat": "image/png"
		},
		"recordingOf": {
			"cid": "Qm",
			"@type": "MusicComposition",
			"@context": "http://schema.org/",
			"name": "The Dopeness",
			"composer":[
				{
					"id": 1,
					"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
					"@context":"http://schema.org",
					"@type":"Person",
					"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
					"email":"email@me.com",
					"familyName":"Shibata",
					"givenName":"パウダー"
				}
			]
		},
		"rights": [
			{
				"cid": "zdpuB12JMCxfPgYxMopQnc8zgpb2W6aqKE85wD3QWgvBfPUWk",
				"@context": "http://coalaip.org",
				"@type": "Right",
				"party": {
					"id": 1, 
					"cid": "zdpuAw6fhie6kVedNUegcqDKoKXA1EwKMD6MtxtUpbj9WjLRC",
					"@context": "http://coalaip.org",
					"@type": "Person",
					"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
					"familyName": "Shibata",
					"givenName": "Jon"
				},
				"percentageShares": 100,
				"validFrom": "",
				"validThrough": ""
			}
		]
	}`,
)

var MusicRecordingPut = []byte(
	`{
		"id": 1,
		"cid": "Qm",
		"@type": "MusicRecording",
		"@context": "http://schema.org/",
		"image":{
			"id": 1,
			"cid":"ImB",
			"@type":"ImageObject",
			"@context":"http://schema.org/",
			"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
			"encodingFormat":"image/jpeg"
		},
		"name": "The Dopeness",
		"duration": "3:12",
		"isrc": "",
		"byArtist": {
			"id": 1,
			"cid":"Qm",
			"@type":"MusicGroup",
			"@context":"http://schema.org/",
			"name": "Powder",
			"description":"Techno",
			"email":"powder@gmail.com",
			"members":[
				{
					"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
					"@context":"http://schema.org",
					"@type":"Person",
					"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
					"email":"email@me.com",
					"familyName":"Shibata",
					"givenName":"パウダー"
				}
			],
			"image":{
				"cid":"ImB",
				"@type":"ImageObject",
				"@context":"http://schema.org/",
				"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
				"encodingFormat":"image/jpeg"
			}
		},
		"audio": {
			"id": 1,
			"cid": "Qm",
			"@type": "AudioObject",
			"@context": "http://schema.org/",
			"contentURL": "http://ipfs.io/ipfs/QmZ",
			"encodingFormat": "image/png"
		},
		"recordingOf": {
			"id": 1,
			"cid": "Qm",
			"@type": "MusicComposition",
			"@context": "http://schema.org/",
			"name": "The Dopeness",
			"composer":[
				{
					"id": 1,
					"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
					"@context":"http://schema.org",
					"@type":"Person",
					"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
					"email":"email@me.com",
					"familyName":"Shibata",
					"givenName":"パウダー"
				}
			]
		}
	}`,
)

var MusicReleasePost = []byte(
	`{
			"cid": "Qm",
			"@type": "MusicRelease",
			"@context": "http://schema.org/",
			"catalogNumber": "1",
			"musicReleaseFormat": "CD",
			"price": 10,
			"releaseOf": {
				"cid": "Qm",
				"@type": "MusicAlbum",
				"@context": "http://schema.org/",
				"albumReleaseType": "EP",
				"albumProductionType": "studio",
				"byArtist": {
					"id": 1,
					"cid":"Qm",
					"@type":"MusicGroup",
					"@context":"http://schema.org/",
					"name": "Powder",
					"description":"Techno",
					"email":"powder@gmail.com",
					"members":[
						{
							"id": 1,
							"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
							"@context":"http://schema.org",
							"@type":"Person",
							"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
							"email":"michelle.wu@consensys.net",
							"familyName":"Shibata",
							"givenName":"パウダー"
						}
					],
					"image":{
						"id": 1,
						"cid":"ImB",
						"@type":"ImageObject",
						"@context":"http://schema.org/",
						"contentURL": "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
						"encodingFormat":"image/jpeg"
					}
				},
				"name": "Random Album Title",
				"tracks": [
					{
						"cid": "Qm",
						"@type": "MusicRecording",
						"@context": "http://schema.org/",
						"name": "The Dopeness",
						"audio": {
							"cid": "Qm",
							"@type": "AudioObject",
							"@context": "http://schema.org/",
							"contentURL": "http://ipfs.io/ipfs/QmZ",
							"encodingFormat": "image/png"
						},
						"recordingOf": {
							"cid": "Qm",
							"@type": "MusicComposition",
							"@context": "http://schema.org/",
							"composer":[
								{
									"id": 1,
									"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
									"@context":"http://schema.org",
									"@type":"Person",
									"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
									"email":"email@me.com",
									"familyName":"Shibata",
									"givenName":"パウダー"
								}
							],
							"name": "The Dopeness"
						}
					}
				]
			},
			"image": {
				"cid": "Qm",
				"@type": "ImageObject",
				"@context": "http://schema.org/",
				"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
				"encodingFormat":"image/jpeg"
			}
		}`,
)

var MusicReleasePut = []byte(
	`{
			"id": 1,
			"cid": "Qm",
			"@type": "MusicRelease",
			"@context": "http://schema.org/",
			"catalogNumber": "1",
			"musicReleaseFormat": "CD",
			"price": 10,
			"active": true,
			"releaseOf": {
				"id": 1,
				"cid": "Qm",
				"@type": "MusicAlbum",
				"@context": "http://schema.org/",
				"albumReleaseType": "EP",
				"byArtist": {
					"id": 1,
					"cid":"Qm",
					"@type":"MusicGroup",
					"@context":"http://schema.org/",
					"name": "Powder",
					"description":"Techno",
					"email":"powder@gmail.com",
					"members":[
						{
							"id": 1,
							"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
							"@context":"http://schema.org",
							"@type":"Person",
							"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
							"email":"email@me.com",
							"familyName":"Shibata",
							"givenName":"パウダー"
						}
					],
					"image":{
						"id": 1,
						"cid":"Qm",
						"@type":"ImageObject",
						"@context":"http://schema.org/",
						"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
						"encodingFormat":"image/jpeg"
					}
				},
				"name": "Random Album Title",
				"tracks": [
					{
						"id": 1,
						"cid": "Qm",
						"@type": "MusicRecording",
						"@context": "http://schema.org/",
						"name": "The Dopeness",
						"audio": {
							"id": 1,
							"cid": "Qm",
							"@type": "AudioObject",
							"@context": "http://schema.org/",
							"contentURL": "http://ipfs.io/ipfs/QmZ",
							"encodingFormat": "image/png"
						},
						"recordingOf": {
							"id": 1,
							"cid": "Qm",
							"@type": "MusicComposition",
							"@context": "http://schema.org/",
							"name": "The Dopeness"
						}
					}
				]
			},
			"image": {
				"id": 1,
				"cid": "Qm",
				"@type": "ImageObject",
				"@context": "http://schema.org/",
				"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
				"encodingFormat":"image/jpeg"			}
		}`,
)

var MusicPlaylistPost = []byte(
	`{
		"cid": "1",
		"@type": "MusicPlaylist",
		"name": "few_words",
		"@context": "http://coalaip.org",
		"tracks": [
			{
				"id": 1,
				"cid": "Qm",
				"@type": "MusicRecording",
				"@context": "http://schema.org/",
				"createdAt": "2019-01-28T14:48:12.253656-05:00",
				"updatedAt": "2019-01-28T14:48:12.253656-05:00",
				"name": "The Dopeness",
				"genres": [
					"Alternative Rock"
				]
			}
		],
		"byUser": {
			"id": 1,
			"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
			"@context":"http://schema.org",
			"@type":"Person",
			"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
			"email":"michelle.wu@consensys.net",
			"familyName":"Shibata",
			"givenName":"パウダー"
		},
		"image": {
			"id": 1,
			"cid": "p12",
			"@type": "Image",
			"@context": "http://coalaip.org",
			"contentURL": "https://up.jpg",
			"encodingFormat": "image/jpeg"
		}
	}`,
)

var MusicPlaylistPut = []byte(
	`{
		"id": 1,
		"cid": "1",
		"@type": "MusicPlaylist",
		"name": "few_words2",
		"@context": "http://coalaip.org",
		"tracks": [
			{
				"id": 1,
				"cid": "Qm",
				"@type": "MusicRecording",
				"@context": "http://schema.org/",
				"createdAt": "2019-01-28T14:48:12.253656-05:00",
				"updatedAt": "2019-01-28T14:48:12.253656-05:00",
				"name": "The Dopeness",
				"genres": [
					"Alternative Rock"
				]
			}
		],
		"byUser": {
			"id": 1,
			"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
			"@context":"http://schema.org",
			"@type":"Person",
			"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
			"email":"michelle.wu@consensys.net",
			"familyName":"Shibata",
			"givenName":"パウダー"
		},
		"image": {
			"id": 1,
			"cid": "p12",
			"@type": "Image",
			"@context": "http://coalaip.org",
			"contentURL": "https://up.jpg",
			"encodingFormat": "image/jpeg"
		}
	}`,
)

var ReportPost = []byte(
	`{
		"reason": "Copyright infringement",
		"message": "This release infringes copyright",
		"email": "email@me.com",
		"reporter_id": 1,
		"musicrelease_id": 1
	}`,
)

var ReportPut = []byte(
	`{
		"response": "Yerp",
		"reason": "Copyright infringement",
		"message": "This release infringes copyright",
		"email": "email@me.com",
		"reporter_id": 1,
		"musicrelease_id": 1
	}`,
)

var PurchasePost = []byte(
	`{
			"cid": "Qm",
			"@type": "Purchase",
			"@context": "http://schema.org/",
			"txHash": "0x87476c90a1b90a30a6dae5b9fb85ec82aa7dc2fe78ae3db7145c0382dd7152e9",
			"buyer": {
				"id": 1,
				"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
				"@context":"http://schema.org",
				"@type":"Person",
				"ethereumAddress": "0x0431aac01788f21466ce4ad27368dc4b03f89049",
				"email":"email@me.com",
				"familyName":"Shibata",
				"givenName":"パウダー"
			},
			"musicRelease": {
				"id": 1,
				"cid": "Qm",
				"@type": "MusicRelease",
				"@context": "http://schema.org/",
				"catalogNumber": "1",
				"musicReleaseFormat": "CD",
				"price": 10,
				"releaseOf": {
					"id": 1,
					"cid": "Qm",
					"@type": "MusicAlbum",
					"@context": "http://schema.org/",
					"albumReleaseType": "EP",
					"albumProductionType": "studio",
					"byArtist": {
						"id": 1,
						"cid":"Qm",
						"@type":"MusicGroup",
						"@context":"http://schema.org/",
						"name": "Powder",
						"description":"Techno",
						"email":"powder@gmail.com",
						"members":[
							{
								"id": 1,
								"cid":"zdpxApt2XrHHXTtztbn9Vc9ynCEeasdfoZPeLW2et9Me6AEa",
								"@context":"http://schema.org",
								"@type":"Person",
								"ethereumAddress": "0x1",
								"email":"email@me.com",
								"familyName":"Shibata",
								"givenName":"パウダー"
							}
						],
						"image":{
							"id": 1,
							"cid":"ImB",
							"@type":"ImageObject",
							"@context":"http://schema.org/",
							"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
							"encodingFormat":"image/jpeg"
						}
					},
					"name": "Random Album Title",
					"tracks": [
						{
							"id": 1,
							"cid": "Qm",
							"@type": "MusicRecording",
							"@context": "http://schema.org/",
							"name": "The Dopeness",
							"audio": {
								"id": 1,
								"cid": "Qm",
								"@type": "AudioObject",
								"@context": "http://schema.org/",
								"contentURL": "http://ipfs.io/ipfs/QmZ",
								"encodingFormat": "image/png"
							},
							"recordingOf": {
								"id": 1,
								"cid": "Qm",
								"@type": "MusicComposition",
								"@context": "http://schema.org/",
								"name": "The Dopeness"
							}
						}
					]
				},
				"image": {
					"id": 1,
					"cid": "Qm",
					"@type": "ImageObject",
					"@context": "http://schema.org/",
					"contentURL":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/CarsonWentz11.jpg/800px-CarsonWentz11.jpg",
					"encodingFormat":"image/jpeg"				}
			}
		}`,
)

var SignedMessagePost = []byte(
	`{
		"message": "data",
		"@type": "SignedMessage",
		"ethereum-address": "0x3362838B336070c69123d35BC537215EEfdfD59f",
		"signedmessage": "0x592fa743889fc7f92ac2a37bb1f5ba1daf2a5c84741ca0e0061d243a2e6707ba"
	}`,
)

var PlayEventPost = []byte(
	`{
		"playedby_id": 1,
		"musicrecording_id": 1
	}`,
)

var PayEventPost = []byte(
	`{
		"playedby_id": 1,
		"beneficiary_id": 1,
		"musicrecording_id": 1,
		"amount": 0.01,
		"link": "redeemable-link"
	}`,
)
