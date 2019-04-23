package plugin

import (
	"context"
	"time"

	"google.golang.org/grpc/keepalive"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/supergiant/analyze/pkg/plugin/proto"
)

type PluginsSet map[string]*Client

type CloudProviderType string

const (
	AWSCloudProviderType CloudProviderType = "aws"
	DOCloudProviderType  CloudProviderType = "do"
)

type Client struct {
	conn *grpc.ClientConn
	proto.PluginClient
}

func NewClient(pluginServerAddress string) (*Client, error) {
	keepaliveCfg := keepalive.ClientParameters{
		Time:                60 * time.Minute,
		Timeout:             60 * time.Second,
		PermitWithoutStream: false,
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepaliveCfg),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, pluginServerAddress, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to dial plugin server side: %s ", pluginServerAddress)
	}
	c := proto.NewPluginClient(conn)

	return &Client{
		conn:         conn,
		PluginClient: c,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetConnTarget() string {
	return c.conn.Target()
}
