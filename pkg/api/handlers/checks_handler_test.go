package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/supergiant/analyze/pkg/storage"

	"github.com/supergiant/analyze/pkg/storage/mock"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
)

//nolint
const fixtureCheckResult = `{"checkStatus":"GREEN","completedAt":"1970-01-01T00:00:00.000Z","description":"detailed check result description","executionStatus":"no execution errors","id":"uniqueUUID","name":"some interesting check"}`

func TestChecksResultsHandler_ReturnResultsSuccessfully(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeAPI.GetCheckResultsHandler = handlers.NewChecksResultsHandler(mock.GetMockStorage(t, map[string]string{
		storage.CheckResultPrefix + "uniqueUUID": fixtureCheckResult,
	}), logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/checks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	//TODO: investigate why it has extra spaces in the end
	if strings.TrimSpace(rr.Body.String()) != "["+fixtureCheckResult+"]" {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), "["+fixtureCheckResult+"]")
	}
}

func TestChecksResultsHandler_ReturnInternalError(t *testing.T) {
	analyzeAPI := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeAPI.GetCheckResultsHandler = handlers.NewChecksResultsHandler(mock.GetMockBrokenStorage(t), logrus.New())
	server := api.NewServer(analyzeAPI)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/checks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}
}
