package test

import (
	"fmt"
	"testing"

	"github.com/feiyizhou/base/utils"
)

func Test_RandCharStr(t *testing.T) {
	fmt.Println(utils.RandStr(22, utils.AllCharArr))
}

func Test_RandNumStr(t *testing.T) {
	fmt.Println(utils.RandStr(8, utils.NumCharArr))
}
