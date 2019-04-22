package handlers_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/supergiant/analyze/pkg/storage"

	"github.com/supergiant/analyze/pkg/storage/mock"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/models"
)

type pluginConfigFixture models.PluginConfig

const fixturePluginID = "analyze-plugin-sunsetting"

func (pc pluginConfigFixture) getPluginConfig() models.PluginConfig {
	return models.PluginConfig(pc)
}

func (pc pluginConfigFixture) string() string {
	pp := models.PluginConfig(pc)
	b, _ := pp.MarshalBinary()
	return string(b)
}

func newPluginConfigFixture() pluginConfigFixture {

	pluginSpecificConf := map[string]interface{}{
		"Endpoints": []string{"http://fake-qc-address/sunset?node1&node3"},
	}

	pc := models.PluginConfig{
		ExecutionInterval:    180,
		PluginSpecificConfig: pluginSpecificConf,
	}

	return pluginConfigFixture(pc)
}

func toPluginConfig(t *testing.T, body *bytes.Buffer) *models.PluginConfig {
	t.Helper()

	pc := &models.PluginConfig{}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatalf("can't read body: got %+v", err)
	}

	err = pc.UnmarshalBinary(b)
	if err != nil {
		t.Fatalf("can't unmarshal body: got %+v", err)
	}
	return pc
}

func TestPluginConfigHandler_ReturnResultsSuccessfully(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	fixturePluginConfig := newPluginConfigFixture()

	//TODO: create interface for logger, and use dummy logger for tests
	analyzeAPI.GetPluginConfigHandler = handlers.NewPluginConfigHandler(mock.GetMockStorage(t, map[string]string{
		storage.PluginConfigPrefix + fixturePluginID: fixturePluginConfig.string(),
	}), logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/"+fixturePluginID+"/config", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	p := toPluginConfig(t, rr.Body)
	if reflect.DeepEqual(*p, fixturePluginConfig.getPluginConfig()) {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), fixturePluginConfig.string())
	}
}

func TestPluginConfigHandler_ReturnInternalError(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeAPI.GetPluginConfigHandler = handlers.NewPluginConfigHandler(mock.GetMockBrokenStorage(t), logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/"+fixturePluginID+"/config", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}
}

func TestPluginConfigHandler_ReturnNotFound(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeAPI.GetPluginConfigHandler = handlers.NewPluginConfigHandler(mock.GetMockStorage(t, nil), logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/"+fixturePluginID+"/config", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}
}
