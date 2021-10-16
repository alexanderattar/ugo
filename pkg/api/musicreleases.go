package api

import (
	"net/http"
	"strconv"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// MusicReleaseResource contains REST handlers
type MusicReleaseResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs MusicReleaseResource) Routes() chi.Router {
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
			r.Post("/link", rs.Link)
		})

		r.Route("/{id}", func(r chi.Router) {
			// Public routes
			r.Get("/", rs.Get) // GET /{resource}/{id} - read a single resource by id

			// Private routes
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Authenticator)
				r.Use(rs.MusicReleaseAuthorization)
				r.Put("/", rs.Put)                           // PUT /{resource}/{id} - update a single resource by id
				r.Delete("/", rs.Delete)                     // DELETE /{resource}/{id} - delete a single resource by id
				r.Post("/track/{trackId}", rs.AddTrack)      // POST /{resource}/{id}/track/{trackId} - add track to release
				r.Delete("/track/{trackId}", rs.DeleteTrack) // DELETE /{resource}/{id}/track/{trackId} - delete track from release
			})
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(AdminAuthorization)
			r.Get("/inactive", rs.ListInactive) // GET /inactive - get inactive musicreleases
		})
	})

	return r
}

type MusicReleaseListRequest struct {
	ByArtist   *models.IDType `api:"byArtist,@query" validate:"omitempty,gte=0"`
	Limit      int            `api:"limit,@query"    validate:"gte=0"`
	Offset     int            `api:"offset,@query"   validate:"gte=0"`
	OrderBy    string         `api:"orderBy,@query"  validate:"isdefault|oneof=createdAt datePublished releaseOf.name releaseOf.byArtist.name"`
	Descending bool           `api:"descending,@query"`
}

