package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Healthy"})
}

func handleGetTasks(c *gin.Context) {
	var loadedTasks, err = GetAllTasks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": loadedTasks})
}

func handleGetTask(c *gin.Context) {
	var task Task
	if err := c.BindUri(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	taskid := c.Param("id")
	var loadedTask, err = GetTaskByID(taskid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": loadedTask.ID, "Title": loadedTask.Title, "Body": loadedTask.Body})
}

func handleCreateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	id, err := Create(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func handleUpdateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	savedTask, err := Update(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": savedTask})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", hello)
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/tasks/:id", handleGetTask)
	r.GET("/tasks/", handleGetTasks)
	r.PUT("/tasks/", handleUpdateTask)
	r.POST("/tasks/", handleCreateTask)
	return r
}

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	fmt.Printf("mongodb %s mongodb\n", os.Getenv("MONGODB_USERNAME"))
	fmt.Printf("mongodb %s mongodb\n", os.Getenv("MONGODB_ENDPOINT"))
	r := setupRouter()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
