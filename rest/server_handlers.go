package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/breadysimon/goless/jwt"
	"github.com/breadysimon/goless/reflection"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ModelContext struct {
	Login bool
	User  *jwt.User
}

type CreateHandlerPreprocessor interface {
	BeforeDbCreate(s *RestApi, c *gin.Context) error
}
type UpdateHandlerPreprocessor interface {
	BeforeDbUpdate(s *RestApi, c *gin.Context) error
}
type K8sResource interface {
	K8sCreate() error
	K8sRead() error
	K8sList(o interface{}, filter map[string]interface{}) int
	K8sDelete() error
}

func responseJson(c *gin.Context, data interface{}, err error, code int) {
	if err == nil {
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    code,
			"message": err.Error(),
		})
	}
}

func ReadHandler(db *gorm.DB, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		m := reflection.MakeInstance(t)
		var err error
		var code int
		if r, ok := m.(K8sResource); ok {
			code = 40001
			err = r.K8sRead()
		} else {
			code = 40002
			id, _ := strconv.Atoi(c.Param("id"))
			err = mdlRead(db, m, id)
		}
		responseJson(c, m, err, code)
	}
}

func CreateHandler(s *RestApi, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		data := reflection.MakeInstance(t)
		code := 10001
		err := c.BindJSON(data)

		d, ok := data.(CreateHandlerPreprocessor)
		if err == nil && ok {
			code = 10002
			err = d.BeforeDbCreate(s, c)
		}

		if err == nil {
			if r, ok := data.(K8sResource); ok {
				code = 10003
				err = r.K8sCreate()
			} else {
				code = 10004
				err = mdlCreate(s.db, data)
			}
		}
		responseJson(c, data, err, code)
	}
}
func UpdateHandler(s *RestApi, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		code := 20001
		id, err := strconv.Atoi(c.Param("id"))
		data := reflection.MakeInstance(t)
		if err == nil {
			code = 20002
			err = c.BindJSON(&data)
		}
		d, ok := data.(UpdateHandlerPreprocessor)
		if err == nil && ok {
			code = 20002
			err = d.BeforeDbUpdate(s, c)
		}
		if err == nil {
			reflection.SetIdInt(data, id)
			code = 20003
			err = mdlUpdate(s.db, data)
		}
		responseJson(c, data, err, code)
	}
}

func DeleteHandler(db *gorm.DB, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		code := 30001
		id, err := strconv.Atoi(c.Param("id"))
		m := reflection.MakeInstance(t)
		if err == nil {
			if r, ok := t.(K8sResource); ok {
				code = 30003
				err = c.BindJSON(&r)
				if err == nil {
					code = 30004
					err = r.K8sDelete()
				}
			} else {
				reflection.SetIdInt(m, id)
				code = 30002
				err = mdlDelete(db, m)
			}
		}
		responseJson(c, nil, err, code)
	}
}

func ListHandler(s *RestApi, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		m := reflection.MakeInstance(t)
		data := reflection.MakeSlice(t)
		sort, offset, limit, filter := ParsedQuery(c)
		var n int
		if r, ok := m.(K8sResource); ok {
			n = r.K8sList(data, filter)
		} else {
			n = mdlList(s.db, m, data, sort, offset, limit, filter)

		}
		c.Writer.Header().Set("Content-Range", fmt.Sprintf("%s %d-%d/%d", generateResourceName(t), offset, offset+limit-1, n))
		c.JSON(http.StatusOK, data)
	}
}
func ParsedQuery(c *gin.Context) (sort string, offset, limit int, filter map[string]interface{}) {

	if q := c.Query("sort"); q != "" {
		var orderBy []string
		qq := strings.Replace(q, "'", "\"", -1)
		if err := json.Unmarshal([]byte(qq), &orderBy); err != nil {
			log.Errorf("fail to parse json: %v", qq)
		} else {
			sort = fmt.Sprintf("%s %s", orderBy[0], orderBy[1])
		}
	}
	if q := c.Query("range"); q != "" {
		rr := []int{0, 100}
		if err := json.Unmarshal([]byte(q), &rr); err != nil {
			log.Errorf("fail to parse json: %v", q)
		} else {
			offset = rr[0]
			limit = rr[1] + 1 - rr[0]
		}
	}
	if q := c.Query("filter"); q != "" {
		if err := json.Unmarshal([]byte(q), &filter); err != nil {
			log.Errorf("fail to parse json: %v", q)
		}
	}
	return
}
