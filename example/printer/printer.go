package printer

/*
#cgo linux LDFLAGS: -lcups
#cgo darwin LDFLAGS: -lcups
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <cups/cups.h>

typedef struct {
    char** printerNames;
    int size;
} PrinterNamesArray;

PrinterNamesArray getPrinterNames() {
    PrinterNamesArray result;
    result.size = 0;
    result.printerNames = NULL;

    cups_dest_t *dests;
    int num_dests = cupsGetDests(&dests);

    if (num_dests > 0) {
        result.printerNames = (char**)malloc(num_dests * sizeof(char*));
        result.size = num_dests;

        for (int i = 0; i < num_dests; i++) {

            result.printerNames[i] = strdup(dests[i].name);
        }
    }

	cupsFreeDests(num_dests, dests);
	return result;
}

void freePrinterNamesArray(PrinterNamesArray printerNamesArray) {
	for (int i = 0; i < printerNamesArray.size; i++) {
		free(printerNamesArray.printerNames[i]);
	}
	free(printerNamesArray.printerNames);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// 获取本地打印机列表
func GetPrinters() []string {

	result := C.getPrinterNames()
	fmt.Println(result)
	fmt.Println(result.size)
	fmt.Println(result.printerNames)

	// Convert a C string to a Go string
	printerName := cStringsToGoStrings(int(result.size), result.printerNames)
	fmt.Printf("Printer: %s\n", printerName)

	// Free the memory of the array
	C.freePrinterNamesArray(result)

	return printerName

}

// 将char**转换为Go语言的字符串切片
func cStringsToGoStrings(size int, cStrings **C.char) []string {
	var goStrings []string

	// Iterate through the C string array and convert it to a Go language string
	for i := 0; i < size; i++ {
		// Gets a pointer to a C string
		cString := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cStrings)) + uintptr(i)*unsafe.Sizeof(*cStrings)))

		// Checks if the pointer is NULL, indicating the end of the string array
		// After the actual array is taken, it can also get the address of cString, but it is a null pointer,
		// so it is taken in the size range
		if cString == nil {
			break
		}

		goStrings = append(goStrings, C.GoString(cString))
	}
	return goStrings
}
