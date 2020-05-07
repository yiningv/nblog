package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

func GetArticles(c *gin.Context) {
	arg := new(model.Page)
	if err := c.ShouldBind(arg); err != nil {
		fail(c, err.Error())
		log.Error(fmt.Sprintf("GetArticleList request error: %v", err))
		return
	}
	if arg.Pn <= 0 {
		arg.Pn = 1
	}
	pager, err := srv.GetArticles(arg.Pn, arg.Ps)
	if err != nil {
		fail(c, err.Error())
		return
	}
	okData(c, pager)
}

func AddArticle(c *gin.Context) {
	article := new(model.Article)
	if err := c.Bind(article); err != nil {
		fmt.Println(err)
	}
	fmt.Println(article)
	c.JSON(200, "hello")
}
