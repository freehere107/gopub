package confcontrollers

import (
	"encoding/json"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/models"
)

type SaveController struct {
	controllers.BaseController
}

func (c *SaveController) Post() {
	var project models.Project

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &project)
	if err != nil {
		c.SetJson(1, nil, "数据格式错误")
		return
	}

	if project.Id != 0 {
		err = models.UpdateProjectById(&project)
	} else {
		project.UserId = uint(c.User.Id)
		_, err = models.AddProject(&project)
	}

	if err != nil {
		c.SetJson(1, nil, "数据库更新错误")
		return
	}
	c.SetJson(0, project, "修改成功")
	return
}
