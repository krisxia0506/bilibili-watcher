package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 定义了标准的 API 响应结构体。
// Code: 业务状态码, 0 表示成功, 其他表示失败。
// Msg: 响应消息。
// Data: 响应数据 (必须字段，如果 Go 值为 nil，JSON 中会是 null)。
type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"` // 使用 any (Go 1.18+), 移除 omitempty
}

// 定义一些常用的业务状态码 (可以根据需要扩展)
const (
	CodeSuccess       = 0
	CodeError         = -1  // 通用错误
	CodeInvalidParams = 1   // 参数无效
	CodeNotFound      = 2   // 资源未找到
	CodeUnauthorized  = 3   // 未授权
	CodeInternalError = 500 // 内部服务器错误 (与 HTTP 500 对应)
	// ... 其他自定义错误码
)

// Success 返回一个表示成功的 API 响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMessage 返回一个带自定义消息的成功 API 响应。
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code: CodeSuccess,
		Msg:  message,
		Data: data,
	})
}

// Error 返回一个表示失败的 API 响应。
// httpStatusCode 是实际的 HTTP 状态码 (例如 http.StatusBadRequest)。
// businessCode 是自定义的业务错误码。
// message 是错误信息。
func Error(c *gin.Context, httpStatusCode int, businessCode int, message string) {
	c.JSON(httpStatusCode, APIResponse{
		Code: businessCode,
		Msg:  message,
		Data: nil, // Data 字段现在会输出为 null
	})
}

// ErrorWithData 返回一个带数据的失败 API 响应。
func ErrorWithData(c *gin.Context, httpStatusCode int, businessCode int, message string, data interface{}) {
	c.JSON(httpStatusCode, APIResponse{
		Code: businessCode,
		Msg:  message,
		Data: data,
	})
}

// QuickSuccess 是 Success 的快捷方式，如果 data 为 nil，则不包含 data 字段。
func QuickSuccess(c *gin.Context, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}
	// 如果 responseData 为 nil (即没有传入 data), Data 字段会是 null
	c.JSON(http.StatusOK, APIResponse{
		Code: CodeSuccess,
		Msg:  "success",
		Data: responseData,
	})
}

// QuickError 是 Error 的快捷方式，使用通用错误码和消息。
// httpStatusCode 通常是 http.StatusBadRequest 或 http.StatusInternalServerError。
func QuickError(c *gin.Context, httpStatusCode int, message string) {
	businessCode := CodeError
	if httpStatusCode == http.StatusInternalServerError {
		businessCode = CodeInternalError
	}
	c.JSON(httpStatusCode, APIResponse{
		Code: businessCode,
		Msg:  message,
		Data: nil, // Data 字段会输出为 null
	})
}
