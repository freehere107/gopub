package confcontrollers

import (
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/common"
	"github.com/linclin/gopub/src/util/validator"
)

type MyListController struct {
	controllers.BaseController
}

func (c *MyListController) Post() {
	var (
		where           string
		projects, count []orm.Params
		total           int
	)

	p := new(struct {
		Page       int    `json:"page" validate:"gt=0"`
		Length     int    `json:"length" validate:"gte=1"`
		SelectInfo string `json:"select_info"`
	})
	if err := validator.Validate(c.Ctx.Input.RequestBody, p); err != nil {
		c.SetJson(1, nil, err.Error())
		return
	}

	start := (p.Page - 1) * p.Length
	if p.SelectInfo != "" {
		where = "  and(`name` LIKE '%" + p.SelectInfo + "%' )"
	}

	if c.User.Role == 10 {
		where = where + "and  `level`= 2  "
	} else if c.User.Role == 20 {
		where = where + "and  id in (SELECT project_id FROM `group` WHERE `group`.user_id=" + common.GetString(c.User.Id) + " )  "
	}

	o := orm.NewOrm()
	_, _ = o.Raw("SELECT *, (SELECT realname FROM `user` WHERE `user`.id=project.user_id LIMIT 1) as realname,(SELECT realname FROM `user` WHERE `user`.id=project.user_lock LIMIT 1) as lockuser FROM `project`  WHERE 1=1 "+where+" ORDER BY id LIMIT ?,?", start, p.Length).Values(&projects)
	_, _ = o.Raw("SELECT count(id) FROM `project` WHERE 1=1 " + where).Values(&count)

	if len(count) > 0 {
		total = common.GetInt(count[0]["count(id)"])
	}

	c.SetJson(0, map[string]interface{}{"total": total, "currentPage": p.Page, "table_data": projects}, "")
	return
}
