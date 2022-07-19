package MongDB

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var mgoCli *mongo.Client

func InitEngine() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb+srv://zhanglinghua:zlh20011228@blog.egnsig0.mongodb.net/?retryWrites=true&w=majority")

	// 连接到MongoDB
	mgoCli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		InitEngine()
	}
	return mgoCli
}
func GetDatabaseConnection(database string) *mongo.Database {
	return GetMgoCli().Database(database)
}
func InsertOne(database *mongo.Database, TableName string, Data interface{}) (primitive.ObjectID, error) {
	var collection = database.Collection(TableName)
	var iResult *mongo.InsertOneResult
	var err error
	if iResult, err = collection.InsertOne(context.TODO(), Data); err != nil {
		fmt.Print(err)
		return primitive.NewObjectID(), err
	}
	id := iResult.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID", id.Hex())
	return id, nil
}
func GetOne(database *mongo.Database, TableName string, Data interface{}, Projection interface{}) interface{} {
	var collection = database.Collection(TableName)
	var result bson.M
	// get the result and make the projection
	err := collection.FindOne(context.TODO(), Data, options.FindOne().SetProjection(Projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return errors.New("尚未找到匹配的对象")
		}
		return err.Error()
	} else {
		return result
	}
}

func GetAll(database *mongo.Database, TableName string, Data interface{}, Projection interface{}) ([]bson.M, error) {
	var collection = database.Collection(TableName)
	// get the result
	cursor, err := collection.Find(context.TODO(), Data, options.Find().SetProjection(Projection))
	if err != nil {
		return nil, err
	}
	var resultArray []bson.M
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		resultArray = append(resultArray, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return resultArray, nil
}
