package exmaple1

import (
	"errors"
	"reflect"
	"unsafe"
)

func underlyingArr(slice []int) (*[4]int, error) {
	var arr *[4]int
	if len(slice) != 4 {
		return arr, errors.New("len of slice must be exactly 4")
	}
	p := unsafe.Pointer(&slice)
	h := (*reflect.SliceHeader)(p)
	arr = (*[4]int)(unsafe.Pointer(h.Data))
	return arr, nil
}
