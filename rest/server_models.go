package rest

import (
	"fmt"
	"strings"

	"github.com/breadysimon/goless/reflection"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func mdlRead(db *gorm.DB, o interface{}, id int) (err error) {
	if err = db.First(o, id).Error; err != nil {
		log.Error(err)
	}
	return
}

func mdlCreate(db *gorm.DB, o interface{}) (err error) {
	return db.Create(o).Error
}
func mdlUpdate(db *gorm.DB, o interface{}) (err error) {
	return db.Save(o).Error
}
func mdlDelete(db *gorm.DB, o interface{}) (err error) {
	if err = db.Delete(o).Error; err != nil {
		log.Error(err)
	}
	return
}
func mdlList(db *gorm.DB, o interface{}, rows interface{}, sort string, offset, limit int, filter map[string]interface{}) (
	count int) {
	search := reflection.GetSearchableFieldNames(o)
	page, all := query(db, sort, offset, limit, filter, search)
	all.Model(o).Count(&count)
	page.Find(rows)
	return
}
func query(db *gorm.DB, sort string, offset, limit int, filter map[string]interface{}, searchFields []string) (page, all *gorm.DB) {
	all = db

	if fuzzyQuery, ok := filter["q"]; ok {
		var phases []string
		for _, i := range searchFields {
			phases = append(phases,
				fmt.Sprintf(" %s like \"%%%s%%\" ", i, fuzzyQuery.(string)))
		}
		all = all.Where(strings.Join(phases, " OR "))
		delete(filter, "q")
	}
	all = all.Where(filter)
	page = all
	if sort != "" {
		page = page.Order(sort)
	}
	if offset > 0 {
		page = page.Offset(offset)
	}
	if limit > 0 {
		page = page.Limit(limit)
	}
	return
}
