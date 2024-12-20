package bootstrap

import (
	"be/models"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(config *Config) *gorm.DB {
	user := config.DBUser
	password := config.DBPass
	host := config.DBHost
	port := config.DBPort
	dbname := config.DBName

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// migrate schema
	// db.AutoMigrate(&models.AllCourses{})
	// db.AutoMigrate(&models.User{})
	// db.AutoMigrate(&models.Course{})
	// db.AutoMigrate(&models.Department{})
	// db.AutoMigrate(&models.Tuition{})
	// db.AutoMigrate(&models.RegisteredCourse{})
	// db.AutoMigrate(&models.PrerequisiteCourse{})
	// db.AutoMigrate(&models.Session{})

	if err := db.AutoMigrate(&models.AllCourses{}); err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Course{}); err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Department{}); err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Tuition{}); err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.RegisteredCourse{}); err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.PrerequisiteCourse{}); err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Session{}); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
	return db
}
