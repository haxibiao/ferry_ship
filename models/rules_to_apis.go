/*
 * @Author: Bin
 * @Date: 2021-07-07
 * @FilePath: /ferry_ship/models/rules_to_apis.go
 */
package models

import (
	"time"
)

type RulesToApis struct {
	Id        int64
	Rules     *Rules    `orm:"rel(fk)"`
	Apis      *Apis     `orm:"rel(fk)"`
	CreatedAt time.Time `orm:"type(datetime)"`
	UpdatedAt time.Time `orm:"type(datetime)"`
}

func init() {
	// orm.RegisterModel(new(RoleUser))
}

func (m *RulesToApis) TableName() string {
	return "rules_to_apis"
}
