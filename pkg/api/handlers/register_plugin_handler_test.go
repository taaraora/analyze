package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
)

func TestRegisterPluginHandler_ReturnCreated(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	fixturePlugins1 := newPluginFixture("123456798")
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockStorage(t, nil)
	analyzeAPI.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(fixturePlugins1.string()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusCreated, rr.Body.String())
	}

	pBody := toPlugin(t, rr.Body)
	if !reflect.DeepEqual(*pBody, fixturePlugins1.getPlugin()) {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), fixturePlugins1.string())
	}

	b, _ := strg.Get(context.TODO(), models.PluginPrefix, "123456798")
	p := &models.Plugin{}
	if err := (p).UnmarshalBinary(b.Payload()); err != nil {
		t.Fatalf("handler put in storage something broken")
	}

	if !reflect.DeepEqual(*p, fixturePlugins1.getPlugin()) {
		t.Fatalf("storage returned unexpected content: got %v want %v", string(b.Payload()), fixturePlugins1)
	}
}

func TestRegisterPluginHandler_ReturnUpdated(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	fixturePlugins := newPluginFixture("123456798")
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockStorage(t, map[string]string{
		models.PluginPrefix + "123456798": fixturePlugins.string(),
	})
	analyzeAPI.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	fixturePlugins.Name = "new super name of the plugin"

	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(fixturePlugins.string()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	p := toPlugin(t, rr.Body)

	if !reflect.DeepEqual(*p, fixturePlugins.getPlugin()) {
		t.Fatalf("handler returned unexpected body: got %+v want %+v", *p, fixturePlugins.getPlugin())
	}

	b2, _ := strg.Get(context.TODO(), models.PluginPrefix, "123456798")

	var buffer bytes.Buffer
	_, err = buffer.Write(b2.Payload())
	if err != nil {
		t.Fatal("can't write to buffer")
	}
	p2 := toPlugin(t, &buffer)

	if !reflect.DeepEqual(*p2, fixturePlugins.getPlugin()) {
		t.Fatalf("storage returned unexpected content: got %+v want %+v", *p2, fixturePlugins.getPlugin())
	}
}

func TestRegisterPluginHandler_ReturnInternalError(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	fixturePlugins1 := newPluginFixture("123456798")
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockBrokenStorage(t)
	analyzeAPI.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(fixturePlugins1.string()))
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

func TestRegisterPluginHandler_ReturnBadRequest(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	fixturePlugins1 := newPluginFixture("123456798")
	//TODO: create interface for logger, and use dummy logger for tests
	strg := storage.GetMockStorage(t, nil)
	analyzeAPI.RegisterPluginHandler = handlers.NewRegisterPluginHandler(strg, logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	fixturePlugins1.ID = " "
	req, err := http.NewRequest("POST", "/api/v1/plugins", strings.NewReader(fixturePlugins1.string()))
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
