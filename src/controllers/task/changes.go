package taskcontrollers

import (
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/models"
)

type ChangesController struct {
	controllers.BaseController
}

func (c *ChangesController) Get() {
	taskId, _ := c.GetInt("taskId", 0)

	if taskId == 0 {
		c.SetJson(1, nil, "Parameter error")
		return
	}
	o := orm.NewOrm()

	var task models.Task
	_ = o.Raw("SELECT * FROM task where task.id = ?", taskId).QueryRow(&task)

	project, _ := models.GetProjectById(task.ProjectId)

	if project.RepoType == "git" {
		var lastTask models.Task
		_ = o.Raw("SELECT * FROM task where project_id = ? AND status=3 order by task.id DESC LIMIT 1", task.ProjectId).QueryRow(&lastTask)

		s := components.BaseComponents{}
		s.SetProject(project)
		s.SetTask(&task)

		g := components.BaseGit{}
		g.SetBaseComponents(s)
		files, _ := g.DiffBetweenCommits(task.Branch, task.CommitId, lastTask.CommitId)

		var fileInfos []map[string]string
		if len(files) < 10 && len(files) > 0 {
			for _, filepath := range files {
				fileInfo, _ := g.GetLastModifyInfo(task.Branch, filepath)
				fileInfo["path"] = filepath
				fileInfos = append(fileInfos, fileInfo)
			}
		}
		c.SetJson(0, fileInfos, "")
		return
	} else {
		c.SetJson(1, nil, "Project is not git")
	}

}
