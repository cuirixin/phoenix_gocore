package wechat

import (
	"fmt"

	"github.com/cuirixin/phoenix_gocore/thirdpart/gopay"
	"github.com/cuirixin/phoenix_gocore/thirdpart/gopay/wechat"
)

func Reverse() {
	//初始化微信客户端
	//    appId：应用ID
	//    MchID：商户ID
	//    ApiKey：Key值
	//    isProd：是否是正式环境
	client := wechat.NewClient("wxdaa2ab9ef87b5497", "1368139502", "GFDS8j98rewnmgl45wHTt980jg543abc", false)

	// 初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gopay.GetRandomString(32))
	bm.Set("out_trade_no", "6aDCor1nUcAihrV5JBlI09tLvXbUp02B")
	bm.Set("sign_type", wechat.SignType_MD5)

	//请求撤销订单，成功后得到结果，沙箱环境下，证书路径参数可传空
	wxRsp, err := client.Reverse(bm, "", "", "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response:", wxRsp)
}
