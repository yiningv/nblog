package dao

import (
	"github.com/smartystreets/goconvey/convey"
	"github.com/yiningv/nblog/model"
	"testing"
)

func TestDaoAddCategory(t *testing.T) {
	convey.Convey("AddCategory", t, func(ctx convey.C) {
		var (
			arg = &model.Category{
				Name:         "Golang",
				Slug:         "/golang",
				ArticleCount: 0,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCategory(arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
