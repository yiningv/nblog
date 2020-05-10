package router

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"github.com/yiningv/nblog/service"
)

// 标签管理：标签列表、添加标签、删除标签、修改标签

// 标签列表，分页获取
func GetTags(c *gin.Context) {
	arg := new(model.Page)
	if err := c.ShouldBind(arg); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("GetTags request error: %v", err))
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	pager, err := service.Tag.GetTags(arg.Pn, arg.Ps)
	if err != nil {
		fail(c, err.Error())
		return
	}
	okData(c, pager)
}

// 添加标签
func AddTag(c *gin.Context) {
	cate := new(model.Tag)
	if err := c.Bind(cate); err != nil {

	}
}

// 删除标签(批量删除)
func DeleteTag(c *gin.Context) {
	arg := &struct {
		Ids []int `json:"ids" form:"ids"`
	}{}
	if err := c.ShouldBind(arg); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("DeleteTag request params error: %v", err))
		return
	}
	if err := service.Tag.BatchDeleteTag(arg.Ids); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("srv.BatchDeleteTag error: %v", err))
		return
	}
	ok(c)
}

// 修改标签
func UpdateTag(c *gin.Context) {

}
