package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/breadysimon/goless/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type RestApi struct {
	login  bool
	err    error
	models []interface{}
	db     *gorm.DB
}

func NewRestApi(t ...interface{}) *RestApi {
	s := &RestApi{}
	for _, v := range t {
		s.models = append(s.models, v)
	}
	return s
}
func (s *RestApi) Error() error {
	return s.err
}
func (s *RestApi) Connect(dialect, conn string) *RestApi {
	if s.err == nil {
		if s.db, s.err = gorm.Open(dialect, conn); s.err == nil {
			s.db.SingularTable(true)
			s.db.LogMode(true)
			s.db.DB().SetMaxIdleConns(10)
			s.db.DB().SetMaxOpenConns(100)

			//AutoMigrate##
			s.db.AutoMigrate(s.models...)
		}
	}
	return s
}
func (s *RestApi) CloseDb() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *RestApi) makeRouters(apigrp *gin.RouterGroup, t interface{}) {

	name := guessResourceName(t)

	apigrp.GET("/"+name, ListHandler(s.db, t))
	apigrp.GET("/"+name+"/:id", ReadHandler(s.db, t))
	apigrp.POST("/"+name, CreateHandler(s.db, t))
	apigrp.PUT("/"+name+"/:id", UpdateHandler(s.db, t))
	apigrp.DELETE("/"+name+"/:id", DeleteHandler(s.db, t))

}

func (s *RestApi) CreateEndpoints(r *gin.Engine, grp string, login bool) *RestApi {
	apiGrp := r.Group(grp)

	if login {
		authMiddleware := jwt.InitAuth(r)
		apiGrp.Use(authMiddleware.MiddlewareFunc())
		r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
			// claims := jwt.ExtractClaims(c)
			// log.Printf("NoRoute claims: %#v\n", claims)
			c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		})
	}

	for _, v := range s.models {
		s.makeRouters(apiGrp, v)
	}
	return s
}

func (s *RestApi) Serve(addr string, port int, apiRoot string, login bool) *http.Server {
	r := gin.New()
	{
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
		// gin.SetMode(gin.TestMode)

	}
	// CORS control
	{
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AllowCredentials = true
		config.AddExposeHeaders("Content-Range")
		r.Use(cors.New(config))
	}

	s.CreateEndpoints(r, apiRoot, login)

	svr := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", addr, port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Infof("Start server. Listening at :%d", port)

	go func() {
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	return svr
}
