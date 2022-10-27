package biz

import (
	"go-micro.dev/v4/util/log"
	v1 "k8s.io/api/apps/v1"
	"strconv"
)

type Pod struct {
	ID            int64     `gorm:"primary_key;not_null;auto_increment" json:"id"`
	Name          string    `gorm:"unique_index;not_null" json:"name"`
	Namespace     string    `json:"namespace"`
	TeamID        string    `json:"team_id"`
	Replicas      int32     `json:"replicas"`
	CpuMax        float32   `json:"cpu_max"`
	CpuMin        float32   `json:"cpu_min"`
	MemoryMax     float32   `json:"memory_max"`
	MemoryMin     float32   `json:"memory_min"`
	Ports         []PodPort `gorm:"Foreignkey:PodID" json:"ports"`
	Envs          []PodEnv  `gorm:"Foreignkey:PodID" json:"envs"`
	PullPolicy    string    `json:"pull_policy"`
	RestartPolicy string    `json:"restart_policy"`
	Type          string    `json:"type"`
	Image         string    `json:"image"`
}

type PodPort struct {
	ID       int64  `gorm:"primary_key;not_null;auto_increment" json:"id"`
	PodID    int64  `json:"pod_id"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

type PodEnv struct {
	ID    int64  `gorm:"primary_key;not_null;auto_increment" json:"id"`
	PodID int64  `json:"pod_id"`
	Key   string `json:"key"`
	Val   string `json:"val"`
}

type PodRepo interface {
	Add(*Pod) (int64, error)
	DeleteByID(int64) error
	FindByID(int64) (*Pod, error)
	Update(*Pod) error
	FindAll() ([]*Pod, error)
}

type K8SRepo interface {
	CreateToK8S(*Pod) error
	UpdateToK8S(*Pod) error
	DeleteFromK8S(*Pod) error
	GetPodByName(*Pod) (*v1.Deployment, error)
	SetDeployment(*Pod)
}

type PodUsecase struct {
	localPodRepo PodRepo
	k8sRepo      K8SRepo
}

func (p *PodUsecase) Add(pod *Pod) (int64, error) {
	if err := p.k8sRepo.CreateToK8S(pod); err != nil {
		log.Error(err)
		return -1, err
	}
	newID, err := p.localPodRepo.Add(pod)
	if err != nil {
		return -1, err
	}
	log.Info("创建Pod成功")
	return newID, nil
}

func (p *PodUsecase) DeleteByID(id int64) error {
	pod, err := p.FindByID(id)
	if err != nil {
		return err
	}
	if err = p.k8sRepo.DeleteFromK8S(pod); err != nil {
		return err
	}
	if err = p.localPodRepo.DeleteByID(id); err != nil {
		return err
	}
	log.Info("删除Pod ID: " + strconv.Itoa(int(pod.ID)) + " 成功")
	return nil
}

func (p *PodUsecase) FindByID(id int64) (*Pod, error) {
	return p.localPodRepo.FindByID(id)
}

func (p *PodUsecase) Update(pod *Pod) error {
	if err := p.k8sRepo.UpdateToK8S(pod); err != nil {
		return err
	}
	if err := p.localPodRepo.Update(pod); err != nil {
		return err
	}
	log.Info("更新Pod ID: " + strconv.Itoa(int(pod.ID)) + " 成功")
	return nil
}

func (p *PodUsecase) FindAll() ([]*Pod, error) {
	return p.localPodRepo.FindAll()
}

func NewPodUsecase(localPodRepo PodRepo, k8sRepo K8SRepo) *PodUsecase {
	return &PodUsecase{localPodRepo: localPodRepo, k8sRepo: k8sRepo}
}
