package windows

/*
   #cgo windows LDFLAGS: -lwinspool
   #include <stdio.h>
   #include <windows.h>

   typedef struct {
       char** printerNames;
       int size;
   } PrinterNamesArray;

   PrinterNamesArray getPrinterNames() {
       printf("Error getting buffer size\n");

   	PrinterNamesArray result;
       result.size = 0;
       result.printerNames = NULL;

       DWORD dwNeeded, dwReturned;
       PRINTER_INFO_2 *pPrinterInfo;

       // The buffer size needed to get printer information
       if (!EnumPrinters(PRINTER_ENUM_LOCAL, NULL, 2, NULL, 0, &dwNeeded, &dwReturned)) {
           if (GetLastError() != ERROR_INSUFFICIENT_BUFFER) {
               printf("Error getting buffer size\n");
               return result;
           }
       }

       // Allocate a large enough buffer
       pPrinterInfo = (PRINTER_INFO_2 *)malloc(dwNeeded);
       if (pPrinterInfo == NULL) {
           printf("Error allocating buffer\n");
           return result;
       }

       // Get a list of printers
       if (!EnumPrinters(PRINTER_ENUM_LOCAL, NULL, 2, (LPBYTE)pPrinterInfo, dwNeeded, &dwNeeded, &dwReturned)) {
           printf("Error getting printer list\n");
           free(pPrinterInfo);
           return result;
       }

       // Print printer list
       // printf("Printers:\n");
       // for (DWORD i = 0; i < dwReturned; i++) {
       //     printf("%s\n", pPrinterInfo[i].pPrinterName);
       // }

   	if (dwReturned > 0) {
           result.printerNames = (char**)malloc(dwReturned * sizeof(char*));
           result.size = dwReturned;

           for (DWORD i = 0; i < dwReturned; i++) {

               result.printerNames[i] = strdup(pPrinterInfo[i].pPrinterName);
           }
       }

       free(pPrinterInfo);

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

func main() {
	fmt.Println(GetPrinters())
}

func GetPrinters() []string {

	result := C.getPrinterNames()
	//fmt.Println(result)
	//fmt.Println(result.size)
	//fmt.Println(result.printerNames)

	// Convert a C string to a Go string
	printerName := cStringsToGoStrings(int(result.size), result.printerNames)
	//fmt.Printf("Printer: %s\n", printerName)

	// Free the memory of the array
	C.freePrinterNamesArray(result)

	return printerName

}

// Convert char** to a string slice in Go
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
