package rest

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/breadysimon/goless/jwt"
	"github.com/gertd/go-pluralize"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type RestApi struct {
	Login  bool
	err    error
	models []interface{}
	db     *gorm.DB
	svr    *http.Server
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

	name := generateResourceName(t)

	apigrp.GET("/"+name, ListHandler(s, t))
	apigrp.GET("/"+name+"/:id", ReadHandler(s.db, t))
	apigrp.POST("/"+name, CreateHandler(s, t))
	apigrp.PUT("/"+name+"/:id", UpdateHandler(s, t))
	apigrp.DELETE("/"+name+"/:id", DeleteHandler(s.db, t))

}

func (s *RestApi) CreateEndpoints(r *gin.Engine, grp string, login bool) *RestApi {
	apiGrp := r.Group(grp)
	s.Login = login
	if s.Login {
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

func (s *RestApi) Server(addr string, port int, apiRoot string, login bool) *RestApi {
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

	s.svr = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", addr, port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

func (s *RestApi) Start() *RestApi {
	if s.err == nil {
		// async
		go func() {
			log.Infof("Start server. %s", s.svr.Addr)

			if err := s.svr.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()
	}
	return s
}
func (s *RestApi) Shutdown() *RestApi {
	if s.svr != nil {
		s.svr.Shutdown(context.TODO())
	}
	return s
}
func (s *RestApi) Run() *RestApi {
	if s.err == nil {
		log.Infof("Start server. %s", s.svr.Addr)
		if err := s.svr.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
	return s
}

// generateResourceName generates a lower case, prual name for api url path,
// according to name of the struct type of the object
func generateResourceName(obj interface{}) string {
	tp := reflect.TypeOf(obj).String()
	reg := regexp.MustCompile(`.*\.`)
	name := reg.ReplaceAllString(tp, "")
	pluralize := pluralize.NewClient()
	return strings.ToLower(pluralize.Plural(name))
}
