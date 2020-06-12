package confcontrollers

import (
	"github.com/astaxie/beego/logs"
	"github.com/linclin/gopub/src/controllers"

	"encoding/json"
	"github.com/linclin/gopub/src/models"
)

type ConfController struct {
	controllers.BaseController
}

func (c *ConfController) Get() {
	project, _ := models.GetProjectById(c.Project.Id)
	c.SetJson(0, project, "")
	return

}
func (c *ConfController) Post() {
	var project models.Project
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &project)
	err = models.UpdateProjectById(&project)
	logs.Info(err)
	return
}
