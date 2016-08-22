package common

import (
	"os"
	"path/filepath"
)

// AppDir 应用目录绝对路径
func AppDir() string {
	return filepath.Dir(os.Args[0])
}
