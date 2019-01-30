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

const fixturePlugins1 = `{"description":"detailed plugin description","id":"123456798","installedAt":"1970-01-01T00:00:00.000Z","name":"the name of the plugin","serviceLabels":{"app":"test"},"serviceName":"name of k8s service which is front of plugin deployment","status":"plugin status","version":"2.0.0"}`
const fixturePlugins2 = `{"description":"detailed plugin description2","id":"1234567980","installedAt":"1970-01-01T00:00:00.001Z","name":"the name of the plugin2","serviceLabels":{"app2":"test"},"serviceName":"name of k8s service which is front of plugin deployment2","status":"plugin status2","version":"2.0.1"}`

func TestPluginsHandler_ReturnResultsSuccessfully(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginsHandler = handlers.NewPluginsHandler(storage.GetMockStorage(t, map[string]string{
		models.PluginPrefix + "123456798": fixturePlugins1,
		models.PluginPrefix + "1234567980": fixturePlugins2,
	}), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	body := rr.Body.String()
	if  !strings.Contains(body, fixturePlugins1) || !strings.Contains(body, fixturePlugins2){
		t.Fatalf("handler returned unexpected body: got %v", body)
	}
}

func TestPluginsHandler_ReturnInternalError(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginsHandler = handlers.NewPluginsHandler(storage.GetMockBrokenStorage(t), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusInternalServerError, rr.Body.String())
	}
}