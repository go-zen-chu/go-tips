//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient interface {
	GetConfigMap(namespace string, configMap string) (*v1.ConfigMap, error)
	UpdateConfigMap(namespace string, configMap *v1.ConfigMap) error
}

type kubeClient struct {
	ctx    context.Context
	client kubernetes.Interface
}

func NewKubernetesClient() (KubernetesClient, error) {
	client, err := newClient()
	if err != nil {
		return nil, fmt.Errorf("failed init kubeclient: %s", err)
	}
	return &kubeClient{
		ctx:    context.Background(),
		client: client,
	}, nil
}

func newClient() (kubernetes.Interface, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeConfig)
}

func (c *kubeClient) GetConfigMap(namespace string, configMap string) (*v1.ConfigMap, error) {
	cm, err := c.client.CoreV1().ConfigMaps(namespace).Get(c.ctx, configMap, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get configmap %s/%s: %s", namespace, configMap, err)
	}
	return cm, err
}

func (c *kubeClient) UpdateConfigMap(namespace string, configMap *v1.ConfigMap) error {
	_, err := c.client.CoreV1().ConfigMaps(namespace).Update(c.ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("updating configmap %s/%s: %s", namespace, configMap, err)
	}
	return nil
}
