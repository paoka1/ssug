package utils

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ResultSuccess(data interface{}) Result {
	return Result{
		Code: 0,
		Msg:  "操作成功",
		Data: data,
	}
}

func ResultFail(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: "",
	}
}

// ResultFailWD Result Fail With Data
func ResultFailWD(code int, msg string, data interface{}) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
