/**
 * @Author: victor
 * @Description:
 * @File:  consumer_demo.gp
 * @Version: 1.0.0
 * @Date: 2020/4/5 11:49 上午
 */

package main

import (
	"fmt"
	"path"
	"time"
	"github.com/gogap/logs"

	"github.com/cuirixin/phoenix_gocore/thirdpart/aliyun/ali_mns"
)

const queue = "test-common"

var MnsConsumerClient ali_mns.MNSClient


func main(){

	queueN := ali_mns.NewMNSQueue(queue, MnsConsumerClient)

	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		for {
			select {
			case resp := <-respChan:
				{
					logs.Pretty("response:", resp)
					logs.Debug("response body:", string(resp.MessageBody[:]))
					logs.Debug("change the visibility: ", resp.ReceiptHandle)
					if ret, e := queueN.ChangeMessageVisibility(resp.ReceiptHandle, 5); e != nil {
						logs.Error(e)
					} else {
						logs.Pretty("visibility changed", ret)
					}

					// TODO 业务代码

					// logs.Debug("delete it now: ", resp.ReceiptHandle)
					// if e := queueN.DeleteMessage(resp.ReceiptHandle); e != nil {
					// 	logs.Error(e)
					// }
				}
			case err := <-errChan:
				{
					if err!=nil {
						// logs.Error(err)
						time.Sleep(time.Second * 5)
					}
				}
			}
		}

	}()

	queueN.ReceiveMessage(respChan, errChan)
	for {
		time.Sleep(time.Second * 1)
	}

}

