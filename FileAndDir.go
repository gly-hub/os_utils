package os_utils

/***********************************************************
*              date:     20200716                          *
*              author:   vangogh                           *
*              describe:   os基础功能方法编写               *
*==========================================================*/

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

//@description 文件路径生成
func PathJoint(elem ...string) string {
	res_path := ""
	for _, str := range elem {
		if res_path != "" {
			res_path = fmt.Sprintf("%s/%s", res_path, str)
		} else {
			res_path = str
		}
	}
	return res_path
}

//@description 文件路径验证
func PathExist(path string) (bool, error) {
	ret := false
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return ret, err
}

//@description 文件内容读取
func ReadFile(path string) (string, error) {
	if isexist, err := PathExist(path); !isexist && err != nil {
		return "", err
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), err
}

//@description 文件写入 model模式只存在写入w、追加a
func WriteFile(path, content, model string) error {
	var file *os.File
	if isexist, err := PathExist(path); !isexist && err != nil {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
	} else {
		if model == "w" {
			file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		} else if model == "a" {
			file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0666)
		} else {
			err = errors.New("model模式只存在写入w、追加a。")
		}
		if err != nil {
			return err
		}
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()

	return nil
}

//@description 获取指定目录下的文件列表 list
func GetFileFromDir2LIST(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	file_list := make([]string, len(files))
	for k, file := range files {
		file_list[k] = file.Name()
	}
	return file_list, nil
}

type F struct {
	Name string
	Date string
	Size int64
}

//@description 获取指定目录下的文件列表 object
func GetFilesFromDir2OBJ(path string) ([]F, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	file_list := make([]F, len(files))
	for k, file := range files {
		f := F{}
		f.Name = file.Name()
		f.Date = file.ModTime().Format("2006-01-02 15:04:05")
		f.Size = file.Size() / 1024
		file_list[k] = f
	}
	return file_list, nil
}

//@description 文件移动
func MoveFile(oldpath, newpath string) error {
	//验证文件是否存在
	if exist, err := PathExist(oldpath); !exist && err != nil {
		return errors.New("oldpath不存在")
	}
	//验证新文件地址是否已存在
	if exist, err := PathExist(newpath); exist && err == nil {
		err := os.Remove(newpath)
		if err != nil {
			return err
		}
	}
	// 文件移动
	err := os.Rename(oldpath, newpath)
	if err != nil {
		return err
	}
	return nil
}
