package rest

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// var log *logging.Logger = logging.GetLogger()
type Template struct {
	ID   int    `gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Name string `gorm:"column:name;type:varchar(50);unique_index" json:"name"`
}

func TestDataOps(t *testing.T) {
	SetupServer()
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

func TestTypeOf(t *testing.T) {
	var a interface{}

	c := Template{ID: 123}
	a = &c
	assert.Equal(t, "credentials", guessResourceName(a))
	id := GetId(a).Int()
	fmt.Println(id)
	l := []Template{}
	a = &l
	fmt.Println(reflect.TypeOf(a))
}

var svr *http.Server

func SetupServer() {
	svr = NewRestApi(&T{}, &Template{}).
		Connect("sqlite3", ":memory:").
		Serve("", 8881, "/api/v1/", true)

}
func TearDownServer() {
	svr.Shutdown(context.TODO())
}
