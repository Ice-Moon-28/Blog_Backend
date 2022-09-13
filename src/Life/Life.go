package Note

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"zhanglinghua_blog/src/MongDB"
	"zhanglinghua_blog/src/Util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	result, _ := MongDB.GetOne(note, "life", bson.D{{"_id", _id}}, bson.D{})
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
	_, err = MongDB.UpDateOne(note, "life", bson.D{{"$set", bson.D{{"markdown", json["markdown"]}}},
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
	id, err := MongDB.InsertOne(note, "life", bson.D{{"markdown", json["markdown"]},
		{"title", json["title"]},
		{"category", json["category"]}})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	}
	// 更新目录信息
	id, err = MongDB.InsertOne(note, "lifecategory", bson.D{{"id", id}, {"title", json["title"]}, {"category", json["category"]}})
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
	succuess, err := MongDB.DeleteOne(note, "life", bson.D{{"_id", _id}})
	succuess1, err := MongDB.DeleteOne(note, "lifecategory", bson.D{{"id", _id}})
	// use this method return json obj
	if succuess && succuess1 {
		context.String(http.StatusOK, "删除成功")
	} else {
		context.String(http.StatusInternalServerError, err.Error())
	}
	// use this method return string context.string()
}

func GetInfo(context *gin.Context) {
	result, err := MongDB.GetAll(note, "lifecategory", bson.D{}, bson.D{{"category", 1}, {"id", 1}, {"title", 1}, {"_id", 1}}, &options.FindOptions{})
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
	result, err := MongDB.GetAll(note, "lifecategory", bson.D{}, bson.D{{"category", 1}, {"_id", 1}}, &options.FindOptions{})
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
// month day schedule {date , content}
func AddCalendar(context *gin.Context) {
	json := make(map[string]interface{})
	context.BindJSON(&json)
	month,_ := json["month"].(float64)
	day,_ := json["day"].(float64)
	log.Println(json,json["month"],json["day"],"---")
	res, err := MongDB.GetOne(note,"calendar",bson.D{
		{"month",month},
		{"day",day}},
		bson.D{})
	log.Println(res["_id"], err,"err---")
	if err != nil{
		// 插入blog数据
		_, err2 := MongDB.InsertOne(note, "calendar", bson.D{{"month", month},
			{"day", day},
			{"schedule", json["schedule"]},
		})
		if err2 == nil{
			context.String(http.StatusOK,"")
		}else{
			context.String(500,err2.Error())
		}
	}else{
		log.Println(json["schedule"],res["_id"])
		success, _ := MongDB.UpDateOne(note,"calendar",bson.D{ 
			{"$set",bson.D{{"schedule",json["schedule"]}}}},bson.D{{"_id",res["_id"]}})
		if success{
			context.String(200,"succuess")
		}else{
			context.String(500,err.Error())
		}
	}

}
// 获取每个月的所有date信息
func GetDayCalendar(context *gin.Context) {
	month , err:= strconv.ParseFloat( context.Query("month") , 64)
	if err != nil{
		context.String(400,err.Error())
	}
	day , err1:= strconv.ParseFloat( context.Query("day"),64)
	if err1 != nil{
		context.String(400,err1.Error())
	}
	log.Println(month,day,"--")
	result , err := MongDB.GetOne(note,"calendar",bson.D{{"month",month},{"day",day}},bson.D{})
	context.JSON(http.StatusOK,result)
}
// 获取每个月的所有date信息
func GetAllCalendar(context *gin.Context) {
	month , err:= strconv.ParseFloat( context.Query("month") , 64)
	if err != nil{
		context.String(400,err.Error())
	}
	result , _ := MongDB.GetAll(note,"calendar",bson.D{{"month",month}},bson.D{{"day",1}},&options.FindOptions{})
	context.JSON(http.StatusOK,result)
}
func DeleteCalendar(context *gin.Context) {
	month ,_:= strconv.ParseFloat( context.Query("month"),64)
	day,_ := strconv.ParseFloat( context.Query("day"),64)
	log.Println(day,month,"删除该事务")
	res, err := MongDB.GetOne(note,"calendar",bson.D{
		{"month",month},
		{"day",day}},
		bson.D{})
	success,err := MongDB.DeleteOne(note,"calendar",bson.D{{"_id",res["_id"]}})
	if success{
		context.JSON(http.StatusOK,"success")
	}else{
		context.String(500,err.Error())
	}
}
func ModifyCalendar(context *gin.Context){
	json := make(map[string]interface{})
	context.BindJSON(&json)
	res, err := MongDB.GetOne(note,"calendar",bson.D{
		{"month",json["month"]},
		{"day",json["day"]}},
		bson.D{})
	if err == nil{
		log.Println(json["schedule"],res["_id"])
		success, _ := MongDB.UpDateOne(note,"calendar",bson.D{ 
			{"$set",bson.D{{"schedule",json["schedule"]}}}},bson.D{{"_id",res["_id"]}})
		if success{
			context.String(200,"succuess")
		}else{
			context.String(500,err.Error())
		}
	}else{
		context.String(500,err.Error())
	}
}