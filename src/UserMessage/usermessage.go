package UserMessage

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"zhanglinghua_blog/src/MongDB"
)

var note *mongo.Database

func init() {
	note = MongDB.GetDatabaseConnection("blog")
}

func GetAdminMessage(context *gin.Context) {
	var err error
	noteNumber, err := MongDB.CountDoc(note, "note", bson.D{{}})
	if err != nil {
		context.String(500, err.Error())
	}
	blogNumber, err := MongDB.CountDoc(note, "blog", bson.D{{}})
	if err != nil {
		context.String(500, err.Error())
	}
	message, err := MongDB.GetOne(note, "AdminMessage", bson.D{{}}, bson.D{
		{"note", 1},
		{"AliPay", 1}, {"QQ", 1},
		{"Wechat", 1}, {"Github", 1}, {"portrait", 1},
	})
	if err != nil {
		context.String(500, err.Error())
	} else {
		context.JSON(200, bson.M{"blog": blogNumber, "note": noteNumber,
			"AliPay": message["AliPay"], "QQ": message["QQ"],
			"Wechat": message["Wechat"], "Github": message["Github"], "portrait": message["portrait"],
		})
	}
}
