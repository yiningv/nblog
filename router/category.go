package router

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

// 分类管理：分类列表、添加分类、删除分类、修改分类

// 分类列表
func GetCategories(c *gin.Context) {
	cates, err := srv.GetCategories()
	if err != nil {
		fail(c, err.Error())
		return
	}
	okData(c, cates)
}

// 添加分类
func AddCategory(c *gin.Context) {
	cate := new(model.Category)
	if err := c.Bind(cate); err != nil {

	}
}

// 删除分类(批量删除)
func DeleteCategory(c *gin.Context) {
	arg := &struct {
		Ids []int64 `json:"ids" form:"ids"`
	}{}
	if err := c.ShouldBind(arg); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("DeleteCategory request params error: %v", err))
		return
	}
	if err := srv.DeleteCategories(arg.Ids); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("srv.DeleteCategories error: %v", err))
		return
	}
}

// 修改分类
func UpdateCategory(c *gin.Context) {
	arg := new(model.Category)
	if err := c.ShouldBind(arg); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("UpdateCategory request params error: %v", err))
		return
	}
	if arg.ID <= 0 {
		fail(c, "invalidate id")
		log.Error(fmt.Sprintf("UpdateCategory request params id: %d", arg.ID))
		return
	}
	if err := srv.UpdateCategory(arg); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("srv.UpdateCategory error: %v", err))
		return
	}
}