// List all of the given MusicRelease objects
func (rs MusicReleaseResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List MusicRelease")

	var req MusicReleaseListRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	musicreleases, err := (&models.MusicRelease{}).All(rs.DB, req.ByArtist, &models.SelectQuery{
		Limit:      req.Limit,
		Offset:     req.Offset,
		OrderBy:    req.OrderBy,
		Descending: req.Descending,
	})
	if err != nil {
		lg.Errorf("Error in List MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewMusicReleaseListResponse(musicreleases))
}

type MusicReleaseListInactiveRequest struct {
	Limit  int `api:"limit,@query"  validate:"gte=0"`
	Offset int `api:"offset,@query" validate:"gte=0"`
}

// ListInactive lists all inactive releases
func (rs MusicReleaseResource) ListInactive(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List inactive MusicReleases")

	var req MusicReleaseListInactiveRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List inactive MusicReleases (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicreleases, err := (&models.MusicRelease{}).AllInactive(rs.DB, &models.SelectQuery{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		lg.Errorf("Error in List inactive MusicReleases (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewMusicReleaseListResponse(musicreleases))
}

type MusicReleaseGetRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Get MusicRelease by id
func (rs MusicReleaseResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get MusicRelease")

	var req MusicReleaseGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicrelease, err := (&models.MusicRelease{}).GetByID(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Get MusicRelease (%v)", err)
		render.Render(w, r, Error404)
		return
	}
	if !musicrelease.Active {
		render.Render(w, r, Error400(errors.New("MusicRelease has been removed")))
		return
	}

	render.JSON(w, r, musicrelease)
}

type MusicReleasePostRequest struct {
	Body models.MusicRelease `api:"body,@body"`
}

// Post a MusicRelease
func (rs MusicReleaseResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post MusicRelease")

	var req MusicReleasePostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicrelease, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in POST MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicrelease, err = (&models.MusicRelease{}).GetByID(rs.DB, musicrelease.ID)
	if err != nil {
		lg.Errorf("Error getting MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewMusicReleaseResponse(musicrelease))
}

type MusicReleaseLinkRequest struct {
	Body models.MusicRelease `api:"body,@body"`
}

// Post a MusicRelease
func (rs MusicReleaseResource) Link(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post MusicRelease")

	var req MusicReleaseLinkRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicrelease, err := req.Body.Link(rs.DB)

	if err != nil {
		lg.Errorf("Error in POST MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewMusicReleaseResponse(musicrelease))
}

type MusicReleasePutRequest struct {
	ID   int64               `api:"id,@url_param" validate:"gte=0"`
	Body models.MusicRelease `api:"body,@body"`
}

// Put MusicRelease
func (rs MusicReleaseResource) Put(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Put MusicRelease")

	var req MusicReleasePutRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in PUT MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}
	req.Body.ID = req.ID

	musicrelease, err := req.Body.Update(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in PUT MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicrelease, err = (&models.MusicRelease{}).GetByID(rs.DB, musicrelease.ID)
	if err != nil {
		lg.Errorf("Error getting MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewMusicReleaseResponse(musicrelease))
}

type MusicReleaseDeleteRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Delete MusicRelease by id
func (rs MusicReleaseResource) Delete(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Delete MusicRelease")

	var req MusicReleaseDeleteRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in DELETE MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	err := (&models.MusicRelease{}).Delete(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Delete MusicRelease (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

// AddTrack to MusicRelease
func (rs MusicReleaseResource) AddTrack(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("AddTrack to MusicRelease")
	id := chi.URLParam(r, "id")
	trackID := chi.URLParam(r, "trackId")

	idInt, err := strconv.ParseInt(id, 10, 64)
	obj := &models.MusicRelease{}
	musicrelease, err := obj.GetByID(rs.DB, idInt)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	err = obj.AddTrack(rs.DB, musicrelease.ReleaseOfID, trackID)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "INSERT")
}

// DeleteTrack from MusicRelease
func (rs MusicReleaseResource) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("DeleteTrack from MusicRelease")
	id := chi.URLParam(r, "id")
	trackID := chi.URLParam(r, "trackId")

	idInt, err := strconv.ParseInt(id, 10, 64)
	obj := &models.MusicRelease{}
	musicrelease, err := obj.GetByID(rs.DB, idInt)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	err = obj.DeleteTrack(rs.DB, musicrelease.ReleaseOfID, trackID)

	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

// MusicReleaseRequest structure
// Any fields to be overidden go here
type MusicReleaseRequest struct {
	*models.MusicRelease
}

// Bind pre-processes any fields after a decode
func (req *MusicReleaseRequest) Bind(r *http.Request) error {
	return nil
}

// MusicReleaseResponse structure
// Add any extra fields to the response here
type MusicReleaseResponse struct {
	*models.MusicRelease
}

// NewMusicReleaseResponse creates a response with the model plus any other data
func NewMusicReleaseResponse(obj *models.MusicRelease) *MusicReleaseResponse {
	resp := &MusicReleaseResponse{MusicRelease: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *MusicReleaseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// MusicReleaseListResponse is an array of response objects
type MusicReleaseListResponse []*MusicReleaseResponse

// NewMusicReleaseListResponse creates a ListResponse
func NewMusicReleaseListResponse(musicreleases []*models.MusicRelease) []render.Renderer {
	list := []render.Renderer{}
	for _, musicrelease := range musicreleases {
		list = append(list, NewMusicReleaseResponse(musicrelease))
	}
	return list
}

// Render post-processes the data before a response is returned
func (resp *MusicReleaseTrackResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// MusicReleaseTrackResponse is an array of response objects
type MusicReleaseTrackResponse string

// MusicReleaseAuthorization is a default authorizor middleware to enforce resource access
func (rs MusicReleaseResource) MusicReleaseAuthorization(next http.Handler) http.Handler {
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
        INNER JOIN musicalbum ON musicgroup.id = musicalbum.by_artist_id
        INNER JOIN musicrelease ON musicalbum.id = musicrelease.release_of_id
        WHERE musicrelease.id = $1;
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

		if claims["admin"] == true {
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
