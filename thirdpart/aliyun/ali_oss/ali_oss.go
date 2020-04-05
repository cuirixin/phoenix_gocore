/**********************************************
** @Des: 阿里云OSS操作库
** @Author: victor
** @Date:   2017-12-12 10:10:00
** @Last Modified by:   victor
** @Last Modified time: 2017-12-12 10:10:00
***********************************************/

package ali_oss

import (
	"fmt"
	//"io/ioutil"
	//"strings"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSS struct {
	Client *oss.Client
}

func NewAliOSS(endpoint, accessKeyId, secretAccessKey string) *OSS {
	ossClient, err := oss.New(endpoint, accessKeyId, secretAccessKey)
	if err!=nil {
		panic(err)
	}
	return &OSS{Client: ossClient}
}

func (oss *OSS) PutOssObjectFromFile(bucket, objectPath, filePath string) bool {
	ossBucket, _ := oss.Client.Bucket(bucket)
	err := ossBucket.PutObjectFromFile(objectPath, filePath)
	if err!=nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (oss *OSS) GetOssObjectToFile(bucket, objectPath, filePath string) bool {
	ossBucket, _ := oss.Client.Bucket(bucket)
	err := ossBucket.GetObjectToFile(objectPath, filePath)
	if err!=nil {
		fmt.Println(err)
		return false
	}
	return true
}





// func main() {
	// client, _ := oss.New("http://oss-cn-hangzhou.aliyuncs.com",
	// 		"YourAccessKeyId",
	// 		"YourAccessKeySecret")
	// bucket, _ := client.Bucket("my-bucket")
	// // 字符串上传下载
	// err := bucket.PutObject("my-object-1", strings.NewReader("Hello Oss"))
	// rd, err := bucket.GetObject("my-object-1")
	// data, err := ioutil.ReadAll(rd)
	// rd.Close()
	// fmt.Println(string(data))
	// // 文件上传下载
	// err = bucket.PutObjectFromFile("my-object-2", "mypic.jpg")
	// err = bucket.GetObjectToFile("my-object-2", "mynewpic.jpg")
	// // 分片并发，断点续传上传/下载
	// err = bucket.UploadFile("my-object-3", "mypic.jpg", 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	// err = bucket.DownloadFile("my-object-3", "mynewpic.jpg", 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	// // 查看Object
	// lsRes, err := bucket.ListObjects()
	// fmt.Println("my objects:", lsRes.Objects)
	// // 上面的err都需要处理，此处略
	// if err != nil {
	// 		// TODO
	// }
// }