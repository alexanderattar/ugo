package api

import (
	"database/sql"
	"net/http"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi"
	// "github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// PayEventResource contains REST handlers
type PayEventResource struct {
	DB *sqlx.DB
}

// Routes creates a REST router for the Resource
func (rs PayEventResource) Routes() chi.Router {
	r := chi.NewRouter()
	// Seek, verify and validate JWT tokens
	r.Use(jwtauth.Verifier(TokenAuth))

	// Protected routes
	r.Group(func(r chi.Router) {
		// Handle valid / invalid tokens
		// This can be modified using the Authenticator method in auth.go
		r.Use(jwtauth.Authenticator)
		r.Get("/unpaid", rs.GetUnpaidForUser)        // GET /unpaid
		r.Post("/mark-as-paid", rs.MarkAsPaid)       // POST /mark-as-paid - mark PayEvents as paid
		r.Post("/get-unused-link", rs.GetUnusedLink) // POST /get-unused-link - grab one unused PayEvent link from the DB (and mark it as paid)
		r.Post("/", rs.Post)                         // POST /{resource} - create a new resource and persist it
	})

	return r
}

type PayEventPostRequest struct {
	NoBeneficiary bool            `api:"no_beneficiary,@query"`
	Body          models.PayEvent `api:"body,@body"`
}

// Post a PayEvent
func (rs PayEventResource) Post(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post PayEvent")

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
	var req PayEventPostRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("error decoding PayEvent.Post request (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	//
	// Create the data
	//
	if !req.NoBeneficiary {
		req.Body.PlayedByID = &user.ID
	}

	payeventID, err := req.Body.Create(rs.DB)
	if err != nil {
		lg.Errorf("error in POST PayEvent (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	payevent, err := (&models.PayEvent{}).Get(rs.DB, payeventID)
	if err != nil {
		lg.Errorf("error getting PayEvent (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, payevent)
}

type PayEventGetUnpaidForUserRequest struct {
	NoBeneficiary bool `api:"no_beneficiary,@query"`
	Limit         int  `api:"limit,@query"    validate:"gte=0"`
	Offset        int  `api:"offset,@query"   validate:"gte=0"`
}

func (rs PayEventResource) GetUnpaidForUser(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Get PayEvent GetUnpaidForUser")

	//
	// Decode the request body
	//
	var req PayEventGetUnpaidForUserRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("error decoding PayEvent.GetUnpaidForUser request (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

	var userID *models.IDType
	if !req.NoBeneficiary {
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
		userID = &user.ID
	}

	//
	// Fetch the data
	//
	unpaid, err := (&models.PayEvent{}).GetUnpaidForUser(rs.DB, userID, &models.SelectQuery{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, r, unpaid)
}

type PayEventMarkAsPaidRequest struct {
	IDs []models.IDType `api:"body,@body"`
}

// Mark a set of PayEvents as paid
func (rs PayEventResource) MarkAsPaid(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post PayEvent MarkAsPaid")

	//
	// Decode the request body
	//
	var req PayEventMarkAsPaidRequest
	if err := DecodeRequest(r, &req); err != nil {
		lg.Errorf("error decoding PayEvent.MarkAsPaid request (%v)", err)
		render.Render(w, r, Error400(err))
		return
	}

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
	// Update the rows
	//
	err = (&models.PayEvent{}).MarkAsPaid(rs.DB, req.IDs, user.ID)
	if err != nil {
		lg.Errorf("error marking PayEvents as paid (%v)", err)
		render.Render(w, r, Error500(err))
		return
	}

	// Return empty
	render.JSON(w, r, map[string]interface{}{})

}

type PayEventGetUnusedLinkRequest struct {
}

// In this situation, where AddBeneficiary was specified, we're taking a generic payment link
// that hasn't been assigned to anybody and are assigning it to a user and marking it paid
// all at once.  At present, each user is only allowed a single payment link, so if we see
// that they've
func (rs PayEventResource) GetUnusedLink(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Post PayEvent GetUnusedLink")

	//
	// Get the user's ID from the JWT
	//
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 401)
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
		http.Error(w, err.Error(), 403)
		return
	}

	//
	// Check to see if they've done this before.  If so, disallow it.
	//
	paid, err := (&models.PayEvent{}).GetPaidForUser(rs.DB, user.ID, &models.SelectQuery{})
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), 500)
		return
	} else if len(paid) > 0 {
		http.Error(w, "Cannot redeem more than one prepaid link.", 400)
		return
	}

	//
	// Find an unused link
	//
	var nobody *models.IDType
	unpaid, err := (&models.PayEvent{}).GetUnpaidForUser(rs.DB, nobody, &models.SelectQuery{Limit: 1})
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), 500)
		return
	} else if len(unpaid) == 0 {
		http.Error(w, "no unclaimed links remaining", 500)
		return
	}

	//
	// Update the rows
	//
	err = (&models.PayEvent{}).MarkPrepaidAsClaimed(rs.DB, unpaid[0].ID, user.ID)
	if err != nil {
		lg.Errorf("error marking PayEvents as paid (%v)", err)
		render.Render(w, r, Error500(err))
		return
	}

	// Return empty
	render.JSON(w, r, unpaid[0])

}
