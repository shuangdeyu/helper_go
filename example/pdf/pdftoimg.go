package pdf

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"os"
)

// 将PDF转换为图片
func PdfToImage() {
	doc, err := fitz.New("/Users/admin/Downloads/testpdf.pdf")
	if err != nil {
		panic(err)
	}

	defer doc.Close()

	/*tmpDir, err := os.MkdirTemp(os.TempDir(), "fitz")
	if err != nil {
		panic(err)
	}*/

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		// 转jpg
		/*img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join("", fmt.Sprintf("test%03d.jpg", n)))
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}
		f.Close()*/

		// 转png，可直接操作字节流[]byte
		png, err := doc.ImagePNG(n, 300)
		if err != nil {
			panic(err)
		}

		fileObj, err := os.OpenFile(fmt.Sprintf("test%03d.png", n), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Failed to open the file", err.Error())
			os.Exit(2)
		}
		if _, err := fileObj.Write(png); err == nil {
			fmt.Println("Successful writing to thr file with os.OpenFile and *File.Write method.", png)
		}
	}
}
