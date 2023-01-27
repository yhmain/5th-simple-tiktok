package util

// 该包定义各种error对应的常量及内容说明

import (
	"errors"
	"fmt"
)

//响应结构体
type Response struct {
	StatusCode int    `json:"status_code"`          //状态码
	StatusMsg  string `json:"status_msg,omitempty"` //返回状态描述
}

//供其他.go文件使用自定义报错信息
var (
	Success    = NewResponse(0, "Success")
	ServiceErr = NewResponse(1, "Service is unable to start successfully")

	// 用户注册遇到的错误
	TooLongInputErr     = NewResponse(10001, "输入长度不得超过32！")
	NameAlreadyExistErr = NewResponse(10002, "用户名已存在！")
	InsertNewUserErr    = NewResponse(10003, "新用户注册失败！")

	// 用户登录遇到的错误
	WrongPassword = NewResponse(11001, "用户登录失败！")
	WrongUserID   = NewResponse(11002, "错误的用户ID")

	// 用户Toekn遇到的错误
	GenTokenFailedErr = NewResponse(20001, "用户Token生成失败！")
	ParseTokenErr     = NewResponse(20002, "用户Token解析失败！")
	WrongTokenErr     = NewResponse(20003, "用户Token校验失败！")

	// 用户视频流遇到的错误
	InvalidTimeErr = NewResponse(30001, "非法的时间戳格式")
	FileParseErr   = NewResponse(30002, "文件解析失败！")
	FileUploadErr  = NewResponse(30003, "文件上传失败！")
	ImgUploadErr   = NewResponse(30004, "图片上传失败！")
	VideoUploadErr = NewResponse(30005, "视频上传失败！")
	VideoInsertErr = NewResponse(30006, "视频保存失败！")

	// ffmpeg截取视频帧遇到的错误
	ReadFrameErr = NewResponse(31001, "ffmpeg解析视频失败！")
	DecodeBufErr = NewResponse(31002, "imaging解析数据流失败！")

	// 赞操作遇到的错误
	FavActionErr = NewResponse(40001, "赞操作失败！")

	// Redis操作遇到的错误
	RedisConnErr = NewResponse(50001, "Redis连接失败！")
	RedisKeysErr = NewResponse(50002, "Redis的Keys命令失败！")
)

//返回一个错误信息的字符串
func (e Response) Error() string {
	return fmt.Sprintf("status_code=%d, status_msg=%s", e.StatusCode, e.StatusMsg)
}

//支持自定义一个Response结构体
func NewResponse(code int, msg string) Response {
	return Response{code, msg}
}

////支持自定义一个Response结构体（不带code）
func (e Response) WithMessage(msg string) Response {
	e.StatusMsg = msg
	return e
}

// ConvertErr convert error to Response（把系统的error类型转化成自定义的Response结构体类型）
func ConvertErr(err error) Response {
	Err := Response{}
	if errors.As(err, &Err) {
		return Err
	}
	s := ServiceErr
	s.StatusMsg = err.Error()
	return s
}
