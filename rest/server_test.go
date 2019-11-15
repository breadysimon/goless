package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/breadysimon/goless/reflection"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

type T struct {
	ID   int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name string `gorm:"type:varchar(50);unique_index" json:"name" rest:"search"`
	Memo string `json:"memo" rest:"search"`
}
type TX struct {
	ID   int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name string `gorm:"type:varchar(50);unique_index" json:"name" rest:"search"`
	XXX  string `json:"memo" rest:"search"`
}

var r *gin.Engine

type kv map[string]string

func TestReflect(t *testing.T) {
	sss := Setup()
	f := reflection.GetSearchableFieldNames(&T{Name: "123", Memo: "234"})
	fmt.Println(f)

	sss.db.Create(&T{Name: "123", Memo: "124"})
	sss.db.Create(&T{Name: "1231"})

	xx := reflection.MakeSlice(&T{})
	sss.db.Find(xx)
	fmt.Println(xx)

}

func TestApi(t *testing.T) {
	sss := Setup()
	{
		// create item
		url := "/api/v1/ts"
		body := strings.NewReader(`{"name":"123"}`)
		req, _ := http.NewRequest("POST", url, body)
		req.Header.Add("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get right status code.")

		respBody, _ := ioutil.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), "123", "Should find string in response body.")

		var p T
		err := json.Unmarshal(respBody, &p)
		assert.Nil(t, err)
		assert.Equal(t, 1, p.ID)
	}

	{
		// get item
		url := "/api/v1/ts/1"
		req, _ := http.NewRequest("GET", url, nil)
		// req.Header.Add("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get right status code.")

		respBody, _ := ioutil.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), "123", "Should find string in response body.")

		var p T
		err := json.Unmarshal(respBody, &p)
		assert.Nil(t, err)
		assert.Equal(t, 1, p.ID)
	}
	{
		// create item
		url := "/api/v1/ts/1"
		body := strings.NewReader(`{"name":"abc"}`)
		req, _ := http.NewRequest("PUT", url, body)
		req.Header.Add("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get right status code.")

		respBody, _ := ioutil.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), "abc", "Should find string in response body.")

		var p T
		err := json.Unmarshal(respBody, &p)
		assert.Nil(t, err)
		assert.Equal(t, 1, p.ID)
	}
	{
		// list item
		sss.db.Create(&T{Name: "abc", Memo: "xyzab"})
		sss.db.Create(&T{Name: "ab1", Memo: "xyzab"})
		sss.db.Create(&T{Name: "b2", Memo: "xyzab"})
		sss.db.Create(&T{Name: "ab3", Memo: "xyza"})
		url := "/api/v1/ts"
		query := kv{
			"filter": `{"q":"ab"}`,
			"range":  "[0,2]",
			"sort":   `["id","DESC"]`}

		req, _ := http.NewRequest("GET", url, nil)
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		req.Header.Add("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get right status code.")

		respBody, _ := ioutil.ReadAll(resp.Body)
		assert.Contains(t, string(respBody), "ab3", "Should find string in response body.")

		assert.Equal(t, "ts 0-2/4", resp.Header.Get("Content-Range"))

	}
	{
		// delete item
		url := "/api/v1/ts/1"
		req, _ := http.NewRequest("DELETE", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should get right status code.")
	}
}
func Setup() *RestApi {

	r = gin.New()
	{
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
		gin.SetMode(gin.TestMode)
	}
	// CORS control
	{
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AllowCredentials = true
		config.AddExposeHeaders("Content-Range")
		r.Use(cors.New(config))
	}
	return NewRestApi(&T{}, &TX{}).
		Connect("sqlite3", ":memory:").
		CreateEndpoints(r, "/api/v1", false)
}
