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

type replacePluginConfigHandler struct {
	storage storage.Interface
	log     logrus.FieldLogger
}

func NewReplacePluginConfigHandler(
	storage storage.Interface,
	logger logrus.FieldLogger,
) operations.ReplacePluginConfigHandler {
	return &replacePluginConfigHandler{
		storage: storage,
		log:     logger,
	}
}

func (h *replacePluginConfigHandler) Handle(params operations.ReplacePluginConfigParams) middleware.Responder {
	h.log.Debugf("got request at: %v, request: %+v", time.Now(), params)
	defer h.log.Debugf("request processing finished at: %v, request: %+v", time.Now(), params)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if "" == strings.TrimSpace(params.PluginID) {
		r := operations.NewReplacePluginConfigDefault(http.StatusBadRequest)
		message := pluginIDValidationErrorMessage
		r.Payload = &models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return r
	}

	if 0 == params.Body.ExecutionInterval {
		r := operations.NewReplacePluginConfigDefault(http.StatusBadRequest)
		message := "execution interval can't be 0"
		r.Payload = &models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return r
	}

	pc := &models.PluginConfig{
		EtcdEndpoints:        params.Body.EtcdEndpoints,
		ExecutionInterval:    params.Body.ExecutionInterval,
		PluginSpecificConfig: params.Body.PluginSpecificConfig,
	}

	rawPluginConfig, err := pc.MarshalBinary()
	if err != nil {
		r := operations.NewReplacePluginConfigDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	//just replace all entity content
	_, err = h.storage.Get(ctx, storage.PluginConfigPrefix, params.PluginID)

	if err == storage.ErrNotFound {
		err = h.storage.Put(ctx, storage.PluginConfigPrefix, params.PluginID, storageMessage(rawPluginConfig))
		if err != nil {
			r := operations.NewReplacePluginConfigDefault(http.StatusInternalServerError)
			msg := err.Error()
			r.Payload = &models.Error{
				Code:    http.StatusInternalServerError,
				Message: &msg,
			}
			return r
		}

		return operations.NewReplacePluginConfigOK()
	}

	if err != nil {
		r := operations.NewReplacePluginConfigDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	err = h.storage.Put(ctx, storage.PluginConfigPrefix, params.PluginID, storageMessage(rawPluginConfig))
	if err != nil {
		r := operations.NewReplacePluginConfigDefault(http.StatusInternalServerError)
		msg := err.Error()
		r.Payload = &models.Error{
			Code:    http.StatusInternalServerError,
			Message: &msg,
		}
		return r
	}

	return operations.NewReplacePluginConfigOK()
}
