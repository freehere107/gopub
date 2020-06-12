package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/util/validator"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/models"
	"strings"
)

// 基类
type BaseController struct {
	beego.Controller
	Project *models.Project
	Task    *models.Task
	User    *models.User
}

type DefaultParams struct {
	ProjectId string `json:"projectId,omitempty"`
	TaskId    string `json:"taskId,omitempty"`
}

// Prepare implemented Prepare method for baseRouter.
func (c *BaseController) Prepare() {

	// 获取panic
	defer func() {
		if err := recover(); err != nil {
			var buf = make([]byte, 1024)
			trace := runtime.Stack(buf, false)
			logs.Error("控制器错误:", err, string(buf[0:trace]))
		}
	}()

	var (
		defaultParams            DefaultParams
		taskId, projectId, token string
	)
	_ = validator.Validate(c.Ctx.Input.RequestBody, &defaultParams)

	if c.Ctx.Input.Param(":taskId") != "" {
		taskId = c.Ctx.Input.Param(":taskId")
	} else if c.GetString("taskId") != "" {
		taskId = c.GetString("taskId")
	} else {
		taskId = defaultParams.TaskId
	}

	if taskId != "" {
		c.Task, _ = models.GetTaskById(common.GetInt(taskId))
	}

	if c.Ctx.Input.Param(":projectId") != "" {
		projectId = c.Ctx.Input.Param(":projectId")
	} else if c.GetString("projectId") != "" {
		projectId = c.GetString("projectId")
	} else {
		projectId = defaultParams.ProjectId
	}
	if projectId != "" {
		c.Project, _ = models.GetProjectById(common.GetInt(projectId))
	}

	if ah := c.Ctx.Input.Header("Authorization"); ah != "" {
		if len(ah) > 5 && strings.ToUpper(ah[0:5]) == "TOKEN" {
			token = ah[6:]
			if token != "" {
				var users []models.User
				o := orm.NewOrm()
				if s, _ := o.Raw("SELECT * FROM `user` WHERE auth_key= ?", token).QueryRows(&users); s > 0 {
					c.User = &(users[0])
				}
			}
		}
	}
}
func (c *BaseController) SetJson(code int, data interface{}, Msg string) {
	if code == 0 {
		if Msg == "" {
			Msg = "success"
		}
		c.Data["json"] = map[string]interface{}{"code": code, "msg": Msg, "data": data}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"code": code, "msg": Msg, "data": data}
	c.ServeJSON()

}
