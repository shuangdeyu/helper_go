package comhelper

import (
	"bufio"
	"encoding/base64"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
)

/**
 * 检查文件夹是否存在
 */
func IsDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 * 将图片转换为base64格式
 */
func ImgToBase64(path string) (string, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer imgFile.Close()

	//create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	//read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	//convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64str := base64.StdEncoding.EncodeToString(buf)
	return "data:image/jpeg;base64," + imgBase64str, nil
}

// 获取当前目录
func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) // 返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	return strings.Replace(dir, "\\", "/", -1)       //将\替换成/
}

// 判断目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 正序读取文件
func FileReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// 倒序读取文件
func FileReverseRead(name string, lineNum uint) ([]string, error) {
	// 打开文件
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fs.Size()

	var offset int64 = -1   // 偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) // 用于读取单个字节
	lineStr := ""           // 存放一行的数据
	buff := make([]string, 0, 100)
	for (-offset) <= fileSize {
		// 通过Seek函数从末尾移动游标然后每次读取一个字节
		file.Seek(offset, io.SeekEnd)
		_, err := file.Read(char)
		if err != nil {
			return buff, err
		}
		if char[0] == '\n' {
			offset--  // windows跳过'\r'
			lineNum-- // 到此读取完一行
			if lineStr != "" {
				buff = append(buff, lineStr)
			}
			lineStr = ""
			if lineNum == 0 {
				return buff, nil
			}
		} else {
			lineStr = string(char) + lineStr
		}
		offset--
	}
	buff = append(buff, lineStr)
	return buff, nil
}

/**********************************************************************************
 * ini配置文件读取
 */
var (
	iniCfg *goconfig.ConfigFile
)

func InitIni(path string) {
	var tmpErr error
	iniCfg, tmpErr = goconfig.LoadConfigFile(path)
	if tmpErr != nil {
		panic("读取配置文件失败")
	}
}

func LoadIni(param1 string, param2 string) string {

	result, err := iniCfg.GetValue(param1, param2)
	if err != nil {
		log.Fatal("无法获取键值", err)
	}
	return result
}
