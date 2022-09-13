package main

import (
	"log"
	"net/http"
	"time"
	BlogHandle "zhanglinghua_blog/src/Blog"
	DreamHandle "zhanglinghua_blog/src/Dream"
	ImgHandle "zhanglinghua_blog/src/Img"
	LifeHandle "zhanglinghua_blog/src/Life"
	Logfile "zhanglinghua_blog/src/Logfile"
	NoteHandle "zhanglinghua_blog/src/Note"
	"zhanglinghua_blog/src/UserMessage"
	"zhanglinghua_blog/src/Util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"
var AdminUser *Util.User

func init() {
	AdminUser = Util.GetMyAdminMessage()
	log.Println("解析到的Admin信息如下", AdminUser)
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {
	// 初始化引擎
	engine := gin.Default()
	// 解决CROS问题
	engine.Use(Cors())
	// 日志的中间件
	engine.Use(Logfile.LogMiddleWare())
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
				"auth":   [2]string{"user", "admin"},
			})
		},
		// 身份认证的函数
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			log.Println(userID, password)
			if userID == AdminUser.Username && password == AdminUser.Password {
				return &User{
					UserName:  userID,
					LastName:  "Bo-Yi",
					FirstName: "Wu",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		// 验证token的函数
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Println(data)
			if v, ok := data.(*User); ok && v.UserName == AdminUser.Username {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Println(c.Params, "----", c.GetHeader("Authorization"))
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		//TokenLookup: "header: Authorization",
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Authorization",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	// 登录的接口
	engine.POST("/login", authMiddleware.LoginHandler)
	// 注册一个路由和处理函数
	Blog := engine.Group("/blog")
	{
		Blog.GET("/getBlog", BlogHandle.GetBlog)
		Blog.GET("/getAllBlog", BlogHandle.GetAllBlog)
		Blog.GET("/getBlogAllCateGory", BlogHandle.GetCategory)
		Blog.GET("/getAllTitle", BlogHandle.GetAllTitle)
		// 为需要验证权限的api 验证权限
		Blog.POST("/newBlog", authMiddleware.MiddlewareFunc(), BlogHandle.NewBlog)
		Blog.POST("/updateBlog", authMiddleware.MiddlewareFunc(), BlogHandle.ModifyBlog)
		Blog.GET("/deleteBlog", authMiddleware.MiddlewareFunc(), BlogHandle.DeleteBlog)
	}
	Note := engine.Group("/note")
	{
		Note.GET("/getNote", NoteHandle.GetNote)
		Note.GET("/getNoteInfo", NoteHandle.GetInfo)
		Note.GET("/getNoteAllCateGory", NoteHandle.GetCategory)
		// 为需要验证权限的api 验证权限
		Note.POST("/newNote", authMiddleware.MiddlewareFunc(), NoteHandle.NewNote)
		Note.POST("/updateNote", authMiddleware.MiddlewareFunc(), NoteHandle.UpdateNote)
		Note.GET("/deleteNote", authMiddleware.MiddlewareFunc(), NoteHandle.DeleteNote)
	}
	Life := engine.Group("/life")
	{
		Life.GET("/getLife", LifeHandle.GetNote)
		Life.GET("/getLifeInfo", LifeHandle.GetInfo)
		Life.GET("/getLifeAllCateGory", LifeHandle.GetCategory)
		// 为需要验证权限的api 验证权限
		Life.POST("/newLife", authMiddleware.MiddlewareFunc(), LifeHandle.NewNote)
		Life.POST("/updateLife", authMiddleware.MiddlewareFunc(), LifeHandle.UpdateNote)
		Life.GET("/deleteLife", authMiddleware.MiddlewareFunc(), LifeHandle.DeleteNote)
	}
	Dream := engine.Group("/dream")
	{
		Dream.GET("/all", DreamHandle.AllData)
		Dream.POST("/new", authMiddleware.MiddlewareFunc(), DreamHandle.NewData)
		Dream.POST("/delete", authMiddleware.MiddlewareFunc(), DreamHandle.DeleteData)
	}
	Img := engine.Group("/img")
	{
		Img.POST("/upload", authMiddleware.MiddlewareFunc(), ImgHandle.Upload)
		Img.GET("/get/:id", ImgHandle.Get)
	}
	LifeCalendar := engine.Group("/life/calendar")
	{
		LifeCalendar.GET("/get",LifeHandle.GetDayCalendar)
		LifeCalendar.GET("/delete",LifeHandle.DeleteCalendar)
		LifeCalendar.GET("/getAll",LifeHandle.GetAllCalendar)
		LifeCalendar.POST("/add",LifeHandle.AddCalendar)
		LifeCalendar.POST("/modify",LifeHandle.ModifyCalendar)
	}
	// 获取用户信息的接口
	AdminMessage := engine.Group("/user")
	{
		AdminMessage.GET("/admin", UserMessage.GetAdminMessage)
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
