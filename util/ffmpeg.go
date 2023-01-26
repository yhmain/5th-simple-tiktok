package util

import (
	"bytes"
	"fmt"
	"os"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// 参考github示例  https://github.com/u2takey/ffmpeg-go
//从视频中抽取第 frameNum 帧作为封面
func GetVideoFrameQiNiu(inFileName string, frameNum int) (string, Response) {
	// 使用ffmpeg
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil { // ffmpeg获取数据流
		return "Error", ReadFrameErr
	}
	// 上传到七牛云			???待解决
	url, resp := UploadImgToQiNiu(buf, int64(len(buf.Bytes())), GenFileName())
	if url == "Error" {
		return "Error", ImgUploadErr
	}
	return url, resp
}

//从视频中抽取第 frameNum 帧作为封面
func GetVideoFrame(inFileName, outFileName string, frameNum int) (string, Response) {
	// 使用ffmpeg
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil { // ffmpeg获取数据流
		return "Error", ReadFrameErr
	}
	reader := buf
	img, err := imaging.Decode(reader)
	if err != nil {
		return "Error", DecodeBufErr
	}
	err = imaging.Save(img, outFileName)
	if err != nil {
		return "Error", ImgUploadErr
	}
	return outFileName, Success
}
