package pdf

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

// OfficeToPdf word,excel等转换成pdf
func OfficeToPdf() {
	fileSrcPath := "/Users/admin/Downloads/test.xlsx" //自己机器上的文件地址
	outPath := "/Users/admin/Downloads"               //转出文件的路径
	fileType := "pdf"

	osName := runtime.GOOS //获取系统类型
	switch osName {
	case "darwin": //mac系统
		command := "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		pdfFile, err := funcDocs2Pdf(command, fileSrcPath, outPath, fileType)
		if err != nil {
			println("转化异常：", err.Error())
		}
		fmt.Println("转化后的文件：", pdfFile)
	case "linux":
		command := "libreoffice"
		pdfFile, err := funcDocs2Pdf(command, fileSrcPath, outPath, fileType)
		if err != nil {
			println("转化异常：", err.Error())
		}
		fmt.Println("转化后的文件：", pdfFile)
	case "windows":
		command := "soffice libreoffice" // 因为没有windows机器需要自己测试下这个命令行
		pdfFile, err := funcDocs2Pdf(command, fileSrcPath, outPath, fileType)
		if err != nil {
			println("转化异常：", err.Error())
		}
		fmt.Println("转化后的文件：", pdfFile)
	default:
		fmt.Println("暂时不支持的系统转化:" + runtime.GOOS)
	}
}

/**
*@tips libreoffice 转换指令：
* libreoffice6.2 invisible --convert-to pdf csDoc.doc --outdir /home/[转出目录]
*
* @function 实现文档类型转换为pdf或html
* @param command:libreofficed的命令(具体以版本为准)；win：soffice； linux：libreoffice6.2
*     fileSrcPath:转换文件的路径
*     fileOutDir:转换后文件存储目录
*     converterType：转换的类型pdf/html
* @return fileOutPath 转换成功生成的文件的路径 error 转换错误
 */

func funcDocs2Pdf(command string, fileSrcPath string, fileOutDir string, converterType string) (fileOutPath string, error error) {
	//校验fileSrcPath
	srcFile, erByOpenSrcFile := os.Open(fileSrcPath)
	if erByOpenSrcFile != nil && os.IsNotExist(erByOpenSrcFile) {
		return "", erByOpenSrcFile
	}
	//如文件输出目录fileOutDir不存在则自动创建
	outFileDir, erByOpenFileOutDir := os.Open(fileOutDir)
	if erByOpenFileOutDir != nil && os.IsNotExist(erByOpenFileOutDir) {
		erByCreateFileOutDir := os.MkdirAll(fileOutDir, os.ModePerm)
		if erByCreateFileOutDir != nil {
			fmt.Println("File ouput dir create error.....", erByCreateFileOutDir.Error())
			return "", erByCreateFileOutDir
		}
	}
	//关闭流
	defer func() {
		_ = srcFile.Close()
		_ = outFileDir.Close()
	}()
	//convert
	cmd := exec.Command(command, "--invisible", "--language=zh-CN", "--convert-to", converterType,
		fileSrcPath, "--outdir", fileOutDir)
	byteByStat, errByCmdStart := cmd.Output()
	//命令调用转换失败
	if errByCmdStart != nil {
		return "", errByCmdStart
	}
	//success
	fileOutPath = fileOutDir + "/" + strings.Split(path.Base(fileSrcPath), ".")[0]
	if converterType == "html" {
		fileOutPath += ".html"
	} else {
		fileOutPath += ".pdf"
	}
	fmt.Println("文件转换成功...", string(byteByStat))
	return fileOutPath, nil
}
