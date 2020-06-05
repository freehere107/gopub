package controllers

type LogoutController struct {
	BaseController
}

func (c *LogoutController) Post() {
	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success"}
	c.ServeJSON()
	return
}
