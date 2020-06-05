package confcontrollers

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
		where = "  and(`name` LIKE '%" + p.SelectInfo + "%' )"
	}

	o := orm.NewOrm()
	_, _ = o.Raw("SELECT *, (SELECT realname FROM `user` WHERE `user`.id=project.user_id LIMIT 1) as realname FROM `project`  WHERE 1=1 "+where+" ORDER BY id LIMIT ?,?", start, p.Length).Values(&projects)
	_, _ = o.Raw("SELECT count(id) FROM `project` WHERE 1=1 " + where).Values(&count)

	if len(count) > 0 {
		total = common.GetInt(count[0]["count(id)"])
	}
	c.SetJson(0, map[string]interface{}{"total": total, "currentPage": p.Page, "table_data": projects}, "")

	return

}
