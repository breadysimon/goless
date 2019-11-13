package rest

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type CreatePreProcessor interface {
	BeforeCreate()
}
type UpdatePreProcessor interface {
	BeforeUpdate()
}

func mdlRead(db *gorm.DB, o interface{}, id int) (err error) {
	if err = db.First(o, id).Error; err != nil {
		log.Error(err)
	}
	return
}

func mdlCreate(db *gorm.DB, o interface{}) (err error) {
	if p, ok := o.(CreatePreProcessor); ok {
		p.BeforeCreate()
	}
	if err = db.Create(o).Error; err != nil {
		log.Error(err)
	}
	return
}
func mdlUpdate(db *gorm.DB, o interface{}) (err error) {
	if p, ok := o.(UpdatePreProcessor); ok {
		p.BeforeUpdate()
	}
	if err = db.Save(o).Error; err != nil {
		log.Error(err)
	}
	return
}
func mdlDelete(db *gorm.DB, o interface{}) (err error) {
	if err = db.Delete(o).Error; err != nil {
		log.Error(err)
	}
	return
}
func mdlList(db *gorm.DB, o interface{}, rows interface{}, sort string, offset, limit int, filter map[string]interface{}) (
	count int) {
	search := getSearchableFields(o)
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
