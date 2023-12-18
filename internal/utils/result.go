package utils

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func ResultSuccess(data string) Result {
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

func ResultFailD(code int, msg string, data string) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
