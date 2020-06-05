package wallecontrollers

import (
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/models"
)

type DetectionController struct {
	controllers.BaseController
}

func (c *DetectionController) Get() {
	if c.Project == nil || c.Project.Id == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}

	s := components.BaseComponents{}
	s.SetProject(c.Project)
	s.SetTask(&models.Task{Id: -1})

	// 1. 权限与免密码登录检测
	if err := s.TestSsh(); err != nil {
		c.SetJson(1, nil, "ssh目标机器错误"+err.Error())
		return
	}

	// 2. clone repo
	if c.Project.RepoType == "git" {
		g := components.BaseGit{}
		g.SetBaseComponents(s)

		err := g.UpdateRepo("", "")
		if err != nil {
			// 清空后再试一次 要是不行在退出
			err := s.RemoveLocalProjectWorkspace()
			if err != nil {
				c.SetJson(1, nil, "清空目录错误"+err.Error())
				return
			}
			err = g.UpdateRepo("", "")
			if err != nil {
				c.SetJson(1, nil, "git拉取代码错误"+err.Error())
				return
			}
		}
	}

	// 3. 检测用户是否具有目标机release目录读写权限
	err := s.TestReleaseDir()
	if err != nil {
		c.SetJson(1, nil, "用户不具有目标机release目录读写权限"+err.Error())
		return
	}

	c.SetJson(0, nil, "")
	c.ServeJSON()
}
