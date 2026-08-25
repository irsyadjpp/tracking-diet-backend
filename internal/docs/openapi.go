package docs

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

// NewOpenAPISpec creates a new OpenAPI 3.0 specification for the API
func NewOpenAPISpec() *openapi3.Spec {
	spec := openapi3.NewSpec()

	// Setup spec metadata
	info := openapi3.Info{
		Title:       "Tracking Diet Backend API",
		Description: "Comprehensive diet tracking backend with AI-powered recommendations, goal setting, and social features",
		Version:     "1.0.0",
		Contact: openapi3.Contact{
			Name:  "Irsyad Jamal Pratama Putra",
			Email: "mr.icadj@gmail.com",
		},
		License: openapi3.License{
			Name: "MIT",
		},
	}
	spec.Info.Set(info)

	// Setup servers
	spec.Servers = []openapi3.Server{
		{
			URL:         "http://localhost:8080",
			Description: "Development server",
		},
		{
			URL:         "https://api.tracking-diet.com",
			Description: "Production server",
		},
	}

	// Setup JWT security scheme
	jwtSecurity := openapi3.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "JWT authentication token",
	}
	spec.Components.SecuritySchemes["jwtAuth"] = jwtSecurity

	// Add security requirement to all protected routes
	spec.AddSecurityRequirement("jwtAuth", []string{})

	// Add basic operation for authentication
	err := spec.SetupOperation(http.MethodPost, "/api/v1/auth/register",
		func(op *openapi3.Operation) error {
			op.SetTags("authentication")
			op.SetSummary("Register a new user")
			op.SetID("register")
			
			type RegisterRequest struct {
				FullName    string `json:"full_name" required:"true" minLength:"1" maxLength:"200"`
				Email       string `json:"email" required:"true" format:"email"`
				Password    string `json:"password" required:"true" minLength:"8"`
				Gender      string `json:"gender" enum:"M,F,O"`
				BirthDate   string `json:"birth_date" format:"date"`
				HeightCm    *float64 `json:"height_cm"`
			}
			
			type AuthResponse struct {
				Message string `json:"message"`
			}
			
			if err := op.AddReqStructure(new(RegisterRequest)); err != nil {
				return err
			}
			
			if err := op.AddRespStructure(new(AuthResponse), func(cu *openapi3.ContentUnit) {
				cu.HTTPStatus = http.StatusCreated
			}); err != nil {
				return err
			}
			
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	// Add login operation
	err = spec.SetupOperation(http.MethodPost, "/api/v1/auth/login",
		func(op *openapi3.Operation) error {
			op.SetTags("authentication")
			op.SetSummary("Login user")
			op.SetID("login")
			
			type LoginRequest struct {
				Email    string `json:"email" required:"true" format:"email"`
				Password string `json:"password" required:"true"`
			}
			
			type LoginResponse struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				ExpiresIn    int    `json:"expires_in"`
				UserID       int    `json:"user_id"`
				Email        string `json:"email"`
			}
			
			if err := op.AddReqStructure(new(LoginRequest)); err != nil {
				return err
			}
			
			if err := op.AddRespStructure(new(LoginResponse), func(cu *openapi3.ContentUnit) {
				cu.HTTPStatus = http.StatusOK
			}); err != nil {
				return err
			}
			
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	return spec
}