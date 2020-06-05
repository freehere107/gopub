package controllers

import (
	"github.com/linclin/gopub/src/util/validator"

	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Post() {
	p := new(struct {
		UserPassword string `json:"user_password" validate:"required"`
		UserName     string `json:"user_name" validate:"required"`
	})
	if err := validator.Validate(c.Ctx.Input.RequestBody, p); err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}

	var user models.User
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM `user` WHERE username= ?", p.UserName).QueryRow(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p.UserPassword))
	if err != nil {
		c.SetJson(1, nil, "用户名或密码错误")
		return
	}

	if user.AuthKey == "" {
		user.AuthKey = common.Md5String(user.Username + common.GetString(time.Now().Unix()))
		_ = models.UpdateUserById(&user)
	}
	user.PasswordHash = ""

	c.SetJson(0, user, "")
	return
}
