package test

import (
	"fmt"
	"github.com/shuangdeyu/helper_go/comhelper"
	"testing"
)

func TestInArray(t *testing.T) {
	arr_array := [5]string{"a", "b", "c", "d", "e"}
	val_array := "c"
	ret := comhelper.InArray(arr_array, val_array)
	fmt.Println(ret)

	arr_slice := []string{"a", "b", "c", "d", "e"}
	val_slice := "f"
	ret = comhelper.InArray(arr_slice, val_slice)
	fmt.Println(ret)

	arr_map := map[string]string{"a": "man", "b": "women"}
	val_map := "b"
	ret = comhelper.InArray(arr_map, val_map)
	fmt.Println(ret)
}
