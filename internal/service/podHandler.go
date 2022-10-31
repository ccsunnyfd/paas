package service

import (
	"context"
	v1 "saas/api/pod/v1"
	"saas/internal/biz"
)

type PodHandler struct {
	PodUsecase *biz.PodUsecase
	PodService v1.PodService
}

func (p *PodHandler) AddPod(ctx context.Context, request *v1.AddPodRequest, response *v1.AddPodResponse) error {
	pod := &biz.Pod{
		Name:          request.PodInfo.Name,
		Namespace:     request.PodInfo.Namespace,
		TeamID:        request.PodInfo.TeamId,
		Replicas:      request.PodInfo.Replicas,
		CpuMax:        request.PodInfo.CpuMax,
		CpuMin:        request.PodInfo.CpuMin,
		MemoryMax:     request.PodInfo.MemoryMax,
		MemoryMin:     request.PodInfo.MemoryMin,
		PullPolicy:    request.PodInfo.PullPolicy,
		RestartPolicy: request.PodInfo.Restart_Policy,
		Type:          request.PodInfo.Type,
		Image:         request.PodInfo.Image,
	}

	pod.Ports = make([]biz.PodPort, 0, len(request.PodInfo.Ports))
	for _, item := range request.PodInfo.Ports {
		pod.Ports = append(pod.Ports, biz.PodPort{
			PodID:    item.Pod_ID,
			Port:     item.Port,
			Protocol: item.Protocol,
		})
	}

	pod.Envs = make([]biz.PodEnv, 0, len(request.PodInfo.Envs))
	for _, item := range request.PodInfo.Envs {
		pod.Envs = append(pod.Envs, biz.PodEnv{
			PodID: item.Pod_ID,
			Key:   item.Key,
			Val:   item.Val,
		})
	}

	id, err := p.PodUsecase.Add(pod)
	if err != nil {
		return err
	}

	response = &v1.AddPodResponse{
		ID: id,
	}

	return nil
}

func (p *PodHandler) DeletePod(ctx context.Context, request *v1.DeletePodRequest, response *v1.DeletePodResponse) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) FindPodByID(ctx context.Context, request *v1.FindPodByIDRequest, response *v1.FindPodByIDResponse) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) UpdatePod(ctx context.Context, request *v1.UpdatePodRequest, response *v1.UpdatePodResponse) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) FindAllPods(ctx context.Context, request *v1.FindAllPodsRequest, response *v1.FindAllPodsResponse) error {
	//TODO implement me
	panic("implement me")
}

func NewPodHandler(podUsecase *biz.PodUsecase) *PodHandler {
	return &PodHandler{PodUsecase: podUsecase}
}

