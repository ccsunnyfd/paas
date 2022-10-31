package config

import (
	"go-micro.dev/v4/config"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8S struct {
	ConfigPath string `json:"config-path"`
}

func GetK8SFromConfig(config config.Config, path ...string) *K8S {
	k8sConfig := &K8S{}

	config.Get(path...).Scan(k8sConfig)
	return k8sConfig
}

func GetK8SFromCMD(kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)

}
