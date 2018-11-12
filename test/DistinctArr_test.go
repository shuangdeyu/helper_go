package test

import (
	"fmt"
	"helper_go/comhelper"
	"testing"
)

func TestDistinctArr(t *testing.T) {
	a := []string{"hello", "good", "world", "yes", " ", "hello", "nihao", "shijie", "hello", "yes", "nihao", "good"}
	fmt.Println(a)
	fmt.Println(comhelper.DistinctArrString(a))
}
