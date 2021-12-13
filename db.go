package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout           = 10
	connectionStringTemplate = "mongodb+srv://%s:%s@%s"
)

// GetConnection Retrieves a client to the MongoDB
func getConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	clusterEndpoint := os.Getenv("MONGODB_ENDPOINT")

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client, ctx, cancel
}

// GetAllTasks Retrives all tasks from the db
func GetAllTasks() ([]*Task, error) {
	var tasks []*Task

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("tasks")
	collection := db.Collection("tasks")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &tasks); err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return tasks, nil
}

// GetTaskByID Retrives a task by its id from the db
func GetTaskByID(id string) (*Task, error) {
	var task *Task
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("tasks")
	collection := db.Collection("tasks")
	result := collection.FindOne(ctx, bson.M{"_id": objectId})
	if result == nil {
		return nil, errors.New("could not find a task")
	}

	err2 := result.Decode(&task)

	if err2 != nil {
		log.Printf("Failed marshalling %v", err2)
		return nil, err2
	}
	log.Printf("Tasks: %v", task)
	return task, nil
}

//Create creating a task in a mongo
func Create(task *Task) (primitive.ObjectID, error) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	result, err := client.Database("tasks").Collection("tasks").InsertOne(ctx, task)
	if err != nil {
		log.Printf("Could not create Task: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

//Update updating an existing task in a mongo
func Update(task *Task) (*Task, error) {
	var updatedTask *Task
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": task,
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database("tasks").Collection("tasks").FindOneAndUpdate(ctx, bson.M{"_id": task.ID}, update, &opt).Decode(&updatedTask)
	if err != nil {
		log.Printf("Could not save Task: %v", err)
		return nil, err
	}
	return updatedTask, nil
}

//DOCTOR DB HANDLER
func GetAllDoctors(limit int64, page int64, name string, service string, gender string) ([]*Doctor, int64, error) {
	var doctors []*Doctor

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	options := options.Find()
	options.SetLimit(limit)
	if page > 0 {
		options.SetSkip((page - 1) * limit)
	} else {
		options.SetSkip(0)
	}
	options.SetSort(bson.D{primitive.E{Key: "rating", Value: -1}})
	filter := bson.M{}
	if name != "" {
		//filter = bson.M{"title": params.Title}
		filter["$text"] = bson.M{
			"$search": name,
		}
	}
	if service != "" {
		filter["servicerole"] = service
	}
	if gender != "" {
		filter["gender"] = gender
	}

	db := client.Database("doctors")
	collection := db.Collection("doctors")

	mod := mongo.IndexModel{
		Keys: bson.M{
			"textindex": "text",
		}, Options: nil,
	}

	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return nil, 0, err
	}

	totalcount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	totalpage := int64(math.Ceil(float64(totalcount) / float64(limit)))

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &doctors); err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, 0, err
	}
	return doctors, totalpage, nil
}

func GetDoctorByID(id string) (*Doctor, error) {
	var doctor *Doctor
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("doctors")
	collection := db.Collection("doctors")
	result := collection.FindOne(ctx, bson.M{"_id": objectId})
	if result == nil {
		return nil, errors.New("could not find a task")
	}

	err = result.Decode(&doctor)

	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	log.Printf("Tasks: %v", doctor)
	return doctor, nil
}

func CreateDoctor(doctor *Doctor) (primitive.ObjectID, error) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	result, err := client.Database("doctors").Collection("doctors").InsertOne(ctx, doctor)
	if err != nil {
		log.Printf("Could not create Task: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func UpdateDoctorbyID(doctor *Doctor) (*Doctor, error) {
	var updatedDoctor *Doctor
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": doctor,
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database("doctors").Collection("doctors").FindOneAndUpdate(ctx, bson.M{"_id": doctor.ID}, update, &opt).Decode(&updatedDoctor)
	if err != nil {
		log.Printf("Could not save Task: %v", err)
		return nil, err
	}
	return updatedDoctor, nil
}
