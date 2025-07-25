package kucoin

import (
	"context"
	"github.com/Kucoin/kucoin-go-sdk"
	"pledge-backend-study/db"
	"pledge-backend-study/log"
	"time"
)

const ApiKeyVersionV2 = "2"

// 默认加个
var PlgrPrice = "0.0027"

// 这个是创建一个创建缓冲区大小为2的通道

var PlgrPriceChan = make(chan string, 2)

func GetExchangePrice() {
	log.Logger.Info("GetExchangePrice start")
	//先从redis获取数据
	price, err := db.RedisGetString("plgr_price")
	if err != nil {
		log.Logger.Sugar().Error("get plgr price from redis err ", err)
	} else {
		PlgrPrice = price
	}
	s := kucoin.NewApiService(
		kucoin.ApiKeyOption("key"),
		kucoin.ApiSecretOption("secret"),
		kucoin.ApiPassPhraseOption("passphrase"),
		kucoin.ApiKeyVersionOption(ApiKeyVersionV2),
	)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	rsp, err := s.WebSocketPublicToken(ctx)

	if err != nil {
		log.Logger.Error("GetExchangePrice err ")
		return
	}

	tk := &kucoin.WebSocketTokenModel{}

	if err := rsp.ReadData(tk); err != nil {
		log.Logger.Error("GetExchangePrice err ")
		return
	}
	c := s.NewWebSocketClient(tk)

	connect, errors, err := c.Connect()
	if err != nil {
		log.Logger.Sugar().Errorf("Error: %s", err.Error())
		return
	}
	//创建一个 订阅消息（subscribe message），用于告诉 KuCoin 的 WebSocket 服务器：我想要订阅 /market/ticker:PLGR-USDT 这个行情主题。
	ch := kucoin.NewSubscribeMessage("/market/ticker:PLGR-USDT", false)
	//取消订阅
	uch := kucoin.NewUnsubscribeMessage("/market/ticker:PLGR-USDT", false)

	if err := c.Subscribe(ch); err != nil {
		log.Logger.Sugar().Errorf("Error: %s", err.Error())
		return
	}

	for {
		select {
		case err := <-errors:
			c.Stop()
			log.Logger.Sugar().Errorf("Error: %s", err.Error())
			//取消订阅
			c.Unsubscribe(uch)
			return
		case msg := <-connect:
			t := &kucoin.TickerLevel1Model{}
			if err := msg.ReadData(t); err != nil {
				log.Logger.Sugar().Errorf("Failure to read: %s", err.Error())
				return
			}
			PlgrPriceChan <- t.Price
			PlgrPrice = t.Price
			db.RedisSetString("plgr_price", PlgrPrice, 0)
		}
	}

}
