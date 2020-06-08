package confcontrollers

import (
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/models"
	"github.com/linclin/gopub/src/util/validator"
)

type LockController struct {
	controllers.BaseController
}

func (c *LockController) Post() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}

	p := new(struct {
		ProjectId int `json:"project_id" validate:"gt=0"`
		Act       int `json:"act" validate:"oneof=0 1"`
	})
	if err := validator.Validate(c.Ctx.Input.RequestBody, p); err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}

	project, _ := models.GetProjectById(p.ProjectId)
	if project == nil {
		c.SetJson(1, nil, "No record")
		return
	}

	project.UserLock = 0
	if p.Act == 1 {
		project.UserLock = int(c.User.Id)
	}

	err := models.UpdateProjectById(project)
	if err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}
	c.SetJson(0, nil, "锁定成功")
	return
}
