package printer

import "C"
import (
	"fmt"
	"github.com/phin1x/go-ipp"
)

/**
 * 调用本地打印机，打印
 */
func PrinterIPP() {
	// 创建 IPP 客户端
	client := ipp.NewIPPClient("localhost", 631, "", "", false)
	//client.GetPrinterAttributes()

	err := client.TestConnection()
	if err != nil {
		fmt.Printf("%v", err)
	} else {
		var res []string
		printer, _ := client.GetPrinterAttributes("HP_Color_LaserJet_MFP_M181fw__F51496_", res)
		for _, atr := range printer {
			for _, nameAtr := range atr {
				fmt.Println(nameAtr)
			}
		}
	}
}
