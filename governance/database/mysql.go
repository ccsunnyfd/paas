package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetMysql(user, pwd string, database string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(
		user+":"+pwd+"@/"+database+"?charset=utf8&parseTime=True&loc=Local"),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)
}
