package handlers

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/supergiant/robot/pkg/api/operations"
	"github.com/supergiant/robot/pkg/models"
	"github.com/supergiant/robot/pkg/storage"
)

type checkResultsHandler struct {
	storage storage.Interface
}

func NewCheckResultsHandler(storage storage.Interface) operations.GetCheckResultsHandler {
	return &checkResultsHandler{
		storage: storage,
	}
}

func (h *checkResultsHandler) Handle(params operations.GetCheckResultsParams) middleware.Responder {

	resultsRaw, err := h.storage.GetAll(context.Background(), "/robot/check_results/")

	if err != nil {
		r := operations.NewGetCheckResultsDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	result := &operations.GetCheckResultsOKBody{
		CheckResults: []*models.CheckResult{},
	}

	for _, rawResult := range resultsRaw {
		checkResult := &models.CheckResult{}
		err := checkResult.UnmarshalBinary(rawResult)
		if err != nil {
			r := operations.NewGetCheckResultsDefault(http.StatusInternalServerError)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    http.StatusInternalServerError,
				Message: &msg,
			}
			return r
		}
		result.CheckResults = append(result.CheckResults, checkResult)
	}

	return operations.NewGetCheckResultsOK().WithPayload(result)
}
