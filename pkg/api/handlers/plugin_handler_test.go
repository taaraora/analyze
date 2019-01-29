package handlers_test

import (
	"github.com/sirupsen/logrus"
	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const fixturePlugin = `{"description":"detailed plugin description","id":"123456798","installedAt":"1970-01-01T00:00:00.000Z","name":"the name of the plugin","serviceLabels":{"app":"test"},"serviceName":"name of k8s service which is front of plugin deployment","status":"plugin status","version":"2.0.0"}`

func TestPluginHandler_ReturnResultsSuccessfully(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginHandler = handlers.NewPluginHandler(storage.GetMockStorage(t, map[string]string{
		models.PluginPrefix + "123456798": fixturePlugin,
	}), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/123456798", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	//TODO: investigate why it has extra spaces in the end
	if strings.TrimSpace(rr.Body.String()) != fixturePlugin {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), fixturePlugin)
	}
}

func TestPluginHandler_ReturnInternalError(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginHandler = handlers.NewPluginHandler(storage.GetMockBrokenStorage(t), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/123456798", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}
}

func TestPluginHandler_ReturnNotFound(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginHandler = handlers.NewPluginHandler(storage.GetMockStorage(t, nil), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/123456798", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}
}