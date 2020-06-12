package taskcontrollers

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/library/components"
	"github.com/linclin/gopub/src/models"
	"time"
)

type TaskChartController struct {
	controllers.BaseController
}

var bm, _ = cache.NewCache("memory", `{"interval":3600}`)

func (c *TaskChartController) Get() {
	taskType := c.GetString("taskType")
	var count, totalmem, totalproject, totalpub, totalpubsuccess []orm.Params
	o := orm.NewOrm()

	switch taskType {
	case "day":
		_, _ = o.Raw("SELECT project.`level`,count(task.id) as task_count  FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE TO_DAYS(now()) - TO_DAYS(task.updated_at) = 0 GROUP BY project. LEVEL").Values(&count)
		for _, c := range count {
			c["name"] = GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
	case "week":
		_, _ = o.Raw("SELECT project.`level`,count(task.id) as task_count  FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE YEARWEEK(date_format(task.updated_at,'%Y-%m-%d')) = YEARWEEK(now()) GROUP BY project. LEVEL").Values(&count)
		for _, c := range count {
			c["name"] = GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")

	case "month":
		_, _ = o.Raw("SELECT project.`level`,count(task.id) as task_count  FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE date_format(task.updated_at,'%Y-%m')=date_format(now(),'%Y-%m') GROUP BY project. LEVEL").Values(&count)
		for _, c := range count {
			c["name"] = GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
	case "dayBypro":
		_, _ = o.Raw("SELECT project.`name`,count(task.id) as task_count,project.`level` FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE TO_DAYS(now()) - TO_DAYS(task.updated_at) = 0 and task.status=3 GROUP BY project.id").Values(&count)
		for _, c := range count {
			c["name"] = common.GetString(c["name"]) + "-" + GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")

	case "weekBypro":
		_, _ = o.Raw("SELECT project.`name`,count(task.id) as task_count,project.`level` FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE YEARWEEK(date_format(task.updated_at,'%Y-%m-%d')) = YEARWEEK(now()) and task.status=3 GROUP BY project.id").Values(&count)
		for _, c := range count {
			c["name"] = common.GetString(c["name"]) + "-" + GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
	case "monthBypro":
		_, _ = o.Raw("SELECT project.`name`,count(task.id) as task_count,project.`level` FROM `task` LEFT JOIN project ON task.project_id = project.id WHERE date_format(task.updated_at,'%Y-%m')=date_format(now(),'%Y-%m') and task.status=3 GROUP BY project.id").Values(&count)
		for _, c := range count {
			c["name"] = common.GetString(c["name"]) + "-" + GetProjectLevel(common.GetInt(c["level"]))
		}
		c.SetJson(0, count, "")
	case "total":
		totalJson := map[string]interface{}{}
		num, err := o.Raw("SELECT count(id) as `totalmen` FROM `user`").Values(&totalmem)
		if num > 0 && err == nil {
			totalJson["totalmen"] = common.GetInt(totalmem[0]["totalmen"])
		}
		num, err = o.Raw("SELECT count(DISTINCT name) as `totalproject` from `project`").Values(&totalproject)
		if num > 0 && err == nil {
			totalJson["totalproject"] = common.GetInt(totalproject[0]["totalproject"])
		}
		num, err = o.Raw("SELECT count(id) as `totalpub` from `task`").Values(&totalpub)
		if num > 0 && err == nil {
			totalJson["totalpub"] = common.GetInt(totalpub[0]["totalpub"])
		}
		num, err = o.Raw("SELECT count(id) as `totalpubsuccess` from `task`where status = 3").Values(&totalpubsuccess)
		if num > 0 && err == nil {
			totalJson["totalpubsuccess"] = common.GetInt(totalpubsuccess[0]["totalpubsuccess"])
		}
		if bm.IsExist("hostsum") == false {
			totalJson["hostsum"] = GetHostNum()
		} else {
			totalJson["hostsum"] = bm.Get("hostsum")
		}
		c.SetJson(0, totalJson, "")
	default:
		c.SetJson(1, nil, "未传参数")
	}

}
func GetProjectLevel(level int) string {
	switch level {
	case 1:
		return "测试环境"
	case 2:
		return "预发布环境"
	case 3:
		return "生产环境"
	}
	return "删除项目"
}

func GetHostNum() int {
	o := orm.NewOrm()
	var projects []models.Project
	i, err := o.Raw("SELECT * FROM `project`").QueryRows(&projects)
	var finalRes []string
	if i > 0 && err == nil {
		for _, project := range projects {
			s := components.BaseComponents{}
			s.SetProject(&project)
			ips := s.GetHostIps()
			for _, ip := range ips {
				if !common.InList(string(ip), finalRes) {
					finalRes = append(finalRes, string(ip))
				}
			}
		}
	}
	_ = bm.Put("hostsum", len(finalRes), 1*time.Hour)
	return len(finalRes)
}
