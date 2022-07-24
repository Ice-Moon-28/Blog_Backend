package Note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	fmt.Println(result)
	// use this method return json obj
	context.JSON(http.StatusOK, result)
	// use this method return string context.string()
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

func GetInfo(context *gin.Context) {
	result, err := MongDB.GetAll(note, "notecategory", bson.D{}, bson.D{{"category", 1}, {"id", 1}, {"title", 1}, {"_id", 1}})
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
			fmt.Println(v, str, id, title, ok, ValueIndex)
			if ok && ValueIndex == -1 {
				titleArray = append(titleArray, str)
				resultArray = append(resultArray, info{Title: str, TitleUrl: id, List: []struct {
					Label string
					Url   primitive.ObjectID
				}{}})
			} else {

				resultArray[ValueIndex].List = append(resultArray[ValueIndex].List, struct {
					Label string
					Url   primitive.ObjectID
				}{Label: title, Url: id})
			}
		}
		fmt.Println(resultArray)
		// Json 只能处理大写的对象属性
		context.JSON(http.StatusOK, resultArray)
	}
}
func GetCategory(context *gin.Context) {
	result, err := MongDB.GetAll(note, "notecategory", bson.D{}, bson.D{{"category", 1}, {"_id", 1}})
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
