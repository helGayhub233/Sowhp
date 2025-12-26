package concert

import (
	"Sowhp/concert/logger"
	"fmt"
	"os"
)

func MkdirResport() {
	dir, err := os.Getwd()
	if err != nil {
		logger.DebugError(err)
	}
	logger.Debug(fmt.Sprintf("当前路径:%s", dir))
	Dir_mk(dir + "/result/")
	logger.Debug("当前目录下创建 /result/ 目录成功！")

}

func GetPath() string {
	dir, err := os.Getwd()
	if err != nil {
		logger.DebugError(err)
	}
	return dir
}

func Dir_mk(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			logger.DebugError(err)
		}
		logger.Debug(fmt.Sprintf("当前已创建该目录：%s", path))
		return
	}
}
