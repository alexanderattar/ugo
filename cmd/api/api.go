package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/consensys/ugo/pkg/api"
	"github.com/consensys/ugo/pkg/db"
	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/middleware"
	"github.com/consensys/ugo/pkg/utils"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	newrelic "github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
)

type Api struct {
	DB *sqlx.DB
}

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		// disable to use custom formatted timestamp
		DisableTimestamp: true,
	}

	utils.CreateIpfsShell()
	utils.CreateInfuraQueue()

	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger
	serverCtx := context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Starting API")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set as environment variable")
	}

	db, err := db.NewDB(dbURL)
	if err != nil {
		log.Panic(err)
	}

	// TODO - Better way of handling this?
	// New Relic requires a 40 string license key so we
	// default to this if the env var is not set
	newrelicKey := "----------------------------------------"
	if os.Getenv("NEWRELIC_KEY") != "" {
		newrelicKey = os.Getenv("NEWRELIC_KEY")
	}

	newRelicConfig := newrelic.NewConfig("Ujo", newrelicKey)
	nr, err := newrelic.NewApplication(newRelicConfig)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(lg.RequestLogger(logger))

	cors := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"}, // Use this to allow specific origin hosts
			AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           60, // Maximum value not ignored by any of major browsers
		},
	)

	r.Use(cors.Handler)
	r.Use(middleware.NewRelicMiddleware(nr))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// TODO - Require auth for these endpoints
	r.Post("/api/auth", api.AuthHandler{DB: db}.Auth)
	r.Post("/api/admin", api.AuthHandler{DB: db}.AdminAuth)
	r.Post("/api/dag/put", api.DagPutAndPin)
	r.Post("/api/dag/put/recursive", api.RecursiveDagPutAndPin)

	// Protected routes
	r.Group(func(r chi.Router) {
		// Handle valid / invalid tokens
		// This can be modified using the Authenticator method in auth.go
		r.Use(jwtauth.Authenticator)

		r.Post("/api/s3/sign", api.S3PutObjectHandler)

		// TODO - Build out admin
		// r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		// 	_, claims, err := jwtauth.FromContext(r.Context())
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	w.Write([]byte(fmt.Sprintf("Protected area. hi %v", claims["ethereumAddress"])))
		// })
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("你好, Ujo!"))
		})

		// API endpoints
		r.Mount("/api/persons", api.PersonResource{DB: db}.Routes())
		r.Mount("/api/musicgroups", api.MusicGroupResource{DB: db}.Routes())
		r.Mount("/api/musicreleases", api.MusicReleaseResource{DB: db}.Routes())
		r.Mount("/api/musicrecordings", api.MusicRecordingResource{DB: db}.Routes())
		r.Mount("/api/musicplaylists", api.MusicPlaylistResource{DB: db}.Routes())
		r.Mount("/api/purchases", api.PurchaseResource{DB: db}.Routes())
		r.Mount("/api/reports", api.ReportResource{DB: db}.Routes())
		r.Mount("/api/signedmessages", api.SignedMessageResource{DB: db}.Routes())
		r.Mount("/api/playevents", api.PlayEventResource{DB: db}.Routes())
		r.Mount("/api/payevents", api.PayEventResource{DB: db}.Routes())
	})

	http.ListenAndServe(":"+port, r)
}
