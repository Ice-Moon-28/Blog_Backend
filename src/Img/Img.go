package Img

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"zhanglinghua_blog/src/MongDB"
	"zhanglinghua_blog/src/Util"
)

var Img = MongDB.GetDatabaseConnection("blog")

func ForMatString(fileName string) string {
	return "![](http://" + Util.GetMyAdminMessage().WebSite + "/img/get/" + fileName + ")"
}

// 上传图片文件
func Upload(context *gin.Context) {
	err := context.Request.ParseMultipartForm(200000)
	if err != nil {
		log.Fatal(err)
	}
	var fileNameArray []string
	// 获取表单
	form := context.Request.MultipartForm
	// 获取参数upload后面的多个文件名，存放到数组files里面，
	files := form.File["file"]
	for i, _ := range files {
		var file multipart.File
		file, err = files[i].Open()
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		fileName := files[i].Filename
		fileContent, _ := ioutil.ReadAll(file)
		// 获取对应的字符串id
		id := fmt.Sprintf("%x", Util.GetFileHash256([]byte(fileName)))
		// 先进行删除 再进行添加图片 这样就可以实现同名图片覆盖的效果
		//MongDB.GridfsDelete("image", fileName)
		err = MongDB.GridfsUploadWithID("image", id, fileName, fileContent)
		if err != nil {
			break
		}
		fileNameArray = append(fileNameArray, ForMatString(id))
	}
	if err == nil {
		context.String(http.StatusCreated, strings.Join(fileNameArray, "@@@"))
	} else {
		context.String(400, err.Error())
	}
}

// 获取图片文件
func Get(context *gin.Context) {
	id := context.Param("id")
	result, err := MongDB.GridfsDownload("image", id)
	if err != nil {
		context.String(400, err.Error())
	} else {
		context.Writer.WriteString(string(result))
	}
}
