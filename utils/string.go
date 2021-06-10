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
	"encoding/base64"
	"fmt"
	"github.com/jameskeane/bcrypt"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"time"
	"unsafe"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func AdvancedRandomStr(n int) string {
	// 参考 https://www.toutiao.com/i6861018877808083469/?tt_from=copy_link&utm_campaign=client_share&timestamp=1605235294&app=news_article_lite&utm_source=copy_link&utm_medium=toutiao_ios&use_new_style=1&req_id=2020111310413301012903213005014893&group_id=6861018877808083469
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
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

/**
 * string转换int
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func ParseInt(b string, defInt int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return defInt
	} else {
		return id
	}
}

/**
 * int转换string
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func ParseString(b int) string {
	id := strconv.Itoa(b)
	return id
}

/**
 * 转换浮点数为string
 * @method func
 * @param  {[type]} t *             Tools [description]
 * @return {[type]}   [description]
 */
func ParseFlostToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 5, 64)
}

/**
 * 字符串截取
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func SubString(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	lenStr := len(runeStr)

	if start < 0 {
		start = lenStr + start
	}
	if start > lenStr {
		start = lenStr
	}
	end := start + length
	if end > lenStr {
		end = lenStr
	}
	if length < 0 {
		end = lenStr + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}

/**
 * base64 解码
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func Base64Decode(str string) string {
	s, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(s)
}

func HashPassword(pwd string) string {
	salt, err := bcrypt.Salt(10)
	if err != nil {
		return ""
	}
	hash, err := bcrypt.Hash(pwd, salt)
	if err != nil {
		return ""
	}

	return hash
}

func MachPassword(password, hash string) bool {
	return bcrypt.Match(password, hash)
}
