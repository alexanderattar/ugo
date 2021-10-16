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

// MusicRecordingResource contains REST handlers
type MusicRecordingResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs MusicRecordingResource) Routes() chi.Router {
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
			//r.Put("/", rs.Put)
			// Private routes
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Authenticator)
				r.Use(rs.MusicRecordingAuthorization)
				r.Put("/", rs.Put)       // PUT /{resource}/{id} - update a single resource by id
				r.Delete("/", rs.Delete) // DELETE /{resource}/{id} - delete a single resource by id
			})
		})
	})

	return r
}

type MusicRecordingListRequest struct {
	Offset   int            `api:"offset,@query"   validate:"gte=0"`
	Limit    int            `api:"limit,@query"    validate:"gte=0"`
	ByArtist *models.IDType `api:"byArtist,@query" validate:"omitempty,gte=0"`
}

// List all of the given MusicRecording objects
func (rs MusicRecordingResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List MusicRecording")

	var req MusicRecordingListRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	musicrecordings, err := (&models.MusicRecording{}).All(rs.DB, req.ByArtist, &models.SelectQuery{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		lg.Errorf("Error in List MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewMusicRecordingListResponse(musicrecordings))
}

type MusicRecordingGetRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Get MusicRecording by id
func (rs MusicRecordingResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get MusicRecording")

	var req MusicRecordingGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicrecording, err := (&models.MusicRecording{}).Get(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Get MusicRecording (%v)", err)
		render.Render(w, r, Error404)
		return
	}

	render.JSON(w, r, musicrecording)
}

type MusicRecordingPostRequest struct {
	Body models.MusicRecording `api:"body,@body"`
}

// Post a MusicRecording
func (rs MusicRecordingResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post MusicRecording")

	var req MusicRecordingPostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicrecording, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in POST MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicrecording, err = (&models.MusicRecording{}).Get(rs.DB, musicrecording.ID)
	if err != nil {
		lg.Errorf("Error getting MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewMusicRecordingResponse(musicrecording))
}

type MusicRecordingPutRequest struct {
	ID   int64                 `api:"id,@url_param" validate:"gte=0"`
	Body models.MusicRecording `api:"body,@body"`
}

// Put MusicRecording
func (rs MusicRecordingResource) Put(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Put MusicRecording")

	var req MusicRecordingPutRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Put MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicrecording, err := req.Body.Update(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in PUT MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicrecording, err = (&models.MusicRecording{}).Get(rs.DB, musicrecording.ID)
	if err != nil {
		lg.Errorf("Error getting MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewMusicRecordingResponse(musicrecording))
}

type MusicRecordingDeleteRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Delete MusicRecording by id
func (rs MusicRecordingResource) Delete(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Delete MusicRecording")

	var req MusicRecordingDeleteRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Delete MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	err := (&models.MusicRecording{}).Delete(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Delete MusicRecording (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

// MusicRecordingRequest structure
// Any fields to be overidden go here
type MusicRecordingRequest struct {
	*models.MusicRecording
}

// Bind pre-processes any fields after a decode
func (req *MusicRecordingRequest) Bind(r *http.Request) error {
	return nil
}

// MusicRecordingResponse structure
// Add any extra fields to the response here
type MusicRecordingResponse struct {
	*models.MusicRecording
}

// NewMusicRecordingResponse creates a response with the model plus any other data
func NewMusicRecordingResponse(obj *models.MusicRecording) *MusicRecordingResponse {
	resp := &MusicRecordingResponse{MusicRecording: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *MusicRecordingResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// MusicRecordingListResponse is an array of response objects
type MusicRecordingListResponse []*MusicRecordingResponse

// NewMusicRecordingListResponse creates a response with the models plus any other data
func NewMusicRecordingListResponse(musicrecordings []*models.MusicRecording) []render.Renderer {
	list := []render.Renderer{}
	for _, musicrecording := range musicrecordings {
		list = append(list, NewMusicRecordingResponse(musicrecording))
	}
	return list
}

// MusicRecordingAuthorization is a default authorizor middleware to enforce resource access
func (rs MusicRecordingResource) MusicRecordingAuthorization(next http.Handler) http.Handler {
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

		// Selects all ethereum addresss associated with this musicgroup
		var ethereumAddresses []string
		query := `
        SELECT person.ethereum_address FROM person
        INNER JOIN musicgroup_members ON musicgroup_members.person_id = person.id
        INNER JOIN musicgroup ON musicgroup_members.musicgroup_id = musicgroup.id
        INNER JOIN musicrecording ON musicgroup.id = musicrecording.by_artist_id
        WHERE musicrecording.id = $1;
        `
		err = rs.DB.Select(&ethereumAddresses, query, id)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Check that the user making the request owns the ethereum address that the
		// musicrelease was registered with
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
