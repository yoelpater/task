package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
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

func handleGetDoctors(c *gin.Context) {
	limitquery := c.DefaultQuery("limit", "60")
	pagingquery := c.DefaultQuery("page", "0")
	namequery := c.DefaultQuery("name", "")
	genderquery := c.DefaultQuery("gender", "")
	servicequery := c.DefaultQuery("servicerole", "")

	limit, err2 := strconv.ParseInt(limitquery, 10, 64)
	if err2 != nil {
		// handle error
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something is wrong with limit value"})
		return
	}
	page, err3 := strconv.ParseInt(pagingquery, 10, 64)
	if err3 != nil {
		// handle error
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something is wrong with page value"})
		return
	}
	var loadedDoctors, totalpage, err = GetAllDoctors(limit, page, namequery, servicequery, genderquery)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pages": totalpage, "doctors": loadedDoctors})
}

func handleGetOneDoctor(c *gin.Context) {
	var doctor Doctor
	if err := c.BindUri(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	doctorid := c.Param("id")
	var loadedDoctor, err = GetDoctorByID(doctorid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loadedDoctor)
}

func handleCreateDoctor(c *gin.Context) {
	var doctor Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	id, err := CreateDoctor(&doctor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func handleUpdateDoctor(c *gin.Context) {
	var doctor Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	savedTask, err := UpdateDoctorbyID(&doctor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": savedTask})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.GET("/", hello)
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/tasks/:id", handleGetTask)
	r.GET("/tasks/", handleGetTasks)
	r.PUT("/tasks/", handleUpdateTask)
	r.POST("/tasks/", handleCreateTask)

	r.GET("/doctors/", handleGetDoctors)
	r.POST("/doctors/", handleCreateDoctor)
	r.GET("/doctors/:id", handleGetOneDoctor)
	r.PUT("/doctors/", handleUpdateDoctor)
	return r
}

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	r := setupRouter()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
