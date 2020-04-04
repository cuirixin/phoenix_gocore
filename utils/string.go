/**********************************************
** @Des: This file ...
** @Author: victor
** @Date:   2017-12-12 10:10:00
** @Last Modified by:   victor
** @Last Modified time: 2017-12-12 10:10:00
***********************************************/

package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
	"sort"
	"net/url"
)

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func SizeFormat(size float64) string {
	units := []string{"Byte", "KB", "MB", "GB", "TB"}
	n := 0
	for size > 1024 {
		size /= 1024
		n += 1
	}

	return fmt.Sprintf("%.2f %s", size, units[n])
}

func Password(len int, pwdO string) (pwd string, salt string) {
	salt = RandomStr(6)
	defaultPwd := "tamphoenix"
	if pwdO != "" {
		defaultPwd = pwdO
	}
	pwd = Md5([]byte(defaultPwd + salt))
	return pwd, salt
}

// 生成32位MD5
// func MD5(text string) string{
//    ctx := md5.New()
//    ctx.Write([]byte(text))
//    return hex.EncodeToString(ctx.Sum(nil))
// }

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// #string到int  
// int,err:=strconv.Atoi(string)  
// #string到int64  
// int64, err := strconv.ParseInt(string, 10, 64)  
// #int到string  
// string:=strconv.Itoa(int)  
// #int64到string  
// string:=strconv.FormatInt(int64,10)  

// Map按key正序排序后拼接url
func ParamsSortToUrl(params map[string]string, excludeParams []string) string {
	var keys []string  
	for k := range params {
		exits := false
		for i := range excludeParams {
			if excludeParams[i] == k {
				exits = true
				break
			}
		}
		if !exits {
			keys = append(keys, k)  
		}
	}  
	sort.Strings(keys)
	u := url.Values{}
	for _, k := range keys {
		fmt.Println("Key:", k, "Value:", params[k])  
		u.Set(k, params[k])
	}

	// fmt.Println(u.Get("id"))
	// fmt.Println(u.Add("id", "1"))

	return u.Encode() // a=A&c=C
}