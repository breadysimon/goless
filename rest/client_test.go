package rest

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/breadysimon/goless/reflection"
	"github.com/stretchr/testify/assert"
)

// var log *logging.Logger = logging.GetLogger()
type Template struct {
	ID   int    `gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name string `gorm:"column:name;type:varchar(50);unique_index" json:"name"`
	XXX  string `json:"xxx"`
}

func TestDataOps(t *testing.T) {
	StartServer()
	defer TearDownServer()
	time.Sleep(1 * time.Second)
	api := NewClient("http://127.0.0.1:8881", "/api/v1/", "admin", "admin")
	var tmps []Template
	tmp := Template{ID: 23}
	assert.Nil(t, api.
		Create(&Template{Name: "ttt1"}).
		Find(&tmps).peekResult().
		Create(&Template{Name: "testa"}).
		Find(&Template{ID: 1}).peekResult().
		Find(&tmps).peekResult().
		Update(&Template{ID: 1, Name: "testx"}).peekResult().
		Update(&Template{ID: 1, Name: "testy"}).peekResult().
		Delete(&Template{ID: tmp.ID}).
		Find(&tmps).peekResult().
		Error())
	fmt.Println(tmps)
	fmt.Println(tmp)
}

func TestFindSingle(t *testing.T) {
	StartServer()
	defer TearDownServer()
	time.Sleep(1 * time.Second)
	api := NewClient("http://127.0.0.1:8881", "/api/v1/", "admin", "admin")

	tmp := Template{Name: "ttt1"}
	assert.Nil(t, api.
		Create(&Template{Name: "ttt1", XXX: "ff"}).
		Find(&tmp).peekResult().
		Error())
	assert.Equal(t, "ff", tmp.XXX)
}
func TestTypeOf(t *testing.T) {
	var a interface{}

	c := Template{ID: 123}
	a = &c
	assert.Equal(t, "templates", generateResourceName(a))
	id := reflection.GetId(a).Int()
	fmt.Println(id)
	l := []Template{}
	a = &l
	fmt.Println(reflect.TypeOf(a))
}

var svr *RestApi

func StartServer() {
	svr = NewRestApi(&T{}, &Template{}).
		Connect("sqlite3", ":memory:").
		Server("", 8881, "/api/v1/", true).
		Start()

}
func TearDownServer() {
	svr.Shutdown()
}
