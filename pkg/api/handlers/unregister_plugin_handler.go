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

type unregisterPluginHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewUnregisterPluginHandler(storage storage.Interface, logger logrus.FieldLogger) operations.UnregisterPluginHandler {
	return &unregisterPluginHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *unregisterPluginHandler) Handle(params operations.UnregisterPluginParams) middleware.Responder {
	h.log.Debugf("got request at: %v, request: %+v", time.Now(), params)
	defer h.log.Debugf("request processing finished at: %v, request: %+v", time.Now(), params)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if "" == strings.TrimSpace(params.PluginID) {
		r := operations.NewUnregisterPluginDefault(http.StatusBadRequest)
		message := "plugin id can't be empty"
		r.Payload = &models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return r
	}

	_, err := h.storage.Get(ctx, models.PluginPrefix, params.PluginID)

	if err == storage.ErrNotFound {
		return operations.NewUnregisterPluginNotFound()
	}

	if err != nil {
		r := operations.NewUnregisterPluginDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	// TODO: need transaction or deleted flag
	err = h.storage.Delete(ctx, models.PluginPrefix, params.PluginID)
	if err != nil {
		r := operations.NewUnregisterPluginDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	return operations.NewUnregisterPluginNoContent()
}
