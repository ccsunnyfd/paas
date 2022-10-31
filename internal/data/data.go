package data

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"k8s.io/client-go/kubernetes"
	"saas/governance/config"
)

func NewMysql(user, pwd string, database string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(
		user+":"+pwd+"@/"+database+"?charset=utf8&parseTime=True&loc=Local"),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)
}

func NewK8SClientset(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := config.GetK8SFromCMD(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
