package controllers

import (
	"ferry_ship/models"
	// "strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// SystemController operations for System
type SystemController struct {
	beego.Controller
}

/**
 * @description: 保存消息模版数据
 * @param {*}
 * @return {*}
 */
func (c *SystemController) ApiSaveMessageTemplate() {
	userAssistant(&c.Controller) // 登陆认证

	configName := "plugin_config_message_template"
	config, err := models.GetConfigsDataByName(configName)

	msgSuccess := c.GetString("success", "")
	msgEmpty := c.GetString("empty", "")
	msgFail := c.GetString("fail", "")

	if err == nil && config != nil {

		if msgSuccess != "" {
			config["success"] = msgSuccess
		}
		if msgEmpty != "" {
			config["empty"] = msgEmpty
		}
		if msgFail != "" {
			config["fail"] = msgFail
		}

		if _, err = models.UpdateConfigByData(configName, config); err != nil {
			callBackResult(&c.Controller, 200, err.Error(), nil)
			c.Finish()
		}
	} else {
		models.AddConfigByData(configName, map[string]interface{}{
			"success": msgSuccess,
			"empty":   msgEmpty,
			"fail":    msgFail,
		})
	}

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}

/**
 * @description: 获取消息模版数据
 * @param {*}
 * @return {*}
 */
func (c *SystemController) ApiGetMessageTemplate() {
	userAssistant(&c.Controller) // 登陆认证

	data, err := models.GetConfigsDataByName("plugin_config_message_template")
	if err != nil || data == nil {
		data = map[string]interface{}{
			"success": "", // 成功的消息模版
			"empty":   "", // 空的消息模版
			"fail":    "", // 失败的消息模版
		}
	}

	// success, _ := strconv.Unquote("\"" + data["success"].(string) + "\"")
	// empty, _ := strconv.Unquote("\"" + data["empty"].(string) + "\"")
	// fail, _ := strconv.Unquote("\"" + data["fail"].(string) + "\"")

	callBackResult(&c.Controller, 200, "", map[string]interface{}{
		"success": data["success"].(string), // 成功的消息模版
		"empty":   data["empty"].(string),   // 空的消息模版
		"fail":    data["fail"].(string),    // 失败的消息模版
	})
	c.Finish()
}
