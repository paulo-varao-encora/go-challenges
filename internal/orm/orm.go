package orm

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/go_challenges?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
