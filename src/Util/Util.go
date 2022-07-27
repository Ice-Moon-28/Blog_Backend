package Util

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func ArrayHasValue[T int | string | float32 | float64 | struct{}](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
func GetValueIndexInArray[T int | string | float32 | float64](value T, array []T) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}

// 读取对应Json 文件中的 v对象 v传入地址
func Load(filename string, v interface{}) error {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal([]byte(data), v)
	if err != nil {
		return err
	}
	return nil
}

func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 管理员的数据类型
type User struct {
	Username string
	Password string
	DataBase string
	WebSite  string
}

//获取管理员对应的信息 以及对应的数据库链接信息
func GetMyAdminMessage() *User {
	hasJson, err := IsExists("admin.json")
	if err != nil {
		fmt.Println(err)
	}
	if hasJson {
		var v = &User{}
		err := Load("admin.json", &v)
		if err != nil {
			panic("Load ip.Json" + err.Error())
		}
		return v
	} else {
		file, _ := os.OpenFile("admin.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		file.WriteString("{\n \"Username\" : \"\" ,\n \"Password\": \"\",\n \"DataBase\": \"\"\n \"DataBase\": \"\"\n}")
		defer file.Close()
		panic(errors.New("请按照对应的管理员信息以及MongDB数据库以及网站DNS的连接信息!"))
	}
}

// 获取文件的hash值
func GetFileHash256(fileName []byte) [32]byte {
	currentTime := time.Now().String() //获取当前时间，类型是Go的时间类型Time
	var date1 []byte = append([]byte(currentTime), fileName...)
	var hs = sha256.Sum256(date1)
	return hs
}
