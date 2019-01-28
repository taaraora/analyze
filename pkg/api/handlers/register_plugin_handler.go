package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
)

type registerPluginHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewRegisterPluginHandler(storage storage.Interface, logger logrus.FieldLogger) operations.RegisterPluginHandler{
	return &registerPluginHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *registerPluginHandler) Handle(params operations.RegisterPluginParams) middleware.Responder {
	h.log.Debugf("got request at: %v, request: %+v", time.Now(), params)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if params.Body == nil || "" == strings.TrimSpace(params.Body.ID){
		r := operations.NewGetPluginsDefault(http.StatusBadRequest)
		message := "plugin id can't be empty"
		r.Payload = &models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
	}
	p := &models.Plugin{
		Description:   params.Body.Description,
		ID:            params.Body.ID,
		InstalledAt:   params.Body.InstalledAt,
		Name:          params.Body.Name,
		ServiceLabels: params.Body.ServiceLabels,
		ServiceName:   params.Body.ServiceName,
		Status:        params.Body.Status,
		Version:       params.Body.Version,
	}

	rawPlugin, err := p.MarshalBinary()
	if err != nil {
		r := operations.NewGetPluginsDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	//just replace all entity content
	_, err = h.storage.Get(ctx, models.PluginPrefix, p.ID)

	if err == storage.ErrNotFound {
		err = h.storage.Put(ctx, models.PluginPrefix, p.ID, rawPlugin)
		if err != nil {
			r := operations.NewGetPluginsDefault(http.StatusInternalServerError)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    http.StatusInternalServerError,
				Message: &msg,
			}
			return r
		}

		return operations.NewRegisterPluginCreated().WithPayload(p)
	}

	if err != nil {
		r := operations.NewGetPluginsDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	err = h.storage.Put(ctx, models.PluginPrefix, p.ID, rawPlugin)
	if err != nil {
		r := operations.NewGetPluginsDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	return operations.NewRegisterPluginOK().WithPayload(p)
}
