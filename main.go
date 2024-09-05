package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID          uint       `gorm:"primaryKey;->" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `gorm:"<-:create" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"<-" json:"updated_at"`
}

var DB *gorm.DB

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("DB_PORT"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if DB.AutoMigrate(&Task{}) != nil {
		log.Fatal("Failed to migrate!")
	}
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.AbortWithStatus(400)
		return
	}
	DB.Create(&task)
	c.JSON(201, task)
}

func getTasks(c *gin.Context) {
	var tasks []Task
	DB.Find(&tasks)
	c.JSON(200, tasks)
}

func getTaskById(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if result := DB.First(&task, id); result.Error != nil {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, task)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if result := DB.First(&task, id); result.Error != nil {
		c.AbortWithStatus(404)
		return
	}
	var updateData Task
	if c.BindJSON(&updateData) != nil {
		c.AbortWithStatus(400)
		return
	}

	DB.Model(&task).Updates(updateData)
	c.JSON(200, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if result := DB.First(&task, id); result.Error != nil {
		c.AbortWithStatus(404)
		return
	}
	DB.Delete(&task)
	c.Status(204)
}

func main() {

	router := gin.Default()

	router.POST("/tasks", createTask)
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskById)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)

	if router.Run("localhost:8080") != nil {
		log.Fatal("Error!")
	}
}
