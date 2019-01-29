package api

import (
	"github.com/go-openapi/loads"
	"github.com/supergiant/analyze/pkg/api/operations"
	"testing"
)

func GetTestAPI(t *testing.T) *operations.AnalyzeAPI {
	t.Helper()
	swaggerSpec, err := loads.Analyzed(SwaggerJSON, "2.0")
	if err != nil {
		t.Fatal("unable to create spec analyzed document")
	}

	return operations.NewAnalyzeAPI(swaggerSpec)
}
