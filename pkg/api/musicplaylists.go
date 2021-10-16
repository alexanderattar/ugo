package api

import (
	"net/http"
	"strconv"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// MusicPlaylistResource contains REST handlers
type MusicPlaylistResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs MusicPlaylistResource) Routes() chi.Router {
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
				r.Use(rs.MusicPlaylistAuthorization)
				r.Put("/", rs.Put)                           // PUT /{resource}/{id} - update a single resource by id
				r.Delete("/", rs.Delete)                     // DELETE /{resource}/{id} - delete a single resource by id
				r.Post("/track/{trackId}", rs.AddTrack)      // POST /{resource}/{id}/track/{trackId} - add track to playlist
				r.Delete("/track/{trackId}", rs.DeleteTrack) // DELETE /{resource}/{id}/track/{trackId} - delete track from playlist
			})
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(AdminAuthorization)
		})
	})

	return r
}

// Post a MusicPlaylist
func (rs MusicPlaylistResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post MusicPlaylist")

	var req MusicPlaylistPostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicplaylistID, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in POST MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicplaylist, err := (&models.MusicPlaylist{}).Get(rs.DB, musicplaylistID)
	if err != nil {
		lg.Errorf("Error getting MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewMusicPlaylistResponse(musicplaylist))
}

type MusicPlaylistPutRequest struct {
	ID   int64               `api:"id,@url_param" validate:"gte=0"`
	Body models.MusicPlaylist `api:"body,@body"`
}

// Put MusicPlaylist
func (rs MusicPlaylistResource) Put(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Put MusicPlaylist")

	var req MusicPlaylistPutRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in PUT MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}
	req.Body.ID = req.ID

	_, err := req.Body.Update(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in PUT MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicplaylist, err := (&models.MusicPlaylist{}).Get(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error getting MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewMusicPlaylistResponse(musicplaylist))
}

type MusicPlaylistRequest struct {
	PersonId     *models.IDType `api:"personId,@query"   validate:"omitempty,gte=0"`
	Limit      int            `api:"limit,@query"    validate:"gte=0"`
	Offset     int            `api:"offset,@query"   validate:"gte=0"`
}

// List all of the given MusicPlaylist objects
func (rs MusicPlaylistResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List MusicPlaylist")

	var req MusicPlaylistRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	musicplaylists, err := (&models.MusicPlaylist{}).All(rs.DB, req.PersonId, &models.SelectQuery{
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		lg.Errorf("Error in List MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewMusicPlaylistListResponse(musicplaylists))
}

type MusicPlaylistDeleteRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Delete MusicPlaylist by id
func (rs MusicPlaylistResource) Delete(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Delete MusicPlaylist")

	var req MusicPlaylistDeleteRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in DELETE MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	err := (&models.MusicPlaylist{}).Delete(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Delete MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

// AddTrack to MusicPlaylist
func (rs MusicPlaylistResource) AddTrack(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("AddTrack to MusicPlaylist")
	id := chi.URLParam(r, "id")
	trackID := chi.URLParam(r, "trackId")

	idInt, err := strconv.ParseInt(id, 10, 64)
	obj := &models.MusicPlaylist{}
	musicplaylist, err := obj.Get(rs.DB, idInt)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	err = obj.AddTrack(rs.DB, musicplaylist.ID, trackID)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "INSERT")
}

// DeleteTrack from MusicPlaylist
func (rs MusicPlaylistResource) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("DeleteTrack from MusicPlaylist")
	id := chi.URLParam(r, "id")
	trackID := chi.URLParam(r, "trackId")

	idInt, err := strconv.ParseInt(id, 10, 64)
	obj := &models.MusicPlaylist{}
	musicplaylist, err := obj.Get(rs.DB, idInt)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	err = obj.DeleteTrack(rs.DB, musicplaylist.ID, trackID)

	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

type MusicPlaylistPostRequest struct {
	Body models.MusicPlaylist `api:"body,@body"`
}

// Bind pre-processes any fields after a decode
func (req *MusicPlaylistRequest) Bind(r *http.Request) error {
	return nil
}

// MusicPlaylistResponse structure
// Add any extra fields to the response here
type MusicPlaylistResponse struct {
	*models.MusicPlaylist
}

// NewMusicPlaylistResponse creates a response with the model plus any other data
func NewMusicPlaylistResponse(obj *models.MusicPlaylist) *MusicPlaylistResponse {
	resp := &MusicPlaylistResponse{MusicPlaylist: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *MusicPlaylistResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// MusicPlaylistListResponse is an array of response objects
type MusicPlaylistListResponse []*MusicPlaylistResponse

// NewMusicPlaylistListResponse creates a ListResponse
func NewMusicPlaylistListResponse(musicplaylists []*models.MusicPlaylist) []render.Renderer {
	list := []render.Renderer{}
	for _, musicplaylist := range musicplaylists {
		list = append(list, NewMusicPlaylistResponse(musicplaylist))
	}
	return list
}

// Render post-processes the data before a response is returned
func (resp *MusicPlaylistTrackResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// MusicPlaylistTrackResponse is an array of response objects
type MusicPlaylistTrackResponse string

type MusicPlaylistGetRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Get MusicPlaylist by id
func (rs MusicPlaylistResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get MusicPlaylist")

	var req MusicPlaylistGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get MusicPlaylist (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicplaylist, err := (&models.MusicPlaylist{}).Get(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Get MusicPlaylist (%v)", err)
		render.Render(w, r, Error404)
		return
	}
	render.JSON(w, r, musicplaylist)
}

// PersonAuthorization is a default authorizor middleware to enforce resource access
func (rs MusicPlaylistResource) MusicPlaylistAuthorization(next http.Handler) http.Handler {
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

		var userID int
		query := `
        SELECT musicplaylist.by_user_id FROM musicplaylist
        WHERE musicplaylist.id = $1;
        `

		err = rs.DB.Get(&userID, query, id)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		var ethereumAddress string
		query = `
        SELECT person.ethereum_address FROM person
        WHERE person.id = $1;
        `

		err = rs.DB.Get(&ethereumAddress, query, userID)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Check that the user making the request owns the playlist
		authorized := false
		if ethereumAddress == claims["ethereumAddress"] {
			authorized = true
		}

		if authorized != true {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authorized, pass it through
		next.ServeHTTP(w, r)
	})
}
