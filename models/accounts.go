package models

import (
	"errors"
	"ferry_ship/bot"
	"ferry_ship/helper"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Accounts struct {
	Id          int
	Account     int64     `orm:"unique; index; description(QQ 账号)"`
	Name        string    `orm:"null; size(128); description(QQ 昵称)"`
	Avatar      string    `orm:"null; size(128); description(QQ 账号头像)"`
	Password    string    `orm:"null; size(128);description(QQ 密码)"`
	Md5Password string    `orm:"null; size(128); description(QQ 账号 md5 密码)"`
	User        *Users    `orm:"rel(fk); description(添加 QQ 账号的用户)"`
	AutoLogin   int       `orm:"default(1); description(自动登陆: 启用[1] 停用[0])"`
	Status      int       `orm:"default(0); description(状态: 在线[1] 离线[0])"`
	Created     time.Time `orm:"auto_now_add; type(datetime)"`
	Updated     time.Time `orm:"auto_now; type(datetime)"`
}

func init() {
	// orm.RegisterModel(new(Accounts))
}

// AddAccounts insert a new Accounts into database and returns
// last inserted Id on success.
func AddAccounts(m *Accounts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		account, _ := GetAccountsById(int(id))
		if account != nil {
			// 触发登陆 QQ 账号
			AccountLoginQQ(account)
		}
	}
	return
}

// GetAccountsById retrieves Accounts by Id. Returns error if
// Id doesn't exist
func GetAccountsById(id int) (v *Accounts, err error) {
	o := orm.NewOrm()
	v = &Accounts{Id: id}
	if err = o.QueryTable(new(Accounts)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// 通过 account QQ 账号获取
func GetAccountsByAccount(account int64) (v *Accounts, err error) {
	o := orm.NewOrm()
	v = &Accounts{Account: account}
	if err = o.QueryTable(new(Accounts)).Filter("Account", account).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAccounts retrieves all Accounts matches certain condition. Returns empty list if
// no records exist
func GetAllAccounts(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Accounts))
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

	var l []Accounts
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

// UpdateAccounts updates Accounts by Id and returns error if
// the record to be updated doesn't exist
func UpdateAccountsById(m *Accounts) (err error) {
	o := orm.NewOrm()
	v := Accounts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAccounts deletes Accounts by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAccounts(id int) (err error) {
	o := orm.NewOrm()
	v := Accounts{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Accounts{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 获取全部 bot 账号
func AllAccounts(limit int, page int) (accounts []Accounts, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Accounts))
	_, err = qs.Filter("id__isnull", false).Limit(limit, page).All(&accounts)
	return
}

// Accounts 转 map 数据
func TurnAccountsToMap(account *Accounts) map[string]interface{} {
	return map[string]interface{}{
		"id":         account.Id,
		"name":       account.Name,
		"account":    account.Account,
		"avatar":     account.Avatar,
		"status":     account.Status,
		"auto_login": account.AutoLogin,
		"updated":    account.Updated,
		"created":    account.Created,
	}
}

// 登陆 QQ 账号
func AccountLoginQQ(m *Accounts) (account *Accounts) {

	if m.Password != "" {
		// 密码存在，使用密码登陆
		err := helper.InitBot(m.Account, m.Password)
		fmt.Printf("device lock -> %v\n", err)
	} else if m.Md5Password != "" {
		// 密码 md5 存在，使用 md5 登陆

	}

	account, _ = GetAccountsById(m.Id)
	return account
}

// 刷新全部机器人账号信息
func RefreshAccountBotInfo() error {

	// 重置全部机器人账号登陆状态
	o := orm.NewOrm()
	_, err := o.QueryTable(new(Accounts)).Update(orm.Params{
		"status": 0,
	})
	if err != nil {
		return err
	}

	for accountKey, botValue := range bot.Instances {

		if botValue == nil {
			continue
		}

		o := orm.NewOrm()
		online := 0
		if botValue.Online {
			online = 1
		}
		o.QueryTable(new(Accounts)).Filter("account", accountKey).Update(orm.Params{
			"name":   botValue.Nickname,
			"avatar": "https://q2.qlogo.cn/headimg_dl?spec=100&dst_uin=" + strconv.FormatInt(botValue.Uin, 10),
			"status": online,
		})

	}

	return nil
}

// 重设机器人账号密码
func AccountReupdatePassword(m *Accounts, password string) (account *Accounts, ok bool) {

	o := orm.NewOrm()
	o.QueryTable(new(Accounts)).Filter("id", m.Id).Update(orm.Params{
		"password": password,
	})
	account, err := GetAccountsById(m.Id)

	if account != nil && err == nil {

		// 判断密码是否修改成功
		if account.Password == password {
			return account, true
		}
	}

	return account, false
}
