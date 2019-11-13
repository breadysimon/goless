package jwt

import (
	"time"

	"github.com/breadysimon/goless/logging"
	"github.com/gin-gonic/gin"

	jwt "github.com/appleboy/gin-jwt/v2"
)

var log *logging.Logger = logging.GetLogger()

type signin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

// User demo
type User struct {
	UserId string
}

func InitAuth(r *gin.Engine) (authMiddleware *jwt.GinJWTMiddleware) {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserId,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserId: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals signin
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			return checkAuthentication(loginVals.Username, loginVals.Password)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO: replace mock
			if v, ok := data.(*User); ok && v.UserId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
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
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	return
}

func GetJwtInfo(c *gin.Context) (id string, user *User) {
	claims := jwt.ExtractClaims(c)
	data, _ := c.Get(identityKey)
	id = claims[identityKey].(string)
	user = data.(*User)
	return
}

func checkAuthentication(userID, password string) (*User, error) {
	// TODO: replace mock
	log.Debug(userID, password)
	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &User{
			UserId: userID,
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}
