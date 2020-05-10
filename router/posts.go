package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"github.com/yiningv/nblog/service"
)

func GetPosts(c *gin.Context) {
	arg := new(model.Page)
	if err := c.ShouldBind(arg); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("GetPosts request error: %v", err))
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	pager, err := service.Posts.GetPostsPager(arg.Pn, arg.Ps)
	if err != nil {
		log.Error(fmt.Sprintf("srv.GetPostsPager error: %v", err))
		fail(c, err.Error())
		return
	}
	okData(c, pager)
}

func AddPosts(c *gin.Context) {
	posts := new(model.Posts)
	if err := c.Bind(posts); err != nil {
		fmt.Println(err)
	}
	fmt.Println(posts)
	c.JSON(200, "hello")
}
