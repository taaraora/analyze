package handlers_test

import (
	"context"
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

func TestRegisterPluginHandler_ReturnCreated(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockStorage(t, nil)
	analyzeApi.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(fixturePlugins1))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusCreated, rr.Body.String())
	}

	//TODO: investigate why it has extra spaces in the end
	if strings.TrimSpace(rr.Body.String()) != fixturePlugins1 {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), fixturePlugins1)
	}

	b, _ := strg.Get(context.TODO(), models.PluginPrefix, "123456798")
	p := &models.Plugin{}
	if err := (p).UnmarshalBinary(b); err != nil {
		t.Fatalf("handler put in storage something broken")
	}

	if string(b) != fixturePlugins1 {
		t.Fatalf("storage returned unexpected content: got %v want %v", string(b), fixturePlugins1)
	}
}

func TestRegisterPluginHandler_ReturnUpdated(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockStorage(t, map[string]string{
		models.PluginPrefix + "123456798": fixturePlugins1,
	})
	analyzeApi.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	tmp := strings.Replace(fixturePlugins1, "the name of the plugin", "new super name of the plugin", 1)
	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(tmp))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	//TODO: investigate why it has extra spaces in the end
	if strings.TrimSpace(rr.Body.String()) != tmp {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), tmp)
	}

	b, _ := strg.Get(context.TODO(), models.PluginPrefix, "123456798")
	p := &models.Plugin{}
	if err := (p).UnmarshalBinary(b); err != nil {
		t.Fatalf("handler put in storage something broken")
	}

	if string(b) != tmp {
		t.Fatalf("storage returned unexpected content: got %v want %v", string(b), fixturePlugins1)
	}
}

func TestRegisterPluginHandler_ReturnInternalError(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockBrokenStorage(t)
	analyzeApi.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(fixturePlugins1))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusInternalServerError, rr.Body.String())
	}
}

func TestRegisterPluginHandler_ReturnBadRequest(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockStorage(t, nil)
	analyzeApi.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	tmp := strings.Replace(fixturePlugins1, "123456798", " ", 1)
	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(tmp))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusBadRequest, rr.Body.String())
	}
}