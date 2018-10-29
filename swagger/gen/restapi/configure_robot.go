// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"github.com/rakyll/statik/fs"
	"github.com/supergiant/robot/swagger/gen/models"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/supergiant/robot/swagger/gen/restapi/operations"

	_ "github.com/supergiant/robot/statik"
)

//go:generate swagger generate server --target ../gen --name Robot --spec ../swagger/api-spec.yml --exclude-main

func configureFlags(api *operations.RobotAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.RobotAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.GetRecomendationPluginsHandler = operations.GetRecomendationPluginsHandlerFunc(func(params operations.GetRecomendationPluginsParams) middleware.Responder {

		result := &operations.GetRecomendationPluginsOKBody{
			InstalledRecommendationPlugins: []*models.RecommendationPlugin{},
			TotalCount:                     1,
		}

		result.InstalledRecommendationPlugins = append(result.InstalledRecommendationPlugins, &models.RecommendationPlugin{
			Description: "",
			ID:          "d6fde92930d4715a2b49857d24b940956b26d2d3",
			InstalledAt: "2018-05-04T01:14:52Z",
			Name:        "limit/requests checker",
			Status:      "OK",
			Version:     "v0.0.1",
		})

		return operations.NewGetRecomendationPluginsOK().WithPayload(result)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return uiMiddleware(handler)
}

func uiMiddleware(handler http.Handler) http.Handler {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	fmt.Println("ONE")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			fmt.Println("TWO")
			http.StripPrefix("/swagger-ui/", staticServer).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
