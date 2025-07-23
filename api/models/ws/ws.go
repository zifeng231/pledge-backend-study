package ws

import "pledge-backend-study/log"

// todo websocket
func StartServer() {
	log.Logger.Info("WsServer start")
	for {
		select {
		case price, ok := <-kucoin.PlgrPriceChan:
			if ok {
				Manager.Servers.Range(func(key, value interface{}) bool {
					value.(*Server).SendToClient(price, SuccessCode)
					return true
				})
			}
		}
	}
}
