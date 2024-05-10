package errors

var (
	Success       = NewError(200, "成功")
	BadParameters = NewError(400, "参数错误")
	Unauthorized  = NewError(401, "未授权")
	InternalError = NewError(500, "内部错误")
)
