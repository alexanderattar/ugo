package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

// TokenAuth initializes jwtauth
var TokenAuth *jwtauth.JWTAuth

func init() {
	secret := os.Getenv("UJO_API_SECRET")
	if secret == "" {
		panic("API secret has not been set")
	}

	// Authentication initializes with empty claims
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

// ValidateSignature decodes a signature and recovers the public key / address
func ValidateSignature(message string, signature string) common.Address {
	sig := hexutil.MustDecode(signature)

	if len(sig) != 65 {
		panic(fmt.Errorf("signature must be 65 bytes long"))
	}

	if sig[64] != 27 && sig[64] != 28 {
		panic(fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)"))
	}
	sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	// Hack that to enable a password
	// TODO - Think about a better way to abstract this
	if message == "Admin" {
		message = os.Getenv("UJO_ADMIN_SECRET")
	}

	// Format the ethereum signed message
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	data := crypto.Keccak256([]byte(msg))

	rpk, err := crypto.Ecrecover(data, sig)
	if err != nil {
		panic(err)
	}

	pubKey := crypto.ToECDSAPub(rpk)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr
}

// AuthHandler contains REST handlers
type AuthHandler struct {
	DB *sqlx.DB
}

// Auth Handler authenticates a user via an ethereum address
func (h AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Authenticating")
	if r.Method != "POST" {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}

	data := &AuthRequest{}
	if err := render.Bind(r, data); err != nil {
		lg.Errorf("Error reading request body %v", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Validate Ethereum Signature
	addr := common.HexToAddress(data.EthereumAddress)
	recoveredAddr := ValidateSignature(data.Message, data.Signature)
	if addr != recoveredAddr {
		lg.Errorf("Address mismatch: want: %x have: %x", addr, recoveredAddr)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	_, tokenString, _ := TokenAuth.Encode(
		jwtauth.Claims{"ethereumAddress": data.EthereumAddress, "exp": time.Now().Add(time.Hour * 24).Unix()},
	)

	lg.Infoln("JWT: ", tokenString)

	// Get the person object
	person, err := (&models.Person{}).Get(h.DB, -1, data.EthereumAddress)
	if err != nil {
		// Bypass the this error on initial registration
		// TODO - Refactor the frontend and backend auth interaction to avoid this hack
		if err.Error() == "sql: no rows in result set" {
			err = nil
		} else {
			lg.Errorf("Error in Get Person (%v)", err)
			render.Render(w, r, Error400(err))
			return
		}
	}

	musicgroups := []*models.MusicGroup{}
	if person != nil {
		// Get the musicgroup object
		musicgroups, err = (&models.MusicGroup{}).AllByPersonID(h.DB, person.ID, &models.SelectQuery{Limit: 0})
		for _, mg := range musicgroups {
			mg.Members, err = mg.GetMembers(h.DB, mg.ID)
		}

		if err != nil {
			lg.Errorf("Error in List MusicGroup (%v)", err)
			render.Render(w, r, Error400(err))
			return
		}
	}

	render.Render(
		w, r, NewAuthResponse(&JWT{Token: tokenString}, person, musicgroups),
	)
}

// AdminAuth Handler authenticates a user via an ethereum address
func (h AuthHandler) AdminAuth(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Authenticating")
	if r.Method != "POST" {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}

	data := &AdminAuthRequest{}
	if err := render.Bind(r, data); err != nil {
		lg.Errorf("Error reading request body %v", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Validate Ethereum Signature
	addr := common.HexToAddress(data.EthereumAddress)
	recoveredAddr := ValidateSignature("Admin", data.Signature)
	if addr != recoveredAddr {
		lg.Errorf("Address mismatch: want: %x have: %x", addr, recoveredAddr)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	_, tokenString, _ := TokenAuth.Encode(
		jwtauth.Claims{"ethereumAddress": data.EthereumAddress, "admin": true, "exp": time.Now().Add(time.Hour * 24).Unix()},
	)

	lg.Infoln("JWT: ", tokenString)

	render.Render(
		w, r, NewAdminAuthResponse(&JWT{Token: tokenString}),
	)
}

// JWT model
type JWT struct {
	Token string `json:"token"`
}

// AuthRequest structure
type AuthRequest struct {
	EthereumAddress string `json:"ethereumAddress"`
	Signature       string `json:"signature"`
	Message         string `json:"message"`
}

// AuthResponse structure
// Add any extra fields to the response here
type AuthResponse struct {
	JWT         *JWT                 `json:"jwt"`
	Person      *models.Person       `json:"person"`
	MusicGroups []*models.MusicGroup `json:"musicgroups"`
}

// NewAuthResponse creates a response with the model plus any other data
func NewAuthResponse(obj *JWT, person *models.Person, musicgroups []*models.MusicGroup) *AuthResponse {
	resp := &AuthResponse{JWT: obj, Person: person, MusicGroups: musicgroups}
	return resp
}

// Bind pre-processes any fields after a decode
func (req *AuthRequest) Bind(r *http.Request) error {
	return nil
}

// Render post-processes the data before a response is returned
func (resp *AuthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// AdminAuthRequest structure
type AdminAuthRequest struct {
	EthereumAddress string `json:"ethereumAddress"`
	Signature       string `json:"signature"`
	Message         string `json:"message"`
}

// AdminAuthResponse structure
// Add any extra fields to the response here
type AdminAuthResponse struct {
	JWT *JWT `json:"jwt"`
}

// NewAdminAuthResponse creates a response with the model plus any other data
func NewAdminAuthResponse(obj *JWT) *AdminAuthResponse {
	resp := &AdminAuthResponse{JWT: obj}
	return resp
}

// Bind pre-processes any fields after a decode
func (req *AdminAuthRequest) Bind(r *http.Request) error {
	return nil
}

// Render post-processes the data before a response is returned
func (resp *AdminAuthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
