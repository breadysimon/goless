package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func ReadHandler(db *gorm.DB, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		m := makeItem(t)
		id, _ := strconv.Atoi(c.Param("id"))
		if err := mdlRead(db, m, id); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, m)
		}
	}
}

func CreateHandler(db *gorm.DB, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		m := makeItem(t)
		if err := c.BindJSON(m); err != nil {
			log.Error(err)
		} else {
			if err = mdlCreate(db, m); err == nil {
				c.JSON(http.StatusOK, m)
			}
		}
		//TODO: error handling
	}
}
func UpdateHandler(db *gorm.DB, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		m := makeItem(t)
		if c.BindJSON(&m) == nil {
			SetId(m, id)
			mdlUpdate(db, m)
		}
		c.JSON(http.StatusOK, m)
		//TODO: error
	}
}

func DeleteHandler(db *gorm.DB, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		m := makeItem(t)
		SetId(m, id)
		mdlDelete(db, m)
		//TODO:error
	}
}

func ListHandler(db *gorm.DB, t interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		m := makeItem(t)
		data := makeSlice(t)
		sort, offset, limit, filter := ParsedQuery(c)
		n := mdlList(db, m, data, sort, offset, limit, filter)
		c.Writer.Header().Set("Content-Range", fmt.Sprintf("%s %d-%d/%d", guessResourceName(t), offset, offset+limit-1, n))
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
