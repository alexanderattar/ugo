package api

import (
	"net/http"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// PurchaseResource contains REST handlers
type PurchaseResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs PurchaseResource) Routes() chi.Router {
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

			// Private routes
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Authenticator)
				r.Delete("/", rs.Delete) // DELETE /{resource}/{id} - delete a single resource by id
			})
		})
	})

	return r
}

type PurchaseListRequest struct {
	Limit        int            `api:"limit,@query"     validate:"gte=0"`
	Offset       int            `api:"offset,@query"    validate:"gte=0"`
	ReleaseID    *models.IDType `api:"releaseId,@query" validate:"omitempty,gte=0"`
	EthereumAddr string         `api:"ethereumAddress,@query"`
}

// List all of the given Purchase objects
// @@TODO: pagination
func (rs PurchaseResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List Purchase")

	var req PurchaseListRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List Purchase (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	purchases, err := (&models.Purchase{}).All(rs.DB, req.ReleaseID, req.EthereumAddr, &models.SelectQuery{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewPurchaseListResponse(purchases))
}

type PurchaseGetRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Get Purchase by id
func (rs PurchaseResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get Purchase")

	var req PurchaseGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get Purchase (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	purchase, err := (&models.Purchase{}).Get(rs.DB, req.ID)
	if err != nil {
		render.Render(w, r, Error404)
		return
	}

	render.JSON(w, r, purchase)
}

type PurchasePostRequest struct {
	Body models.Purchase `api:"body,@body"`
}

// Post a Purchase
func (rs PurchaseResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post Purchase")

	var req PurchasePostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post Purchase (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	purchase, err := req.Body.Create(rs.DB)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewPurchaseResponse(purchase))
}

type PurchaseDeleteRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Delete Purchase by id
func (rs PurchaseResource) Delete(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Delete Purchase")

	var req PurchaseDeleteRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Delete Purchase (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	err := (&models.Purchase{}).Delete(rs.DB, req.ID)
	if err != nil {
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

// PurchaseRequest structure
// Any fields to be overidden go here
type PurchaseRequest struct {
	*models.Purchase
}

// Bind pre-processes any fields after a decode
func (req *PurchaseRequest) Bind(r *http.Request) error {
	return nil
}

// PurchaseResponse structure
// Add any extra fields to the response here
type PurchaseResponse struct {
	*models.Purchase
}

// NewPurchaseResponse creates a response with the model plus any other data
func NewPurchaseResponse(obj *models.Purchase) *PurchaseResponse {
	resp := &PurchaseResponse{Purchase: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *PurchaseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PurchaseListResponse is an array of response objects
type PurchaseListResponse []*PurchaseResponse

// NewPurchaseListResponse creates a ListResponse
func NewPurchaseListResponse(purchases []*models.Purchase) []render.Renderer {
	list := []render.Renderer{}
	for _, purchase := range purchases {
		list = append(list, NewPurchaseResponse(purchase))
	}
	return list
}
