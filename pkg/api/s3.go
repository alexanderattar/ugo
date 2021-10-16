package api

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/consensys/ugo/pkg/lg"
	"github.com/go-chi/render"
)

// S3PutObjectHandler generates a signed URL for a PUT object request
func S3PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	lg.RequestLog(r).Infoln("Generating signed URL")
	if r.Method != "POST" {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}

	data := &S3PutObjectRequest{}
	if err := render.Bind(r, data); err != nil {
		lg.Errorf("Error reading request body %v", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		lg.Errorf("Error initializing session %v", err)
		render.Render(w, r, Error400(err))
		return
	}

	// Create S3 service client
	svc := s3.New(sess)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(data.Bucket),
		Key:    aws.String(data.Filename),
	})

	signedURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		lg.Errorf("Error generating signed URL %v", err)
		render.Render(w, r, Error400(err))
		return
	}

	lg.Printf("Generated a signed URL %v", signedURL)
	render.Render(w, r, NewS3PutObjectResponse(signedURL))
}

// S3PutObjectRequest structure
type S3PutObjectRequest struct {
	Filename string `json:"filename"`
	Bucket   string `json:"bucket"`
}

// S3PutObjectResponse structure
// Add any extra fields to the response here
type S3PutObjectResponse struct {
	SignedURL string `json:"signedUrl"`
}

// NewS3PutObjectResponse creates a response with the model plus any other data
func NewS3PutObjectResponse(signedURL string) *S3PutObjectResponse {
	resp := &S3PutObjectResponse{SignedURL: signedURL}
	return resp
}

// Bind pre-processes any fields after a decode
func (req *S3PutObjectRequest) Bind(r *http.Request) error {
	return nil
}

// Render post-processes the data before a response is returned
func (resp *S3PutObjectResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
