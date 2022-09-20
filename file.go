package gutils

import (
	"errors"
	"os"
)

/*
判断文件是否存在
存在则返回 true
不存在则返回  false
*/
func PathExist(path string) (bool, error) {
	//读取文件
	fs, err := os.Stat(path)
	//读取文件正常
	if err == nil {
		//文件已经存在
		if fs.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
