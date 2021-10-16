package api

import (
	"net/http"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

// PersonResource contains REST handlers
type PersonResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs PersonResource) Routes() chi.Router {
	r := chi.NewRouter()
	// Seek, verify and validate JWT tokens
	r.Use(jwtauth.Verifier(TokenAuth))

	// Protected routes
	r.Group(func(r chi.Router) {
		// Public routes
		r.Get("/", rs.List)

		// Private routes
		r.Group(func(r chi.Router) {
			// Handle valid / invalid tokens
			// This can be modified using the Authenticator method in auth.go
			// r.Use(jwtauth.Authenticator)
			r.Post("/", rs.Post) // POST /{resource} - create a new resource and persist it
		})

		r.Route("/{id}", func(r chi.Router) {
			// Public routes
			r.Get("/", rs.Get) // GET /{resource}/{id} - read a single resource by id

			// Private routes
			r.Group(func(r chi.Router) {
				// r.Use(jwtauth.Authenticator)
				// r.Use(rs.PersonAuthorization)
				r.Put("/", rs.Put) // PUT /{resource}/{id} - update a single resource by id
				// r.Delete("/", rs.Delete) // DELETE /{resource}/{id} - delete a single resource by id
			})
		})

		r.Get("/cid/{cid}", rs.Cid) // GET /{resource}/cid/{cid} - read a single resource by cid
	})

	return r
}

type PersonListRequest struct {
	Limit        int    `api:"limit,@query"    validate:"gte=0"`
	Offset       int    `api:"offset,@query"   validate:"gte=0"`
	EthereumAddr string `api:"ethereumAddress,@query"`
}

// List all of the given Person objects
func (rs PersonResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List Person")

	var req PersonListRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	persons, err := (&models.Person{}).All(rs.DB, req.EthereumAddr, &models.SelectQuery{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		lg.Errorf("Error in List Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewPersonListResponse(persons))
}

type PersonGetRequest struct {
	ID models.IDType `api:"id,@url_param" validate:"gte=0"`
}

// Get Person by id
func (rs PersonResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get Person")

	var req PersonGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	person, err := (&models.Person{}).Get(rs.DB, req.ID, "")
	if err != nil {
		lg.Errorf("Error in Get Person (%v)", err)
		render.Render(w, r, Error404)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())

	if claims["ethereumAddress"] != person.EthereumAddress && claims["admin"] != true {
		person.Email = nil
	}

	render.JSON(w, r, person)
}

type PersonPostRequest struct {
	Body models.Person `api:"body,@body"`
}

// Post a Person
func (rs PersonResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post Person")

	var req PersonPostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	personID, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in Post Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	person, err := (&models.Person{}).Get(rs.DB, personID, "")
	if err != nil {
		lg.Errorf("Error in Post Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewPersonResponse(person))
}

// Get Person by CID
func (rs PersonResource) Cid(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get Person by CID")
	cid := chi.URLParam(r, "cid")
	if cid == "" {
		render.Render(w, r, Error404)
		return
	}
	obj := &models.Person{}
	person, err := obj.GetByCID(rs.DB, cid)
	if err != nil {
		lg.Errorf("Error in Get Person by CID (%v)", err)
		render.Render(w, r, Error404)
		return
	}
	render.JSON(w, r, person)
}

type PersonPutRequest struct {
	ID   models.IDType `api:"id,@url_param" validate:"gte=0"`
	Body models.Person `api:"body,@body"`
}

// Put Person
func (rs PersonResource) Put(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Put Person")

	var req PersonPutRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Put Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	personID, err := req.Body.Update(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Put Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	person, err := (&models.Person{}).Get(rs.DB, personID, "")
	if err != nil {
		lg.Errorf("Error in Put Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewPersonResponse(person))
}

type PersonDeleteRequest struct {
	ID models.IDType `api:"id,@url_param" validate:"gte=0"`
}

// Delete Person by id
func (rs PersonResource) Delete(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Delete Person")

	var req PersonDeleteRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Delete Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	err := (&models.Person{}).Delete(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Delete Person (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

// PersonRequest structure
// Any fields to be overidden go here
type PersonRequest struct {
	*models.Person
}

// Bind pre-processes any fields after a decode
func (req *PersonRequest) Bind(r *http.Request) error {
	return nil
}

// PersonResponse structure
// Add any extra fields to the response here
type PersonResponse struct {
	*models.Person
}

// NewPersonResponse creates a response with the model plus any other data
func NewPersonResponse(obj *models.Person) *PersonResponse {
	resp := &PersonResponse{Person: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *PersonResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PersonListResponse is an array of response objects
type PersonListResponse []*PersonResponse

// NewPersonListResponse creates a ListResponse
func NewPersonListResponse(persons []*models.Person) []render.Renderer {
	list := []render.Renderer{}
	for _, person := range persons {
		list = append(list, NewPersonResponse(person))
	}
	return list
}

// PersonAuthorization is a default authorizor middleware to enforce resource access
func (rs PersonResource) PersonAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		id := chi.URLParam(r, "id")

		// Selects all ethereum addresss associated with this person
		var ethereumAddresses []string
		query := `
        SELECT person.ethereum_address FROM person
        WHERE person.id = $1;
        `

		err = rs.DB.Select(&ethereumAddresses, query, id)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Check that the user making the request owns the ethereum address that the
		// person was registered with
		authorized := false
		for _, address := range ethereumAddresses {
			if address == claims["ethereumAddress"] {
				authorized = true
			}
		}

		if authorized != true {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authorized, pass it through
		next.ServeHTTP(w, r)
	})
}
