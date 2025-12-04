package modules

import (
	"github.com/youbuwei/doeot-go/pkg/biz"
	"gorm.io/gorm"
)

// Factory 用于创建模块实例。
type Factory func(db *gorm.DB) biz.Module

// factories 由 zz_modules_gen.go 的 init() 填充。
var factories []Factory

// All 返回当前项目中所有业务模块实例。
func All(db *gorm.DB) []biz.Module {
	res := make([]biz.Module, 0, len(factories))
	for _, f := range factories {
		res = append(res, f(db))
	}
	return res
}
