package taskcontrollers

import (
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/util/validator"
)

type ListController struct {
	controllers.BaseController
}

func (c *ListController) Post() {
	var (
		projects     []orm.Params
		count        []orm.Params
		start, total int
	)
	p := new(struct {
		My         int    `json:"my" validate:"gte=-1"`
		Page       int    `json:"page" validate:"gt=0"`
		Length     int    `json:"length" validate:"gte=1"`
		SelectInfo string `json:"select_info"`
	})
	if err := validator.Validate(c.Ctx.Input.RequestBody, p); err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}

	start = (p.Page - 1) * p.Length

	where := ""
	if p.SelectInfo != "" {
		where = "  and( project.`name` LIKE '%" + p.SelectInfo + "%' or `user`.realname LIKE '%" + p.SelectInfo + "%'  or task.title LIKE '%" + p.SelectInfo + "%'  )"
	}

	if p.My != 0 {
		where = where + "  and task.user_id=" + common.GetString(p.My)
	}

	o := orm.NewOrm()
	_, _ = o.Raw("SELECT task.id,project.name,project.name,project.level,`user`.realname,task.title,task.action,task.link_id,task.is_run,task.enable_rollback,task.updated_at,task.branch,task.commit_id,task.pms_uwork_id,task.pms_batch_id,task.`status` FROM `task` LEFT JOIN project on task.project_id=project.id   LEFT JOIN `user` on task.user_id=user.id where 1=1 "+where+" order by task.id DESC  LIMIT ? ,?", start, p.Length).Values(&projects)
	_, _ = o.Raw("SELECT count(task.id) FROM `task` LEFT JOIN project on task.project_id=project.id  LEFT JOIN `user` on task.user_id=user.id where 1=1 " + where).Values(&count)

	if len(count) > 0 {
		total = common.GetInt(count[0]["count(task.id)"])
	}

	for _, project := range projects {
		project["status"] = getTaskStatus(common.GetInt(project["status"]))

		if common.GetInt(project["is_run"]) != 0 && common.GetString(project["status"]) != "上线完成" {
			project["status"] = "上线中"
		}

		if common.GetInt(project["level"]) == 3 {
			project["name"] = common.GetString(project["name"]) + "-线上环境"
		}

		if common.GetInt(project["level"]) == 2 {
			project["name"] = common.GetString(project["name"]) + "-预发布环境"
		}
	}
	c.SetJson(0, map[string]interface{}{"total": total, "currentPage": p.Page, "table_data": projects}, "")
	return

}
func getTaskStatus(status int) string {
	ts := map[int]string{
		0: "新建提交",
		1: "新建提交",
		2: "审核拒绝",
		3: "上线完成",
		4: "上线失败",
	}
	return ts[status]
}
