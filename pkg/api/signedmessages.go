package api

import (
	"net/http"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
)

// SignedMessageResource contains REST handlers
type SignedMessageResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs SignedMessageResource) Routes() chi.Router {
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
			r.Use(jwtauth.Authenticator)
			r.Post("/", rs.Post) // POST /{resource} - create a new resource and persist it
		})

		r.Route("/{id}", func(r chi.Router) {
			// Public routes
			r.Get("/", rs.Get) // GET /{resource}/{id} - read a single resource by id
		})

		r.Route("/signedmessage/{signedmessage}", func(r chi.Router) {
			// Public routes
			r.Get("/", rs.SignedMessage) // GET /{resource}/{signedmessage} - read a single resource by signedmessage
		})
	})

	return r
}

type SignedMessageListRequest struct {
	Limit  int `api:"limit,@query"  validate:"gte=0"`
	Offset int `api:"offset,@query" validate:"gte=0"`
}

// List all of the given SignedMessage objects
func (rs SignedMessageResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List SignedMessage")

	var req SignedMessageListRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List SignedMessage (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	signedmessages, err := (&models.SignedMessage{}).All(rs.DB, &models.SelectQuery{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		lg.Errorf("Error in List SignedMessage (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewSignedMessageListResponse(signedmessages))
}

type SignedMessageGetRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Get SignedMessage by id
func (rs SignedMessageResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get SignedMessage")

	var req SignedMessageGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get SignedMessage (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	signedmessage, err := (&models.SignedMessage{}).Get(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Get SignedMessage (%v)", err)
		render.Render(w, r, Error404)
		return
	}

	render.JSON(w, r, signedmessage)
}

type SignedMessageGetBySignedMessageRequest struct {
	SignedMessage string `api:"signedmessage,@url_param" validate:"gte=0"`
}

// SignedMessage gets by signedmessage
func (rs SignedMessageResource) SignedMessage(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get SignedMessage")

	var req SignedMessageGetBySignedMessageRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get SignedMessage (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	signedmessage, err := (&models.SignedMessage{}).GetBySignedMessage(rs.DB, req.SignedMessage)
	if err != nil {
		lg.Errorf("Error in Get SignedMessage (%v)", err)
		render.Render(w, r, Error404)
		return
	}

	render.JSON(w, r, signedmessage)
}

type SignedMessagePostRequest struct {
	Body models.SignedMessage `api:"body,@body"`
}

// Post a SignedMessage
func (rs SignedMessageResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post SignedMessage")

	var req SignedMessagePostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post SignedMessage (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	signedmessageID, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in Post SignedMessage (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	signedmessage, err := (&models.SignedMessage{}).Get(rs.DB, signedmessageID)
	if err != nil {
		lg.Errorf("Error in Post SignedMessage (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewSignedMessageResponse(signedmessage))
}

// SignedMessageRequest structure
// Any fields to be overidden go here
type SignedMessageRequest struct {
	*models.SignedMessage
}

// Bind pre-processes any fields after a decode
func (req *SignedMessageRequest) Bind(r *http.Request) error {
	return nil
}

// SignedMessageResponse structure
// Add any extra fields to the response here
type SignedMessageResponse struct {
	*models.SignedMessage
}

// NewSignedMessageResponse creates a response with the model plus any other data
func NewSignedMessageResponse(obj *models.SignedMessage) *SignedMessageResponse {
	resp := &SignedMessageResponse{SignedMessage: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *SignedMessageResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// SignedMessageListResponse is an array of response objects
type SignedMessageListResponse []*SignedMessageResponse

// NewSignedMessageListResponse creates a ListResponse
func NewSignedMessageListResponse(signedmessages []*models.SignedMessage) []render.Renderer {
	list := []render.Renderer{}
	for _, signedmessage := range signedmessages {
		list = append(list, NewSignedMessageResponse(signedmessage))
	}
	return list
}
