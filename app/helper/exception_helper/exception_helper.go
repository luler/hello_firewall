package exception_helper

import (
	"net/http"
)

type MyException struct {
	Message string
	Code    int
	Data    any
}

// 通用异常
func CommonException(data ...any) {
	myException := MyException{
		Message: "系统异常",
		Code:    http.StatusBadRequest,
		Data:    []int{},
	}
	dataLength := len(data)
	if dataLength >= 1 {
		myException.Message = data[0].(string)
	}
	if dataLength >= 2 {
		myException.Code = data[1].(int)
	}
	if dataLength >= 3 {
		myException.Data = data[2]
	}
	panic(myException)
}
