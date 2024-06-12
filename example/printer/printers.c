// printers.c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <cups/cups.h>

// 定义一个结构体，用于存储打印机名字的数组和数组的大小
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
        // 分配存储打印机名字的数组
        result.printerNames = (char**)malloc(num_dests * sizeof(char*));
        result.size = num_dests;

        for (int i = 0; i < num_dests; i++) {
            printf("dest: %s\n", dests[i].name);

             // 分配每个打印机名字的内存，并将其复制到数组中
            result.printerNames[i] = strdup(dests[i].name);
        }
    } else {
        printf("No printers found.\n");
    }

    /*result.printerNames = (char**)malloc(2 * sizeof(char*));
    result.size = 2;
    result.printerNames[0] = strdup("sasasasas");
    result.printerNames[1] = strdup("klkllklli");*/

    cupsFreeDests(num_dests, dests);

    return result;
}

void freePrinterNamesArray(PrinterNamesArray printerNamesArray) {
    for (int i = 0; i < printerNamesArray.size; i++) {
        free(printerNamesArray.printerNames[i]);
    }
    free(printerNamesArray.printerNames);
}
