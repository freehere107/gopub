package taskcontrollers

import (
	"encoding/json"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/models"
)

type SaveController struct {
	controllers.BaseController
}

func (c *SaveController) Post() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}

	var task models.Task
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &task)
	if err != nil {
		c.SetJson(1, nil, "数据库更新错误"+err.Error())
	}

	if task.Id != 0 {
		if err = models.UpdateTaskById(&task); err != nil {
			c.SetJson(1, nil, "数据库更新错误")
			return
		}
		c.SetJson(0, task, "修改成功")
		return
	}

	task.UserId = uint(c.User.Id)
	task.EnableRollback = 1

	if task.Hosts == "" {
		if ss, _ := models.GetProjectById(task.ProjectId); ss != nil {
			task.Hosts = ss.Hosts
			task.HostGroup = ss.HostGroup
			if ss.IsGroup == 1 {
				s := components.BaseComponents{}
				s.SetProject(ss)
				mapHost := s.GetGroupHost()
				for k, v := range mapHost {
					task1 := task
					task1.Hosts = v
					task1.Title = task1.Title + "第" + common.GetString(k) + "批"
					_, _ = models.AddTask(&task1)
				}
				c.SetJson(0, task, "修改成功")
				return
			}
		}
	}

	if _, err = models.AddTask(&task); err != nil {
		c.SetJson(1, nil, "数据库更新错误")
	}
	c.SetJson(0, task, "修改成功")

	return
}
