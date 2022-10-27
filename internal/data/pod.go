package data

import (
	"gorm.io/gorm"
	"saas/internal/biz"
)

type podRepo struct {
	db *gorm.DB
}

func NewPodRepo(db *gorm.DB) biz.PodRepo {
	return &podRepo{db: db}
}

func (p *podRepo) InitTable() error {
	return p.db.AutoMigrate(&biz.Pod{}, &biz.PodPort{}, &biz.Pod{})
}

func (p *podRepo) Add(podInfo *biz.Pod) (int64, error) {
	return podInfo.ID, p.db.Create(podInfo).Error
}

func (p *podRepo) DeleteByID(id int64) error {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	if err := p.db.Where("id = ?", id).Delete(&biz.Pod{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := p.db.Where("pod_id = ?", id).Delete(&biz.PodEnv{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := p.db.Where("pod_id = ?", id).Delete(&biz.PodPort{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (p *podRepo) FindByID(id int64) (*biz.Pod, error) {
	pod := &biz.Pod{}
	return pod, p.db.Preload("Ports").Preload("Envs").First(pod, id).Error
}

func (p *podRepo) Update(podInfo *biz.Pod) error {
	return p.db.Model(podInfo).Updates(podInfo).Error
}

func (p *podRepo) FindAll() (podList []*biz.Pod, err error) {
	return podList, p.db.Find(podList).Error
}
