package wallecontrollers

import (
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/models"
)

type TagController struct {
	controllers.BaseController
}

func (c *TagController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}

	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{})

	g := components.BaseGit{}
	g.SetBaseComponents(s)

	res, err := g.GetTagList()

	if err != nil {
		c.SetJson(1, nil, "获取tag错误—"+err.Error())
		return
	}
	c.SetJson(0, res, "")
	return
}
