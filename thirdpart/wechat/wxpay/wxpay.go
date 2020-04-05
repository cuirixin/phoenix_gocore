package wxpay

import (
	"fmt"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/cuirixin/phoenix_gocore/utils"
)

const bodyType = "application/xml; charset=utf-8"
const url_transfer_info = "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"
const url_unifiedorder = "https://api.mch.weixin.qq.com/pay/unifiedorder"
const url_refund = "https://api.mch.weixin.qq.com/secapi/pay/refund"
const url_orderquery = "https://api.mch.weixin.qq.com/pay/orderquery"
const url_refundquery = "https://api.mch.weixin.qq.com/pay/refundquery"

// 微信支付配置
type Config struct {
	AppId  string // 微信公众平台应用ID
	MchId  string // 微信支付商户平台商户号
	ApiKey string // 微信支付商户平台API密钥
	CertFile string
	KeyFile string
	RootCAFile string
	NotifyUrl string
}

// API客户端
type WxPay struct {
	config *Config
	stdClient *http.Client
	tlsClient *http.Client
}

// 实例化API客户端
func NewWxPay(config *Config) *WxPay {
	c := &WxPay{
		config: config,
		stdClient: &http.Client{},
	}

	// 附着商户证书
	err := c.WithCert(c.config.CertFile, c.config.KeyFile, c.config.RootCAFile)
	if err != nil {
		fmt.Println("创建链接失败", err)
	}
	return c
}

// 设置请求超时时间
func (c *WxPay) SetTimeout(d time.Duration) {
	c.stdClient.Timeout = d
	if c.tlsClient != nil {
		c.tlsClient.Timeout = d
	}
}

// 附着商户证书
func (c *WxPay) WithCert(certFile, keyFile, rootcaFile string) error {
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}
	rootca, err := ioutil.ReadFile(rootcaFile)
	if err != nil {
		return err
	}
	return c.WithCertBytes(cert, key, rootca)
}

func (c *WxPay) WithCertBytes(cert, key, rootca []byte) error {
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(rootca)
	if !ok {
		return errors.New("failed to parse root certificate")
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      pool,
	}
	trans := &http.Transport{
		TLSClientConfig: conf,
	}
	c.tlsClient = &http.Client{
		Transport: trans,
	}
	return nil
}

// 发送请求
func (c *WxPay) Post(url string, params Params, tls bool) (Params, error) {
	var httpc *http.Client
	if tls {
		if c.tlsClient == nil {
			return nil, errors.New("tls client is not initialized")
		}
		httpc = c.tlsClient
	} else {
		httpc = c.stdClient
	}
	resp, err := httpc.Post(url, bodyType, c.Encode(params))
	if err != nil {
		return nil, err
	}
	return c.Decode(resp.Body), nil
}

// XML解码
func (c *WxPay) Decode(r io.Reader) Params {
	var (
		d      *xml.Decoder
		start  *xml.StartElement
		params Params
	)
	d = xml.NewDecoder(r)
	params = make(Params)
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			start = &t
		case xml.CharData:
			if t = bytes.TrimSpace(t); len(t) > 0 {
				params.SetString(start.Name.Local, string(t))
			}
		}
	}
	return params
}

// XML编码
func (c *WxPay) Encode(params Params) io.Reader {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)
	return &buf
}

// 验证签名
func (c *WxPay) CheckSign(params Params) bool {
	return params.GetString("sign") == c.Sign(params)
}

// 生成签名
func (c *WxPay) Sign(params Params) string {
	var keys = make([]string, 0, len(params))
	for k, _ := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.GetString(k)) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params.GetString(k))
			buf.WriteString(`&`)
		}
	}
	buf.WriteString(`key=`)
	buf.WriteString(c.config.ApiKey)

	sum := md5.Sum(buf.Bytes())
	str := hex.EncodeToString(sum[:])

	return strings.ToUpper(str)
}

func (c*WxPay) GenParams(p Params) Params {
	params := make(Params)
	params.SetString("appid", c.config.AppId)
	params.SetString("mch_id", c.config.MchId)
	params.SetString("nonce_str", utils.RandomStr(10))  // 随机字符串
	for k, v := range p {
		params.SetString(k, v)
	}
	params.SetString("sign", c.Sign(params))
	return params
}

// 发送查询企业付款请求
// out_trade_no 商户调用企业付款API时使用的商户订单号
func (c *WxPay) QueryTransfer(out_transfer_no string) (Params, error) {
	params := c.GenParams(map[string]string{"partner_trade_no": out_transfer_no})
	ret, err := c.Post(url_transfer_info, params, true)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	return ret, err
}

const (
	DEVICE_INFO_WEB = "web"
	TRADE_TYPE_JSAPI = "JSAPI" // 公众号支付
	TRADE_TYPE_NATIVE = "NATIVE" // 原生扫码支付
	TRADE_TYPE_APP = "APP" // APP支付
)

// 统一下单
// return prepay_id 
// map[
// 	prepay_id:xxx 统一下单ID
// ]
func (c *WxPay) Charge(openid, amount_fen, out_trade_no, body, device_info, trade_type, client_ip string) (Params, error) {
	params := c.GenParams(map[string]string{
		"out_trade_no": out_trade_no, // 商户订单号
		"device_info": device_info, // 设备信息，默认web
		"body": body, // 支付订单信息
		"fee_type": "CNY",
		"trade_type": trade_type,
		"notify_url": c.config.NotifyUrl,
		"openid": openid,
		"spbill_create_ip": client_ip,
		"total_fee": amount_fen,
	})
	ret, err := c.Post(url_unifiedorder, params, true)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	return ret, err
}

// 支付订单查询
// return prepay_id 统一下单ID
// map[
// 	return_code:SUCCESS
// 	result_code:SUCCESS
// 	trade_state_desc:订单未支付
// 	return_msg:OK
// 	trade_state:NOTPAY 支付状态
// ]
func (c *WxPay) QueryCharge(out_trade_no string) (Params, error) {
	params := c.GenParams(map[string]string{
		"out_trade_no": out_trade_no, // 商户订单号
	})
	ret, err := c.Post(url_orderquery, params, true)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	return ret, err
}

// 退款
// map[
// 	refund_id: xxxx 退款流水号
// ]
func (c *WxPay) Refund(pay_transaction_id, pay_amount_fen, out_refund_no, amount_fen string) (Params, error) {
	params := c.GenParams(map[string]string{
		"transaction_id": pay_transaction_id, // 微信支付订单的流水号
		"out_refund_no": out_refund_no, // 退款商户订单号
		"total_fee": pay_amount_fen, // 订单总金额
		"refund_fee": amount_fen, // 支付订单信息
		"refund_fee_type": "CNY",
		"op_user_id": c.config.MchId, // 操作员帐号, 默认为商户号
	})
	ret, err := c.Post(url_refund, params, true)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	return ret, err
}

// 退款订单查询
// return prepay_id 统一下单ID
// map[
// 	return_code:SUCCESS
// 	result_code:SUCCESS
// 	trade_state_desc:订单未支付
// 	return_msg:OK
// 	trade_state:NOTPAY 退款状态
// ]
func (c *WxPay) QueryRefund(refund_id string) (Params, error) {
	params := c.GenParams(map[string]string{
		"refund_id": refund_id, // 退款流水号
	})
	ret, err := c.Post(url_refundquery, params, true)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	return ret, err
}


