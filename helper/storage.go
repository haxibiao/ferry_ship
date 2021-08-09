/*
 * @Author: Bin
 * @Date: 2021-08-05
 * @FilePath: /ferry_ship/helper/storage.go
 */
package helper

import (
	"encoding/json"
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

type StorageConfig struct {
	Id   int
	Name string
	Data string
}

// 辅助机器人插件获取配置数据
func GetConfigsDataByName(name string) (data map[string]interface{}, error error) {

	var c StorageConfig

	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	qb.Select("configs.id", "configs.name", "configs.data").From("configs").Where("name = ?")

	// 导出 SQL 语句
	sql := qb.String()

	// 执行 SQL 语句
	o := orm.NewOrm()
	error = o.Raw(sql, name).QueryRow(&c)

	if c.Id == 0 {
		return nil, errors.New("config is nil")
	}

	dataStr := c.Data
	var data_obj interface{}
	json.Unmarshal([]byte(dataStr), &data_obj)
	data = data_obj.(map[string]interface{})

	return data, error
}
