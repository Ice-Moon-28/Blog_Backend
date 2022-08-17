package Note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"zhanglinghua_blog/src/MongDB"
	"zhanglinghua_blog/src/Util"
)

var note *mongo.Database

func init() {
	note = MongDB.GetDatabaseConnection("blog")
}
func GetNote(context *gin.Context) {
	id := context.Query("id")
	// if we want to make Query according to _id,must do this convertion
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		context.String(http.StatusOK, err.Error())
		return
	}
	result, _ := MongDB.GetOne(note, "note", bson.D{{"_id", _id}}, bson.D{})
	// use this method return json obj
	context.JSON(http.StatusOK, result)
	// use this method return string context.string()
}
func UpdateNote(context *gin.Context) {
	var err error
	var id primitive.ObjectID
	json := make(map[string]string)
	context.BindJSON(&json)
	// 插入blog数据
	fmt.Println("upadate---", json)
	id, err = primitive.ObjectIDFromHex(json["_id"])
	_, err = MongDB.UpDateOne(note, "note", bson.D{{"$set", bson.D{{"markdown", json["markdown"]}}},
		{"$set", bson.D{{"title", json["title"]}}},
		{"$set", bson.D{{"category", json["category"]}}},
		{"$set", bson.D{{"_id", id}}},
	}, bson.D{{"_id", id}})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	} else {
		context.String(200, "更新成功")
	}
}

func NewNote(context *gin.Context) {
	json := make(map[string]string)
	context.BindJSON(&json)
	// 插入blog数据
	id, err := MongDB.InsertOne(note, "note", bson.D{{"markdown", json["markdown"]},
		{"title", json["title"]},
		{"category", json["category"]}})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	}
	// 更新目录信息
	id, err = MongDB.InsertOne(note, "notecategory", bson.D{{"id", id}, {"title", json["title"]}, {"category", json["category"]}})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	}
}

type info struct {
	Title    string
	TitleUrl primitive.ObjectID
	List     []struct {
		Label string
		Url   primitive.ObjectID
	}
}

func DeleteNote(context *gin.Context) {
	id := context.Query("_id")
	// if we want to make Query according to _id,must do this convertion
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		context.String(http.StatusOK, err.Error())
		return
	}
	succuess, err := MongDB.DeleteOne(note, "note", bson.D{{"_id", _id}})
	succuess1, err := MongDB.DeleteOne(note, "notecategory", bson.D{{"id", _id}})
	// use this method return json obj
	if succuess && succuess1 {
		context.String(http.StatusOK, "删除成功")
	} else {
		context.String(http.StatusInternalServerError, err.Error())
	}
	// use this method return string context.string()
}

func GetInfo(context *gin.Context) {
	result, err := MongDB.GetAll(note, "notecategory", bson.D{}, bson.D{{"category", 1}, {"id", 1}, {"title", 1}, {"_id", 1}}, &options.FindOptions{})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	} else {
		var resultArray []info
		var titleArray []string
		for _, v := range result {
			str, ok := v["category"].(string)
			title, _ := v["title"].(string)
			id, _ := v["id"].(primitive.ObjectID)
			ValueIndex := Util.GetValueIndexInArray(str, titleArray)
			if ok && ValueIndex == -1 {
				titleArray = append(titleArray, str)
				resultArray = append(resultArray, info{Title: str, TitleUrl: id, List: []struct {
					Label string
					Url   primitive.ObjectID
				}{{Label: title, Url: id}}})
			} else {
				resultArray[ValueIndex].List = append(resultArray[ValueIndex].List, struct {
					Label string
					Url   primitive.ObjectID
				}{Label: title, Url: id})
			}
		}
		// Json 只能处理大写的对象属性
		context.JSON(http.StatusOK, resultArray)
	}
}
func GetCategory(context *gin.Context) {
	result, err := MongDB.GetAll(note, "notecategory", bson.D{}, bson.D{{"category", 1}, {"_id", 1}}, &options.FindOptions{})
	var categoryArray []string
	for i := range result {
		// 类型断言
		str, ok := result[i]["category"].(string)
		if ok && !Util.ArrayHasValue(str, categoryArray) {
			categoryArray = append(categoryArray, str)
		}
	}
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	} else {
		context.JSON(http.StatusOK, categoryArray)
	}
}
