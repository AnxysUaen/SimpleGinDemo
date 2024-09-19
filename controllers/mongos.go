package controllers

import (
	"SimpleGinDemo/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client // 全局Mongo客户端

func initMongoClient() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

// 接收的Log参数
type InsertData struct {
	Name      string    `bson:"Name"`
	RequestID string    `bson:"RequestID"`
	Type      string    `bson:"Type"`
	JSONData  bson.M    `bson:"JSONData"`
	Time      time.Time `bson:"Time"`
}

func SaveLog(c *gin.Context) {
	if client == nil {
		fmt.Println("Initializing MongoDB client...")
		initMongoClient()
	} else {
		// 尝试ping MongoDB以检查连接是否仍然活跃
		err := client.Ping(context.TODO(), nil)
		if err != nil {
			fmt.Println("Lost connection to MongoDB, reconnecting...")
			initMongoClient()
		}
	}

	var requestData models.RequestData
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据requestData.DataSource选择合适的集合
	collection := client.Database("bmp_m_logs").Collection(requestData.DataSource)

	var result map[string]interface{}
	err := json.Unmarshal([]byte(requestData.JSONData), &result)
	if err != nil {
		log.Fatal(err)
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", requestData.Time)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}

	insertData := InsertData{
		RequestID: requestData.RequestID,
		Name:      requestData.Name,
		JSONData:  bson.M(result),
		Time:      parsedTime,
		Type:      requestData.Type,
	}

	// 插入文档
	insertResult, err := collection.InsertOne(context.TODO(), insertData)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert log"})
		return
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	c.JSON(http.StatusOK, gin.H{"message": "Log inserted successfully"})
}
