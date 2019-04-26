package analyze

import (
	"context"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"

	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/plugin"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"github.com/supergiant/analyze/pkg/storage"
)

func CheckJob(pluginID string, pluginClient *plugin.Client, stor storage.Interface) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		defer cancel()
		checkResponse, err := pluginClient.Check(ctx, &proto.CheckRequest{})
		if err != nil {

			return errors.Errorf("unable to execute check for pluginClient: %s, error: %v", pluginID, err)
		}
		if checkResponse.Error != "" {

			return errors.Errorf("pluginClient: %s, returned error: %s", pluginID, checkResponse.Error)
		}
		if checkResponse.Result == nil {

			return errors.Errorf("pluginClient: %s, returned nil Result", pluginID)
		}

		r := checkResponse.Result

		var currentTime = time.Now()
		checkResult := models.CheckResult{
			CheckStatus:     r.GetStatus().String(),
			CompletedAt:     strfmt.DateTime(currentTime),
			Description:     string(r.GetDescription().Value),
			ExecutionStatus: r.GetExecutionStatus(),
			ID:              pluginID,
			Name:            r.GetName(),
		}

		bytes, err := checkResult.MarshalBinary()
		if err != nil {
			return errors.Errorf("unable to marshal check result, pluginClient: %s, returned error: %s", pluginID, err)
		}

		err = stor.Put(ctx, storage.CheckResultPrefix, pluginID, msg(bytes))
		if err != nil {
			return errors.Errorf("unable to store check result, pluginClient: %s, returned error: %s", pluginID, err)
		}

		return nil
	}
}

type msg []byte

func (d msg) Payload() []byte {
	return d
}
