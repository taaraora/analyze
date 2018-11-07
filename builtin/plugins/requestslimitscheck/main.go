package requestslimitscheck

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/supergiant/robot/pkg/plugin/proto"
)

type resourceRequirementsPlugin struct {
	config *proto.PluginConfig
}

func NewPlugin() proto.PluginClient {
	return &resourceRequirementsPlugin{}
}

func (u *resourceRequirementsPlugin) Check(ctx context.Context, in *proto.CheckRequest, opts ...grpc.CallOption) (*proto.CheckResponse, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	var result = &proto.CheckResult{
		ExecutionStatus: "OK",
		Status:          proto.CheckStatus_UNKNOWN_CHECK_STATUS,
		Name:            "Resources (CPU/RAM) requests and limits Check",
		Description:     "Resources (CPU/RAM) requests and limits where checked on nodes of k8s cluster. Results: \n",
		Actions: []*proto.Action{
			&proto.Action{
				ActionId:    "1",
				Description: "Dismiss notification",
			},
			&proto.Action{
				ActionId:    "2",
				Description: "Set missing requests/limits",
			},
		},
	}

	nodes, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		result.ExecutionStatus = err.Error()

		return &proto.CheckResponse{Result: result}, nil
	}
	for _, node := range nodes.Items {
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + node.Name,
		})
		if err != nil {
			if err != nil {
				result.ExecutionStatus = err.Error()

				return &proto.CheckResponse{Result: result}, nil
			}
		}

		for _, pod := range pods.Items {
			result.Description += " PodName: " + pod.Name
			for _, container := range pod.Spec.Containers {
				description, status := describeResourceRequirements(container)
				result.Description += description
				setHigher(result, status)
			}
		}
	}

	return &proto.CheckResponse{Result: result}, nil
}

func (u *resourceRequirementsPlugin) Action(ctx context.Context, in *proto.ActionRequest, opts ...grpc.CallOption) (*proto.ActionResponse, error) {
	panic("implement me")
}

func (u *resourceRequirementsPlugin) Configure(ctx context.Context, in *proto.PluginConfig, opts ...grpc.CallOption) (*proto.Empty, error) {
	return nil, nil
}

func (u *resourceRequirementsPlugin) Stop(ctx context.Context, in *proto.Stop_Request, opts ...grpc.CallOption) (*proto.Stop_Response, error) {
	panic("implement me")
}

func (u *resourceRequirementsPlugin) Info(ctx context.Context, in *proto.Empty, opts ...grpc.CallOption) (*proto.PluginInfo, error) {
	return &proto.PluginInfo{
		Id:      "supergiant-resources-requests-and-limits-check-plugin",
		Version: "v0.0.1",
		Name:    "Resources (CPU/RAM) requests and limits",
		Description: "This plugin checks resources (CPU/RAM) requests and limits on pods. " +
			"It returns Green when limits and requests are set, Yellow when limits are not set, " +
			"and Red when requests are not set.",
	}, nil
}

func setHigher(ch *proto.CheckResult, status proto.CheckStatus) {
	if status > ch.Status {
		ch.Status = status
	}
}

func describeResourceRequirements(container v1.Container) (string, proto.CheckStatus) {
	var resultDescription = " ContainerName: " + container.Name
	var resultStatus = proto.CheckStatus_GREEN
	var limitIsAbsent bool
	var requestIsAbsent bool

	if !container.Resources.Limits.Cpu().IsZero() {
		resultDescription += "CPU limit: " + container.Resources.Limits.Cpu().String() + " "
	} else {
		resultDescription += "CPU limit: is Not Set "
		limitIsAbsent = true
	}

	if !container.Resources.Limits.Memory().IsZero() {
		resultDescription += "RAM limit: " + container.Resources.Limits.Memory().String() + " "
	} else {
		resultDescription += "RAM limit: is Not Set "
		limitIsAbsent = true
	}

	if !container.Resources.Requests.Cpu().IsZero() {
		resultDescription += "CPU request: " + container.Resources.Requests.Cpu().String() + " "
	} else {
		resultDescription += "CPU request: is Not Set "
		requestIsAbsent = true
	}

	if !container.Resources.Requests.Memory().IsZero() {
		resultDescription += "RAM request: " + container.Resources.Requests.Memory().String() + " "
	} else {
		resultDescription += "RAM request: is Not Set "
		requestIsAbsent = true
	}

	if limitIsAbsent {
		resultStatus = proto.CheckStatus_YELLOW
	}

	if requestIsAbsent {
		resultStatus = proto.CheckStatus_RED
	}

	return resultDescription, resultStatus
}
