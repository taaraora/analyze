package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
)

type pluginHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewPluginHandler(storage storage.Interface, logger logrus.FieldLogger) operations.GetPluginHandler {
	return &pluginHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *pluginHandler) Handle(params operations.GetPluginParams) middleware.Responder {
	h.log.Debugf("got request at: %v, request: %+v", time.Now(), params)
	defer h.log.Debugf("request processing finished at: %v, request: %+v", time.Now(), params)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	pluginRaw, err := h.storage.Get(ctx, models.PluginPrefix, params.PluginID)

	if err == storage.ErrNotFound {
		r := operations.NewGetPluginDefault(http.StatusNotFound)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusNotFound,
			Message: &msg,
		}
		return r
	}

	if err != nil {
		r := operations.NewGetPluginDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	p := &models.Plugin{}
	err = p.UnmarshalBinary(pluginRaw)
	if err != nil {
		r := operations.NewGetPluginDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	return operations.NewGetPluginOK().WithPayload(p)
}
