package shared

import (
	"os"
	"path/filepath"
)

// WriteFileOnce 如果文件不存在则写入，否则跳过。
func WriteFileOnce(path string, content []byte) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}
