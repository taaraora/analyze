package underutilizednodes

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/supergiant/robot/pkg/plugin/proto"
)

type uuNodesPlugin struct {
	config *proto.PluginConfig
}

func NewPlugin() proto.PluginClient {
	return &uuNodesPlugin{}
}

func (u *uuNodesPlugin) Check(ctx context.Context, in *proto.CheckRequest, opts ...grpc.CallOption) (*proto.CheckResponse, error) {
	return nil, nil
}

func (u *uuNodesPlugin) Action(ctx context.Context, in *proto.ActionRequest, opts ...grpc.CallOption) (*proto.ActionResponse, error) {
	panic("implement me")
}

func (u *uuNodesPlugin) Configure(ctx context.Context, in *proto.PluginConfig, opts ...grpc.CallOption) (*proto.Empty, error) {
	return nil, nil
}

func (u *uuNodesPlugin) Stop(ctx context.Context, in *proto.Stop_Request, opts ...grpc.CallOption) (*proto.Stop_Response, error) {
	panic("implement me")
}

func (u *uuNodesPlugin) Info(ctx context.Context, in *proto.Empty, opts ...grpc.CallOption) (*proto.PluginInfo, error) {
	return &proto.PluginInfo{
		Id:          "supergiant-underutilized-nodes-plugin",
		Version:     "v0.0.1",
		Name:        "Underutilized nodes fixing plugin",
		Description: "This plugin checks nodes using high intelligent Kellys approach to find underutilized nodes, than calculates how it is possible to fix that",
	}, nil
}
