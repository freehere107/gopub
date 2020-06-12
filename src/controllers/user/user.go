package usercontrollers

import (
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/controllers"
)

type UserController struct {
	controllers.BaseController
}

func (c *UserController) Get() {
	o := orm.NewOrm()
	var (
		users []orm.Params
		res   orm.Params
	)
	userId, _ := c.GetInt("id")
	if userId == 0 {
		_, _ = o.Raw("SELECT * FROM `user` ").Values(&users)
		c.SetJson(0, users, "")
		return
	}
	if i, _ := o.Raw("SELECT * FROM `user` where id = ? ", userId).Values(&users); i > 0 {
		res = users[0]
	}
	c.SetJson(0, res, "")
	return
}
