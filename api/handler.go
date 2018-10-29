package api

import (
	"context"
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

	pluginRaw, err := h.storage.GetAll(context.Background(), "/robot/plugins/")

	if err != nil {
		r := operations.NewGetRecommendationPluginsDefault(500)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    500,
			Message: &msg,
		}
		return r
	}

	result := &operations.GetRecommendationPluginsOKBody{
		InstalledRecommendationPlugins: []*models.RecommendationPlugin{},
		TotalCount:                     int64(len(pluginRaw)),
	}

	for _, rawPlugin := range pluginRaw {
		p := &models.RecommendationPlugin{}
		err := p.UnmarshalBinary(rawPlugin)
		if err != nil {
			r := operations.NewGetRecommendationPluginsDefault(500)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    500,
				Message: &msg,
			}
			return r
		}
		result.InstalledRecommendationPlugins = append(result.InstalledRecommendationPlugins, p)
	}

	return operations.NewGetRecommendationPluginsOK().WithPayload(result)
}
