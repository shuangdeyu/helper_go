package comhelper

import (
	"bufio"
	"encoding/base64"
	"log"
	"os"
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
