package kube

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientSet *kubernetes.Clientset
	logger    logrus.FieldLogger
}

func NewKubeClient(logger logrus.FieldLogger) (*Client, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the client
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientSet: clientSet,
		logger:    logger,
	}, nil
}

func (c *Client) GetService(serviceName string, labelsSet map[string]string) (corev1.Service, error) {
	var labelsSelector = labels.SelectorFromSet(labels.Set(labelsSet))
	var options = metav1.ListOptions{
		LabelSelector: labelsSelector.String(),
		FieldSelector: fields.OneTermEqualSelector("metadata.name", serviceName).String(),
	}

	serviceList, err := c.clientSet.CoreV1().Services("").List(options)
	if err != nil {
		return corev1.Service{}, errors.Wrap(err, "can't list services")
	}

	if len(serviceList.Items) != 1 {
		return corev1.Service{}, errors.Wrap(err, "multiple services are deployed")
	}

	return serviceList.Items[0], nil
}
