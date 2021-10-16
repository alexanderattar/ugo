package api

import (
	"net/http"
	"os"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	m "github.com/keighl/mandrill"
)

// MusicGroupResource contains REST handlers
type MusicGroupResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs MusicGroupResource) Routes() chi.Router {
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

		r.Get("/cid/{cid}", rs.Cid) // GET /{resource}/cid/{cid} - read a single resource by cid

		r.Route("/{id}", func(r chi.Router) {
			// Public routes
			r.Get("/", rs.Get) // GET /{resource}/{id} - read a single resource by id
			// Private routes
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Authenticator)
				r.Use(rs.MusicGroupAuthorization)
				r.Put("/", rs.Put)       // PUT /{resource}/{id} - update a single resource by id
				r.Delete("/", rs.Delete) // DELETE /{resource}/{id} - delete a single resource by id
			})
		})
	})

	return r
}

type MusicGroupListRequest struct {
	Offset       int            `api:"offset,@query" validate:"gte=0"`
	Limit        int            `api:"limit,@query"  validate:"gte=0"`
	EthereumAddr string         `api:"ethereumAddress,@query"`
	PersonID     *models.IDType `api:"personId,@query" validate:"omitempty,gte=0"`
}

// List all of the given MusicGroup objects
func (rs MusicGroupResource) List(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("List MusicGroup")

	var req MusicGroupListRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in List MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	musicgroups, err := (&models.MusicGroup{}).All(rs.DB, req.EthereumAddr, req.PersonID, &models.SelectQuery{Limit: req.Limit, Offset: req.Offset})
	if err != nil {
		lg.Errorf("Error in List MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	for _, mg := range musicgroups {
		mg.Members, err = mg.GetMembers(rs.DB, mg.ID)
		if err != nil {
			lg.Errorf("Error in List MusicGroup (%v)", err)
			render.Render(w, r, Error400(err))
			return
        }

        if mg.Members == nil {
            return
        }

		// TODO - This should be abstracted to a read authorization middleware 
		// Omit email fields for user privacy. Ideally we'd like to flip this logic, making it omit by
		// default, and adding it in when we want. 
        for _, m := range mg.Members {
            m.Email = nil
        }
	}

	render.RenderList(w, r, NewMusicGroupListResponse(musicgroups))
}

type MusicGroupGetRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Get MusicGroup by id
func (rs MusicGroupResource) Get(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get MusicGroup")

	var req MusicGroupGetRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicgroup, err := (&models.MusicGroup{}).Get(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Get MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
    }

	// TODO - This should be abstracted to a read authorization middleware 
    // Omit emails if logged in user isn't in group
    if musicgroup.Members != nil {
		_, claims, _ := jwtauth.FromContext(r.Context())
		
		// Check if the current user is in the group
        var isInGroup bool
        for _, m := range musicgroup.Members {
            if claims["ethereumAddress"] == m.EthereumAddress {
                isInGroup = true
            }
        }

		// If not, all the member emails should be omitted 
        if !isInGroup {
            for _, m := range musicgroup.Members {
                m.Email = nil
            }
        }
    }

	render.JSON(w, r, musicgroup)
}

type MusicGroupGetByCIDRequest struct {
	CID string `api:"cid,@url_param" validate:"gte=0"`
}

// Cid MusicGroup by id
func (rs MusicGroupResource) Cid(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get MusicGroup by CID")

	var req MusicGroupGetByCIDRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Get MusicGroup by CID (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicgroup, err := (&models.MusicGroup{}).GetByCID(rs.DB, req.CID)
	if err != nil {
		lg.Errorf("Error in Get MusicGroup by CID (%v)", err)
		render.Render(w, r, Error404)
		return
	}

	render.JSON(w, r, musicgroup)
}

type MusicGroupPostRequest struct {
	Body models.MusicGroup `api:"body,@body"`
}

// Post a MusicGroup
func (rs MusicGroupResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post MusicGroup")

	var req MusicGroupPostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Post MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicgroup, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("Error in Post MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicgroup, err = (&models.MusicGroup{}).Get(rs.DB, musicgroup.ID)
	if err != nil {
		lg.Errorf("Error getting MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewMusicGroupResponse(musicgroup))

	if os.Getenv("TEST_ENV") == "Yes" {
		return
	}

	sendRegisteredEmail(musicgroup, "artist-registered")
}

type MusicGroupPutRequest struct {
	ID   int64             `api:"id,@url_param" validate:"gte=0"`
	Body models.MusicGroup `api:"body,@body"`
}

// Put MusicGroup
func (rs MusicGroupResource) Put(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Put MusicGroup")

	var req MusicGroupPutRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Put MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	musicgroup, err := req.Body.Update(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Put MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Populate sub-object fields
	musicgroup, err = (&models.MusicGroup{}).Get(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Getting MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewMusicGroupResponse(musicgroup))
}

type MusicGroupDeleteRequest struct {
	ID int64 `api:"id,@url_param" validate:"gte=0"`
}

// Delete MusicGroup by id
func (rs MusicGroupResource) Delete(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Delete MusicGroup")

	var req MusicGroupDeleteRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("Error in Delete MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	err := (&models.MusicGroup{}).Delete(rs.DB, req.ID)
	if err != nil {
		lg.Errorf("Error in Delete MusicGroup (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "DELETE")
}

// MusicGroupRequest structure
// Any fields to be overidden go here
type MusicGroupRequest struct {
	*models.MusicGroup
}

// Bind pre-processes any fields after a decode
func (req *MusicGroupRequest) Bind(r *http.Request) error {
	return nil
}

// MusicGroupResponse structure
// Add any extra fields to the response here
type MusicGroupResponse struct {
	*models.MusicGroup
}

// NewMusicGroupResponse creates a response with the model plus any other data
func NewMusicGroupResponse(obj *models.MusicGroup) *MusicGroupResponse {
	resp := &MusicGroupResponse{MusicGroup: obj}
	return resp
}

// Render post-processes the data before a response is returned
func (resp *MusicGroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// MusicGroupListResponse is an array of response objects
type MusicGroupListResponse []*MusicGroupResponse

// NewMusicGroupListResponse creates a response with the models plus any other data
func NewMusicGroupListResponse(musicgroups []*models.MusicGroup) []render.Renderer {
	list := []render.Renderer{}
	for _, musicgroup := range musicgroups {
		list = append(list, NewMusicGroupResponse(musicgroup))
	}
	return list
}

// MusicGroupAuthorization is a default authorizor middleware to enforce resource access
func (rs MusicGroupResource) MusicGroupAuthorization(next http.Handler) http.Handler {
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
        WHERE musicgroup.id = $1;
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

func sendMusicGroupsEmail(nameEmail map[string]string, mergeVar map[string]interface{}, templateName string) {
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

func sendRegisteredEmail(musicgroup *models.MusicGroup, templateName string) {
	mergeVar := make(map[string]interface{})
	nameEmailMap := make(map[string]string)

	mergeVar["firstName"] = nil
	mergeVar["groupName"] = musicgroup.Name
	mergeVar["groupDescription"] = musicgroup.Description
	mergeVar["groupImageSource"] = *musicgroup.Image.ContentURL

	var email string
	var name string

	for _, member := range musicgroup.Members {
		if member.Email != nil {
			if member.GivenName != nil {
				email = *member.Email
				name = *member.GivenName
				nameEmailMap[name] = email
			}
		}
	}

	sendMusicGroupsEmail(nameEmailMap, mergeVar, templateName)
}
