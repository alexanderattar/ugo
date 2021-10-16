package api

import (
	"net/http"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// PlayEventResource contains REST handlers
type PlayEventResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs PlayEventResource) Routes() chi.Router {
	r := chi.NewRouter()
	// Seek, verify and validate JWT tokens
	r.Use(jwtauth.Verifier(TokenAuth))

	// Protected routes
	r.Group(func(r chi.Router) {
		// Private routes
		r.Group(func(r chi.Router) {
			// Handle valid / invalid tokens
			// This can be modified using the Authenticator method in auth.go
			r.Use(jwtauth.Authenticator)
			r.Post("/", rs.Post) // POST /{resource} - create a new resource and persist it
		})
	})

	return r
}

type PlayEventPostRequest struct {
	Body models.PlayEvent `api:"body,@body"`
}

// Post a PlayEvent
func (rs PlayEventResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post PlayEvent")

	//
	// Get the user's ID from the JWT
	//
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(401), 401)
		return
	} else if token == nil || !token.Valid {
		http.Error(w, http.StatusText(401), 401)
		return
	}

	ethAddr, is := claims["ethereumAddress"].(string)
	if !is {
		http.Error(w, http.StatusText(403), 403)
		return
	}

	user, err := (&models.Person{}).Get(rs.DB, 0, ethAddr)
	if err != nil {
		http.Error(w, http.StatusText(403), 403)
		return
	}

	//
	// Decode the request body
	//
	var req PlayEventPostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post PlayEvent (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	req.Body.PlayedByID = &user.ID

	playeventID, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in POST PlayEvent (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	playevent, err := (&models.PlayEvent{}).Get(rs.DB, playeventID)
	if err != nil {
		lg.Errorf("Error getting PlayEvent (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, playevent)
}
