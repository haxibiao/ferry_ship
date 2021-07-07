package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Apis struct {
	Id       int
	Name     string    `orm:"size(128); description(名称)"`
	Request  int       `orm:"default(1); description(状态: Get请求[1] Post请求[0])"`
	Url      string    `orm:"size(128); description(API 地址)"`
	Args     string    `orm:"size(128); description(请求参数名称)"`
	Analysis string    `orm:"size(128); description(数据解析规则)"`
	Callback string    `orm:"null; size(128); description(请求返回 Json 数据缓存)"`
	User     *Users    `orm:"rel(fk); description(添加接口的用户)"`
	Rules    []*Rules  `orm:"reverse(many); null; on_delete(set_null)"`
	Created  time.Time `orm:"auto_now_add; type(datetime)"`
	Updated  time.Time `orm:"auto_now; type(datetime)"`
}

func init() {
	// orm.RegisterModel(new(Apis))
}

// AddApis insert a new Apis into database and returns
// last inserted Id on success.
func AddApis(m *Apis) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetApisById retrieves Apis by Id. Returns error if
// Id doesn't exist
func GetApisById(id int) (v *Apis, err error) {
	o := orm.NewOrm()
	v = &Apis{Id: id}
	if err = o.QueryTable(new(Apis)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApis retrieves all Apis matches certain condition. Returns empty list if
// no records exist
func GetAllApis(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Apis))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Apis
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateApis updates Apis by Id and returns error if
// the record to be updated doesn't exist
func UpdateApisById(m *Apis) (err error) {
	o := orm.NewOrm()
	v := Apis{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApis deletes Apis by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApis(id int) (err error) {
	o := orm.NewOrm()
	v := Apis{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Apis{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
