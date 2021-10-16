package api

import (
	"net/http"
	"os"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"
	m "github.com/keighl/mandrill"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// ReportResource contains REST handlers
type ReportResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs ReportResource) Routes() chi.Router {
	r := chi.NewRouter()
	// Seek, verify and validate JWT tokens
	r.Use(jwtauth.Verifier(TokenAuth))

	// Protected routes
	r.Group(func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/", rs.Post) // POST /{resource} - create a new resource and persist it
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(AdminAuthorization)
			r.Get("/", rs.List)
		})

		r.Route("/{id}", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(AdminAuthorization)
				r.Get("/", rs.Get)                  // GET /{resource}/{id} - read a single resource by id
				r.Put("/", rs.Put)                  // PUT /{resource}/{id} - update single resource resource and persist it
				r.Put("/resolve", rs.Resolve)       // PUT /{resource}/{id}/resolve - update single resource resource and persist it
				r.Put("/deactivate", rs.Deactivate) // PUT /{resource}/{id}/deactivate - update single resource resource and persist it
			})
		})
	})

	return r
}

type ReportPostRequest struct {
	Body models.Report `api:"body,@body"`
}

// Post a Report
func (rs ReportResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post Report")

	var req ReportPostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	reportID, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in POST Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Reset the object with a Get so the id fields get populated
	report, err := (&models.Report{}).Get(rs.DB, reportID)
	if err != nil {
		lg.Errorf("Error in POST Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewReportResponse(report))

	if os.Getenv("TEST_ENV") == "Yes" {
		return
	}

	sendFlaggedEmail(report, "flagged-content")
}

type ReportListRequest struct {
	Limit     int    `api:"limit,@query"          validate:"gte=0"`
	Offset    int    `api:"offset,@query"         validate:"gte=0"`
	ReleaseID *int64 `api:"musicreleaseID,@query" validate:"omitempty,gte=0"`
}

// List all of the given Report objects
func (rs ReportResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List Report")

	var req ReportListRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	reports, err := (&models.Report{}).All(rs.DB, req.ReleaseID, &models.SelectQuery{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		lg.Errorf("Error in List Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.RenderList(w, r, NewReportListResponse(reports))
}

type ReportGetRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Get Report by id
func (rs ReportResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get Report")

	var req ReportGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	report, err := (&models.Report{}).Get(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Get Report (%v)", err)
		render.Render(w, r, Error404)
		return
	}

	render.JSON(w, r, report)
}

type ReportPutRequest struct {
	ID   int64         `api:"id,@url_param" validate:"gte=0"`
	Body models.Report `api:"body,@body"`
}

// Put Report
func (rs ReportResource) Put(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Put Report")

	var req ReportPutRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Put Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	updated, err := req.Body.Update(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in PUT Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	updated.ID = req.ID
	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewReportResponse(updated))
}

type ReportResolveRequest struct {
	ID   int64         `api:"id,@url_param" validate:"gte=0"`
	Body models.Report `api:"body,@body"`
}

// Resolve Report
func (rs ReportResource) Resolve(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Resolve Report")

	var req ReportResolveRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Resolve Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	report, err := req.Body.Resolve(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in PUT to resolve Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	report.ID = req.ID
	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewReportResponse(report))
}

type ReportDeactivateRequest struct {
	ID   int64         `api:"id,@url_param" validate:"gte=0"`
	Body models.Report `api:"body,@body"`
}

// Deactivate Report
func (rs ReportResource) Deactivate(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Deactivate Report")

	var req ReportDeactivateRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Deactivate Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	report, err := req.Body.Deactivate(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in PUT to deactivate Report (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	report.ID = req.ID
	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewReportResponse(report))

	if os.Getenv("TEST_ENV") == "Yes" {
		return
	}

	sendRemovedEmail(report, "release-removed")
}

// ReportRequest structure
// Any fields to be overidden go here
type ReportRequest struct {
	*models.Report
}

// Bind pre-processes any fields after a decode
func (req *ReportRequest) Bind(r *http.Request) error {
	return nil
}

// ReportResponse structure
// Add any extra fields to the response here
type ReportResponse struct {
	*models.Report
}

// NewReportResponse creates a response with the model plus any other data
func NewReportResponse(obj *models.Report) *ReportResponse {
	resp := &ReportResponse{Report: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *ReportResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ReportListResponse is an array of response objects
type ReportListResponse []*ReportResponse

// NewReportListResponse creates a ListResponse
func NewReportListResponse(reports []*models.Report) []render.Renderer {
	list := []render.Renderer{}
	for _, report := range reports {
		list = append(list, NewReportResponse(report))
	}
	return list
}

func sendFlaggedEmail(report *models.Report, templateName string) {
	mergeVar := make(map[string]interface{})
	nameEmailMap := make(map[string]string)

	mergeVar["firstName"] = nil
	mergeVar["releaseName"] = report.MusicRelease.ReleaseOf.Name
	mergeVar["groupName"] = report.MusicRelease.ReleaseOf.ByArtist.Name
	mergeVar["flaggingReason"] = report.Reason
	mergeVar["flaggingMessage"] = report.Message
	mergeVar["albumCover"] = *report.MusicRelease.Image.ContentURL

	var email string
	var name string

	for _, member := range report.MusicRelease.ReleaseOf.ByArtist.Members {
		email = *member.Email
		name = *member.GivenName
		nameEmailMap[name] = email
	}

	sendReportsEmail(nameEmailMap, mergeVar, templateName)
}

func sendRemovedEmail(report *models.Report, templateName string) {
	mergeVar := make(map[string]interface{})
	nameEmailMap := make(map[string]string)

	mergeVar["firstName"] = nil
	mergeVar["releaseName"] = report.MusicRelease.ReleaseOf.Name
	mergeVar["groupName"] = report.MusicRelease.ReleaseOf.ByArtist.Name
	mergeVar["removalReason"] = report.Reason
	mergeVar["removalMessage"] = report.Response
	mergeVar["albumImageSource"] = *report.MusicRelease.Image.ContentURL

	var email string
	var name string

	for _, member := range report.MusicRelease.ReleaseOf.ByArtist.Members {
		email = *member.Email
		name = *member.GivenName
		nameEmailMap[name] = email
	}

	sendReportsEmail(nameEmailMap, mergeVar, templateName)
}

func sendReportsEmail(nameEmail map[string]string, mergeVar map[string]interface{}, templateName string) {
	failed := false
	mandrillAPIKey := os.Getenv("MANDRILL_API_KEY")
	client := m.ClientWithKey(mandrillAPIKey)

	message := &m.Message{}

	for name, email := range nameEmail {
		message.AddRecipient(email, name, "to")
		message.FromEmail = "no-reply@ujomusic.com"
		message.MergeLanguage = "handlebars"
		mergeVar["firstName"] = name

		merge := m.MapToRecipientVars(email, mergeVar)

		message.MergeVars = []*m.RcptMergeVars{merge}

		responses, err := client.MessagesSendTemplate(message, templateName, nil)

		for _, r := range responses {
			if r.Status != "send" {
				failed = true
				break
			}
		}

		if err != nil || failed != false {
			lg.Errorf("Error in sending email (%v)", err)
			return
		}
	}
}
