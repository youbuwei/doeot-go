// Package boot pkg/boot/container.go
package boot

import (
	digen "github.com/youbuwei/doeot-go/gen/di"
	"go.uber.org/dig"

	"github.com/youbuwei/doeot-go/pkg/config"
)

// NewContainer 创建一个基础容器（方便以后测试等场景复用）
func NewContainer() *dig.Container {
	return dig.New()
}

// BuildContainer 是整个应用的 DI 入口：
// 1）注册框架层依赖（config、logger 等）；
// 2）调用业务侧生成的模块注入（internal/di/gen）。
func BuildContainer() *dig.Container {
	c := NewContainer()

	// 框架层注入
	provideFramework(c)

	// 业务侧注入（包含自动生成 + 手工扩展）
	digen.ProvideModules(c)

	return c
}

// provideFramework 注册框架级别的依赖（与具体业务无关）。
func provideFramework(c *dig.Container) {
	// 配置加载
	_ = c.Provide(config.Load)

	// 未来可以在这里加 logger、tracer、metrics 等：
	// _ = c.Provide(NewLogger)
	// _ = c.Provide(NewTracer)
}
