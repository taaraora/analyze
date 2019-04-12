package handlers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
)

type pluginFixture models.Plugin

func (p pluginFixture) getPlugin() models.Plugin {
	return models.Plugin(p)
}

func (p pluginFixture) string() string {
	pp := models.Plugin(p)
	b, _ := pp.MarshalBinary()
	return string(b)
}

func newPluginFixture(id string) pluginFixture {
	d, _ := strfmt.ParseDateTime("2019-01-02T15:04:05Z07:00")
	p := models.Plugin{
		CheckComponentEntryPoint:    "/check/analyze-plugin-sunsetting-check-main.js",
		Description:                 "detailed plugin description",
		ID:                          id,
		InstalledAt:                 d,
		Name:                        "the name of the plugin",
		ServiceEndpoint:             "dfdsfsdfsd:8089",
		ServiceLabels:               nil,
		SettingsComponentEntryPoint: "/settings/analyze-plugin-sunsetting-settings-main.js",
		Status:                      "OK",
		Version:                     "v2.0.0",
	}

	return pluginFixture(p)
}

func toPlugin(t *testing.T, body *bytes.Buffer) *models.Plugin {
	t.Helper()

	p := &models.Plugin{}
	b, err := ioutil.ReadAll(body)

	if err != nil {
		t.Fatalf("can't read body: got %+v", err)
	}
	err = p.UnmarshalBinary(b)
	if err != nil {
		t.Fatalf("can't unmarshal body: got %+v", err)
	}
	return p
}

func toPlugins(t *testing.T, body *bytes.Buffer) []models.Plugin {
	t.Helper()

	p := make([]models.Plugin, 0)
	b, err := ioutil.ReadAll(body)

	if err != nil {
		t.Fatalf("can't read body: got %+v", err)
	}
	err = json.Unmarshal(b, &p)
	if err != nil {
		t.Fatalf("can't unmarshal body: got %+v", err)
	}
	return p
}

func TestPluginsHandler_ReturnResultsSuccessfully(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	fixturePlugins1 := newPluginFixture("123456798")
	fixturePlugins2 := newPluginFixture("1234567980")

	//TODO: create interface for logger, and use dummy logger for tests
	analyzeAPI.GetPluginsHandler = handlers.NewPluginsHandler(storage.GetMockStorage(t, map[string]string{
		models.PluginPrefix + "123456798":  fixturePlugins1.string(),
		models.PluginPrefix + "1234567980": fixturePlugins2.string(),
	}), logrus.New())
	server := api.NewServer(analyzeAPI)
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

	plugins := toPlugins(t, rr.Body)
	if !reflect.DeepEqual(plugins, []models.Plugin{fixturePlugins1.getPlugin(), fixturePlugins2.getPlugin()}) &&
		!reflect.DeepEqual(plugins, []models.Plugin{fixturePlugins2.getPlugin(), fixturePlugins1.getPlugin()}) {
		t.Fatalf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestPluginsHandler_ReturnInternalError(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeAPI.GetPluginsHandler = handlers.NewPluginsHandler(storage.GetMockBrokenStorage(t), logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins", nil)
	if err != nil {
		t.Fatal(err)
	}
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
