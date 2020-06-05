package controllers

import (
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/models"
	"github.com/linclin/gopub/src/util/validator"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswdController struct {
	BaseController
}

func (c *ChangePasswdController) Post() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}
	p := new(struct {
		NewPassword       string `json:"newpassword" validate:"required"`
		RepeatNewPassword string `json:"repeat_newpassword" validate:"required,eqfield=NewPassword"`
		Uid               string `json:"uid" validate:"required"`
	})

	if err := validator.Validate(c.Ctx.Input.RequestBody, p); err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}

	if common.GetString(c.User.Id) != p.Uid && c.User.Role != 1 {
		c.SetJson(1, nil, "403")
		return
	}

	var user models.User
	o := orm.NewOrm()
	_ = o.Raw("SELECT * FROM `user` WHERE id= ?", p.Uid).QueryRow(&user)

	// 验证旧密码
	password := []byte(p.NewPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user.PasswordHash = string(hashedPassword)
	_ = models.UpdateUserById(&user)

	c.Data["json"] = map[string]interface{}{"code": 0, "msg": "success"}
	c.ServeJSON()
	return

}
