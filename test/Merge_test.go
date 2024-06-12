package test

import (
	"fmt"
	"github.com/shuangdeyu/helper_go/comhelper"
	"testing"
)

func TestMergeString(t *testing.T) {
	arr1 := []string{"apple", "nuojiya", "heimei", "google"}
	arr2 := []string{"facebook", "twitter", "tumbler"}
	arr := comhelper.MergeString(arr1, arr2)
	fmt.Println(arr)
}
