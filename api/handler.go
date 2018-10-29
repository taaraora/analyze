package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/supergiant/robot/pkg/storage"
	"github.com/supergiant/robot/swagger/gen/models"
	"github.com/supergiant/robot/swagger/gen/restapi/operations"
)

type recommendationPluginsHandler struct {
	storage storage.Interface
}

func NewRecommendationPluginsHandler(storage storage.Interface) operations.GetRecommendationPluginsHandler {
	return &recommendationPluginsHandler{
		storage: storage,
	}
}

func (h *recommendationPluginsHandler) Handle(params operations.GetRecommendationPluginsParams) middleware.Responder {

	result := &operations.GetRecommendationPluginsOKBody{
		InstalledRecommendationPlugins: []*models.RecommendationPlugin{},
		TotalCount:                     1,
	}

	result.InstalledRecommendationPlugins = append(result.InstalledRecommendationPlugins, &models.RecommendationPlugin{
		Description: "",
		ID:          "d6fde92930d4715a2b49857d24b940956b26d2d3",
		InstalledAt: "2018-05-04T01:14:52Z",
		Name:        "limit/requests checker",
		Status:      "OK",
		Version:     "v0.0.1",
	})

	return operations.NewGetRecommendationPluginsOK().WithPayload(result)
}
