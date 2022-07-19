package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	BlogHandle "zhanglinghua_blog/src/Blog"
	NoteHandle "zhanglinghua_blog/src/Note"
)

func main() {
	// 初始化引擎
	engine := gin.Default()
	// 解决CROS问题
	engine.Use(Cors())
	// 注册一个路由和处理函数
	Blog := engine.Group("/blog")
	{
		Blog.GET("/getBlog", BlogHandle.GetBlog)
		Blog.GET("/getAllBlog", BlogHandle.GetAllBlog)
		Blog.GET("/getBlogAllCateGory", BlogHandle.GetCategory)
		Blog.POST("/newBlog", BlogHandle.NewBlog)
	}
	Note := engine.Group("/note")
	{
		Note.GET("/getNote", NoteHandle.GetNote)
		Note.GET("/getNoteInfo", NoteHandle.GetInfo)
		Note.GET("/getNoteAllCateGory", NoteHandle.GetCategory)
		Note.POST("/newNote", NoteHandle.NewNote)
	}
	// 绑定端口，然后启动应用
	engine.Run(":9205")
}

/**
* 根请求处理函数
* 所有本次请求相关的方法都在 context 中，完美
* 输出响应 hello, world
 */
func WebRoot(context *gin.Context) {
	context.String(http.StatusOK, "hello, world")
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "*")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
