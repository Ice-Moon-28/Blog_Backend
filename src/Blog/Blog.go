package Blog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"zhanglinghua_blog/src/MongDB"
	"zhanglinghua_blog/src/Util"
)

type BlogContent struct {
	id      string
	content string
}

var blog *mongo.Database

func init() {
	blog = MongDB.GetDatabaseConnection("blog")
}
func GetBlog(context *gin.Context) {
	id := context.Query("id")
	// if we want to make Query according to _id,must do this convertion
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		context.String(http.StatusOK, err.Error())
		return
	}
	result, _ := MongDB.GetOne(blog, "blog", bson.D{{"_id", _id}}, bson.D{})
	fmt.Println(result)
	// use this method return json obj
	context.JSON(http.StatusOK, result)
	// use this method return string context.string()
}

type BlogCategory struct {
	Color    string
	Category string
}

func NewBlog(context *gin.Context) {
	json := make(map[string]interface{})
	context.BindJSON(&json)
	// 插入blog数据
	_, err := MongDB.InsertOne(blog, "blog", bson.D{{"markdown", json["markdown"]},
		{"title", json["title"]},
		{"category", json["category"]},
		{"time", json["time"]},
	})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	}
	var nowCategoryArray []string
	result, err := MongDB.GetAll(blog, "blogcategory", bson.D{}, bson.D{}, &options.FindOptions{})
	if err != nil {
		log.Println("NewBlog_获取Category出错", err.Error())
	}
	for _, v := range result {
		midstr, _ := v["category"].(string)
		nowCategoryArray = append(nowCategoryArray, midstr)
	}
	categoryArray := json["category"].([]interface{})
	fmt.Println(result, json["category"], nowCategoryArray, categoryArray, "result--")
	for _, v := range categoryArray {
		mid := v.(map[string]interface{})
		midstr := mid["category"].(string)
		midColorStr := mid["color"].(string)
		if !Util.ArrayHasValue(midstr, nowCategoryArray) {
			// 更新目录信息
			_, err = MongDB.InsertOne(blog, "blogcategory", bson.D{{"category", midstr}, {"color", midColorStr}})
			if err != nil {
				context.String(http.StatusInternalServerError, err.Error())
			}
		}
	}
}

// 获取所有的Blog
func GetAllBlog(context *gin.Context) {
	page := context.Query("page")
	tip := context.Query("tip")
	fmt.Println("tip", tip, page)
	var LimitNumber int64 = 10
	var DataPrejection bson.D = bson.D{{}}
	pageSkip, ok := strconv.ParseInt(page, 10, 64)
	if ok != nil {
		context.String(http.StatusInternalServerError, ok.Error())
	}
	if tip != "" {
		// 这种方式来处理对象数组的对于具体对象的查询
		DataPrejection = bson.D{{"category.category", tip}}
	}
	// 需要跳过的数据条数
	pageSkip = (pageSkip - 1) * 10
	result, err := MongDB.GetAll(blog, "blog", DataPrejection, bson.D{}, &options.FindOptions{
		Skip:  &pageSkip,
		Sort:  bson.D{{"_id", 1}},
		Limit: &LimitNumber,
	})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	} else {
		fmt.Println(result)
		// Json 只能处理大写的对象属性
		context.JSON(http.StatusOK, result)
	}
}
func GetCategory(context *gin.Context) {
	result, err := MongDB.GetAll(blog, "blogcategory", bson.D{}, bson.D{{"category", 1}, {"color", 1}}, &options.FindOptions{})
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
	} else {
		context.JSON(http.StatusOK, result)
	}
}
