package data

import (
	"context"
	"errors"
	log "go-micro.dev/v4/util/log"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"saas/internal/biz"
	"strconv"
)

type k8sRepo struct {
	k8sClientSet *kubernetes.Clientset
	deployment   *v1.Deployment
}

// CreateToK8S 创建deployment
func (k *k8sRepo) CreateToK8S(pod *biz.Pod) (err error) {
	k.SetDeployment(pod)
	if _, err = k.GetPodByName(pod); err != nil {
		if _, err = k.k8sClientSet.AppsV1().Deployments(pod.Namespace).Create(
			context.TODO(), k.deployment, v12.CreateOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.New("Pod " + pod.Name + "已经存在")
}

// UpdateToK8S 更新deployment
func (k *k8sRepo) UpdateToK8S(pod *biz.Pod) (err error) {
	k.SetDeployment(pod)
	if _, err = k.GetPodByName(pod); err != nil {
		log.Error(err)
		return errors.New("Pod " + pod.Name + " 不存在请先创建")
	}
	if _, err = k.k8sClientSet.AppsV1().Deployments(pod.Namespace).Update(
		context.TODO(), k.deployment, v12.UpdateOptions{}); err != nil {
		log.Error(err)
		return err
	}
	log.Info(pod.Name + " 更新成功")
	return nil
}

func (k *k8sRepo) DeleteFromK8S(pod *biz.Pod) (err error) {
	if err = k.k8sClientSet.AppsV1().Deployments(pod.Namespace).Delete(
		context.TODO(), pod.Name, v12.DeleteOptions{}); err != nil {
		log.Error(err)
		return
	}
	return
}

func (k *k8sRepo) SetDeployment(pod *biz.Pod) {
	deployment := &v1.Deployment{}
	deployment.TypeMeta = v12.TypeMeta{
		Kind:       "deployment",
		APIVersion: "v1",
	}
	deployment.ObjectMeta = v12.ObjectMeta{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Labels: map[string]string{
			"app-name": pod.Name,
			"author":   "Caplost",
		},
	}
	deployment.Name = pod.Name
	deployment.Spec = v1.DeploymentSpec{
		Replicas: &pod.Replicas,
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": pod.Name,
			},
			MatchExpressions: nil,
		},
		Template: v13.PodTemplateSpec{
			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"app-name": pod.Name,
				},
			},
			Spec: v13.PodSpec{
				Containers: []v13.Container{
					{
						Name:            pod.Name,
						Image:           pod.Image,
						Ports:           getContainerPorts(pod),
						Env:             getEnvs(pod),
						Resources:       getResources(pod),
						ImagePullPolicy: getImagePullPolicy(pod),
					},
				},
			},
		},
		Strategy:                v1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}
	k.deployment = deployment
}

func (k *k8sRepo) GetPodByName(pod *biz.Pod) (*v1.Deployment, error) {
	return k.k8sClientSet.AppsV1().Deployments(pod.Namespace).Get(
		context.TODO(), pod.Name, v12.GetOptions{})
}

func getImagePullPolicy(pod *biz.Pod) v13.PullPolicy {
	switch pod.PullPolicy {
	case "Always":
		return "Always"
	case "Never":
		return "Never"
	case "IfNotPresent":
		return "IfNotPresent"
	default:
		return "Always"
	}
}

func getContainerPorts(pod *biz.Pod) (containerPorts []v13.ContainerPort) {
	for _, v := range pod.Ports {
		containerPorts = append(containerPorts, v13.ContainerPort{
			Name:          "port-" + strconv.Itoa(int(v.Port)),
			ContainerPort: v.Port,
			Protocol:      getProtocol(v.Protocol),
		})
	}
	return
}

func getEnvs(pod *biz.Pod) (envVars []v13.EnvVar) {
	for _, v := range pod.Envs {
		envVars = append(envVars, v13.EnvVar{
			Name:  v.Key,
			Value: v.Val,
		})
	}
	return
}

func getResources(pod *biz.Pod) (require v13.ResourceRequirements) {
	require.Limits = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(pod.CpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(pod.MemoryMax), 'f', 6, 64)),
	}
	require.Requests = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(pod.CpuMin), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(pod.MemoryMin), 'f', 6, 64)),
	}
	return
}

func getProtocol(protocol string) v13.Protocol {
	switch protocol {
	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	case "SCTP":
		return "SCTP"
	default:
		return "TCP"
	}
}
