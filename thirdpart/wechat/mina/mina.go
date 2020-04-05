package mina

import (
	"encoding/json"
	"fmt"
	_ "time"

	"github.com/cuirixin/phoenix_gocore/thirdpart/wechat/context"
	"github.com/cuirixin/phoenix_gocore/utils"
)

const jscode2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

// Mina struct
type Mina struct {
	*context.Context
}

// resSession 请求返回结果
type resSession struct {
	utils.CommonError

	Openid    string `json:"openid"`
	SessionKey string  `json:"session_key"`
	Unionid	 string `json:"unionid"`
}

//NewMina init
func NewMina(context *context.Context) *Mina {
	mina := new(Mina)
	mina.Context = context
	return mina
}

//GetSessionByJscode 从服务器中code置换用户session
func (mina *Mina) GetSessionByJscode(code string) (session resSession, err error) {
	var response []byte
	url := fmt.Sprintf(jscode2sessionURL, mina.Context.AppID, mina.Context.AppSecret, code)
	response, err = utils.HTTPGet(url)
	err = json.Unmarshal(response, &session)
	fmt.Println(err)
	if err != nil {
		return session, err
	}
	if session.ErrCode != 0 {
		err = fmt.Errorf("GetSessionByJscode Error : errcode=%d , errmsg=%s", session.ErrCode, session.ErrMsg)
		return session, err
	}
	return session, nil
}
