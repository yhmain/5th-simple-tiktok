package util

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	accessKey = "ByyFFHWJCJRUaJpbquB9d8Nyxuc0O3sPp3waR6ku"
	secretKey = "ea6nktNTugz795s4ZyNXtHsmPA-p5DhnIIhqAFRj"
	bucket    = "simple-tiktok"
	Refer     = "http://rp1950nfc.hn-bkt.clouddn.com/"
)

// 上传文件到七牛云
func UploadToQiNiu(file multipart.File, fileSize int64, filepath string) (string, Response) {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	// 指定地区，不使用cdn，https
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.Put(context.Background(), &ret, upToken, filepath, file, fileSize, &putExtra)
	if err != nil {
		fmt.Println(err)
		return "Error", FileUploadErr
	}
	url := Refer + ret.Key
	return url, Success

}

// 自动生成文件名
func GenFileName() string {
	unix := time.Now().Unix()
	timeByte := []byte(strconv.Itoa(int(unix)))
	md5Str := md5.Sum(timeByte) // 这里是计算了一个md5(time())的字符串作为文件名
	filename := fmt.Sprintf("%x", md5Str)
	return filename
}

// 上传图片到七牛
func UploadImgToQiNiu(file io.Reader, fileSize int64, filepath string) (string, Response) {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	// 指定地区，不使用cdn，https
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.Put(context.Background(), &ret, upToken, filepath, file, fileSize, &putExtra)
	if err != nil {
		fmt.Println(err)
		return "Error", FileUploadErr
	}
	url := Refer + ret.Key
	return url, Success

}
