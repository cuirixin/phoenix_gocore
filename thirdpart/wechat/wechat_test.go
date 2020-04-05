package wechat

import (
	"fmt"
	"testing"

	"github.com/cuirixin/phoenix_gocore/thirdpart/wechat/cache"
)
func TestTime(tttt *testing.T) {

	wCache := cache.NewMemcache("39.106.145.5:11211")

	// opts := &cache.RedisOpts{
	// 	Host: "",
	// 	Password: "",
	// }
	// fmt.Println(1)
	// wCache := cache.NewRedis(opts)
	// fmt.Println(wCache)

	//配置微信参数
	config := &Config{
		AppID:          "wx534743a9806a1f4e",
		AppSecret:      "78d6c7090778c532580bd599aa1c4a84",
		Token:          "kUxrdlE7",
		EncodingAESKey: "bINnUplgPU34aO3VWphoenixg8BIbNLLntLxL92BrR4a6Cq9CX",
		Cache: wCache,
	}

	wc := NewWechat(config)

	// Test 1: JS-SDK
	js := wc.GetJs()
	cfg, _ := js.GetConfig("http://test-mall.91yummy.com/test")

	fmt.Println("AppID", cfg.AppID)
	fmt.Println("NonceStr", cfg.NonceStr)
	fmt.Println("Timestamp", cfg.Timestamp)
	fmt.Println("Signature", cfg.Signature)

	// Test 2: code换取access_token
	oauth := wc.GetOauth()
	fmt.Println(oauth)
	code := "testcode"
	// 通过code换取access_token
	resToken, err := oauth.GetUserAccessToken(code)
	fmt.Println(resToken, err.Error())
	// 拉取用户信息(需scope为 snsapi_userinfo)
	userInfo, err := oauth.GetUserInfo(resToken.AccessToken, resToken.OpenID)
	fmt.Println(userInfo, err)
}