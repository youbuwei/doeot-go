// Package di internal/di/container.go
package di

import "go.uber.org/dig"

// ProvideManual 给业务开发者预留的“手动注入 / 覆盖”扩展点。
// 默认什么都不做，有特殊需求再在这里加 Provide。
func ProvideManual(c *dig.Container) {
	// 示例（需要时自己加）：
	// _ = c.Provide(NewCustomLogger)
	// _ = c.Provide(NewFeatureFlagClient)
}
