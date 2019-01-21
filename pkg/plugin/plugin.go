package plugin

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/supergiant/analyze/pkg/plugin/proto"
)

type Config struct {
	ProviderType CloudProviderType `mapstructure:"cloud_provider_type"`
	// TODO: refactor config when multiple providers become being supported
	AWSAccessKeyID     string        `mapstructure:"aws_access_key_id"`
	AWSSecretAccessKey string        `mapstructure:"aws_secret_access_key"`
	AWSRegion          string        `mapstructure:"aws_region"`
	CheckInterval      time.Duration `mapstructure:"check_interval"`
	CheckTimeout       time.Duration `mapstructure:"check_timeout"`
}

type PluginsSet map[string]*Client

type CloudProviderType string

const (
	AWSCloudProviderType CloudProviderType = "aws"
	DOCloudProviderType  CloudProviderType = "do"
)

type Client struct {
	cfg  *proto.PluginConfig
	conn *grpc.ClientConn
	proto.PluginClient
}

func NewClient(pluginServerAddress string, cfg *proto.PluginConfig) (*Client, error) {
	conn, err := grpc.Dial(pluginServerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to dial plugin server side: %s ", pluginServerAddress)
	}
	c := proto.NewPluginClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	_, err = c.Configure(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to configure plugin")
	}

	return &Client{
		cfg:          cfg,
		conn:         conn,
		PluginClient: c,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (ps PluginsSet) Load(plugin proto.PluginClient, cfg *proto.PluginConfig) error {

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	info, err := plugin.Info(ctx, &empty.Empty{})
	if err != nil {
		return errors.Wrap(err, "unable to get plugin info")
	}

	ps[info.Id] = plugin

	return nil
}

func (c Config) Validate() error {
	if c.CheckInterval.Nanoseconds() <= 0 {
		return errors.New("plugin check interval can't be less or equal to zero")
	}

	if c.CheckTimeout.Nanoseconds() <= 0 {
		return errors.New("plugin check timeout can't be less or equal to zero")
	}

	if c.ProviderType != AWSCloudProviderType {
		return errors.New("Currently only AWS is implemented")
	}

	return nil
}

func (c Config) ToProtoConfig() *proto.PluginConfig {
	return &proto.PluginConfig{
		ProviderType: newProviderType(c.ProviderType),
		CloudProviderConfig: &proto.PluginConfig_AwsConfig{
			AwsConfig: &proto.AwsConfig{
				AccessKeyId:     c.AWSAccessKeyID,
				SecretAccessKey: c.AWSSecretAccessKey,
				Region:          c.AWSRegion,
			},
		},
	}
}

func newProviderType(t CloudProviderType) proto.CloudProviderType {
	switch t {
	case DOCloudProviderType:
		return proto.CloudProviderType_DO
	default:
		return proto.CloudProviderType_AWS
	}
}
