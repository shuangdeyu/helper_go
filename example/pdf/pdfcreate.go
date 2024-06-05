package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/signintech/gopdf"
	"log"
)

func Gofpdf() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	// 英文-系统自带的
	//pdf.SetFont("Arial", "BIUS", 16)
	// 中文字体加载，适用于英文、日文等
	pdf.AddUTF8Font("simsun", "", "./pdf/ttf/simsun.ttf")
	pdf.SetFont("simsun", "", 16)

	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(190, 7, "Welcome こんにちは 你好", "", 0, "CM", false, 0, "")

	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	/*pdf.ImageOptions(
		"topgoer.png",
		80, 20,
		0, 0,
		false,
		gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
		0,
		"",
	)*/

	err := pdf.OutputFileAndClose("hello.pdf")
	fmt.Println(err)
}

func Gopdf() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err := pdf.AddTTFFont("simsun", "./pdf/ttf/simsun.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = pdf.SetFont("simsun", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}
	//pdf.SetGrayFill(0.5)
	pdf.Cell(nil, "Welcome こんにちは 你好呀!")
	pdf.WritePdf("hello.pdf")
}
