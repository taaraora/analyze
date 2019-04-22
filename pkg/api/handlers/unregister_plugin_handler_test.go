package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/supergiant/analyze/pkg/storage"

	"github.com/supergiant/analyze/pkg/storage/mock"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
)

func TestUnregisterPluginHandler_ReturnBadRequest(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	strg := mock.GetMockStorage(t, nil)
	analyzeAPI.UnregisterPluginHandler = handlers.NewUnregisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("DELETE", "/api/v1/plugins/ ", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf(
			"handler returned wrong status code: got %v want %v, body: %v",
			status,
			http.StatusBadRequest,
			rr.Body.String(),
		)
	}
}

func TestUnregisterPluginHandler_ReturnNotFound(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	strg := mock.GetMockStorage(t, nil)
	analyzeAPI.UnregisterPluginHandler = handlers.NewUnregisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("DELETE", "/api/v1/plugins/12341234", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf(
			"handler returned wrong status code: got %v want %v, body: %v",
			status,
			http.StatusNotFound,
			rr.Body.String(),
		)
	}
}

func TestUnregisterPluginHandler_ReturnInternalServerError(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	strg := mock.GetMockBrokenStorage(t)
	analyzeAPI.UnregisterPluginHandler = handlers.NewUnregisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("DELETE", "/api/v1/plugins/12341234", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf(
			"handler returned wrong status code: got %v want %v, body: %v",
			status,
			http.StatusInternalServerError,
			rr.Body.String(),
		)
	}
}

func TestUnregisterPluginHandler_ReturnNoContent(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	fixturePlugins1 := newPluginFixture("123456798")
	//TODO: create interface for logger, and use dummy logger for tests
	strg := mock.GetMockStorage(t, map[string]string{
		storage.PluginPrefix + "123456798": fixturePlugins1.string(),
	})
	analyzeAPI.UnregisterPluginHandler = handlers.NewUnregisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("DELETE", "/api/v1/plugins/123456798", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Fatalf(
			"handler returned wrong status code: got %v want %v, body: %v",
			status,
			http.StatusNoContent,
			rr.Body.String(),
		)
	}
}
