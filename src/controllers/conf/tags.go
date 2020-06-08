package confcontrollers

import (
	"github.com/astaxie/beego/orm"
	"github.com/linclin/gopub/src/controllers"
	"github.com/linclin/gopub/src/library/common"
	"strings"
)

type TagsController struct {
	controllers.BaseController
}

func (c *TagsController) Get() {
	if c.User == nil || c.User.Id == 0 {
		c.SetJson(2, nil, "not login")
		return
	}

	var (
		rows []orm.Params
		a    []string
	)

	o := orm.NewOrm()
	_, _ = o.Raw("SELECT tag FROM `project`").Values(&rows)

	for _, row := range rows {

		tmp := strings.Split(common.GetString(row["tag"]), " ")
		for _, tag := range tmp {
			if tag != "" {
				a = append(a, tag)
			}
		}
	}

	a = common.ArrayUnique(a)
	c.SetJson(0, a, "")
	return

}
