package Dream

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"zhanglinghua_blog/src/MongDB"
)

var dream *mongo.Database

func init() {
	dream = MongDB.GetDatabaseConnection("blog")
}

func AllData(context *gin.Context) {
	result, err := MongDB.GetOne(dream, "dream", bson.D{}, bson.D{{"data", 1}, {"begin-time", 1}})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	} else {
		// Json 只能处理大写的对象属性
		context.JSON(http.StatusOK, result)
	}
}
func NewData(context *gin.Context) {
	json := make(map[string]interface{})
	context.BindJSON(&json)
	id, ok := json["_id"].(string)
	Id, err := primitive.ObjectIDFromHex(id)
	if ok && err == nil {
		MongDB.UpDateOne(dream, "dream", bson.D{{"data", json["Data"]}}, bson.D{{"_id", Id}})
		context.JSON(http.StatusOK, gin.H{"data": json["Data"]})
	} else {
		context.String(400, "当前传递的参数有问题")
	}
}

type DataMessage struct {
	Data []struct{ todo string }
}

func DeleteData(context *gin.Context) {
	json := make(map[string]interface{})
	context.BindJSON(&json)
	id, ok := json["_id"].(string)
	Id, err := primitive.ObjectIDFromHex(id)
	if ok && err == nil {
		MongDB.UpDateOne(dream, "dream", bson.D{{"data", json["Data"]}}, bson.D{{"_id", Id}})
		context.JSON(http.StatusOK, gin.H{"data": json["Data"]})
	} else {
		context.String(400, "当前传递的参数有问题")
	}
}
