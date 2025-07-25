package ws

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"pledge-backend-study/api/models/kucoin"
	"pledge-backend-study/config"
	"pledge-backend-study/log"
	"sync"
	"time"
)

// 初始化manager
var Manager = ServerManager{}

var UserPingPongDurTime = config.Config.Env.WssTimeoutDuration // seconds

const SuccessCode = 0
const PongCode = 1
const ErrorCode = -1

// 匿名嵌入（匿名组合）了 sync.Mutex，所以可以直接调用 s.Lock() 和 s.Unlock()
type Server struct {
	//sync.Mutex 是 Go 标准库中提供的 互斥锁（Mutual Exclusion Lock），用于在并发编程中保护共享资源，防止多个 goroutine 同时访问导致数据竞争
	sync.Mutex
	Id       string
	Socket   *websocket.Conn
	Send     chan []byte
	LastTime int64 // last send time
}
type ServerManager struct {
	// sync.Map 可以零值使用 不需要初始化
	Servers    sync.Map
	Broadcast  chan []byte
	Register   chan *Server
	Unregister chan *Server
}

type Message struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

// 其中 (接收者 接收者类型) 叫做 方法接收者（receiver），是方法和某个类型绑定的关键。
// 就是为了去使用s的一些属性和方法
// s := &Server{}  这是创建一个指向 Server 类型的指针，并初始化为零值（字段都为默认值）。
func (s *Server) SendToClient(data string, code int) {
	s.Lock()
	defer s.Unlock()

	dataBytes, err := json.Marshal(Message{
		Code: code,
		Data: data,
	})
	err = s.Socket.WriteMessage(websocket.TextMessage, dataBytes)
	if err != nil {
		log.Logger.Sugar().Error(s.Id+" SendToClient err ", err)
	}
}

// func (s *Server)：这是一个绑定在 *Server 类型上的方法；
// ReadAndWrite()：方法名，表示这个函数负责处理 WebSocket 的读写；
// 通常这个函数在一个新的 goroutine 中运行，用来处理一个 WebSocket 连接的整个生命周期。
// 这个函数主要完成了以下几件事：
//
// 功能	说明
// 1️⃣ 注册当前连接	把当前 Server 实例加入连接管理器 Manager.Servers
// 2️⃣ 定义清理逻辑	使用 defer 在函数退出时删除连接、关闭 WebSocket、关闭发送通道
// 3️⃣ 启动写协程	从 s.Send 通道读取数据，发送给客户端
// 4️⃣ 启动读协程	从 WebSocket 读取消息，处理心跳（ping/pong）
// 5️⃣ 心跳超时检测	每秒检查一次是否收到心跳，超时则断开连接

// 客户端 <----> 服务器
func (s *Server) ReadAndWrite() {
	//定义错误通道
	//用于在 goroutine 中发生错误时，通知主 goroutine 退出；
	//所有子 goroutine 都可以向这个通道发送错误。
	errChan := make(chan error)

	//将当前连接注册进连接管理器
	Manager.Servers.Store(s.Id, s)

	//defer 表示函数退出时执行这段代码；
	//清理工作包括：
	//从连接池中删除当前连接；
	//关闭 WebSocket 连接；
	//关闭发送通道 s.Send，防止 goroutine 泄漏；
	defer func() {
		Manager.Servers.Delete(s)
		_ = s.Socket.Close()
		close(s.Send)
	}()
	//从 s.Send 通道中接收要发送的消息；
	//如果通道关闭（ok == false），说明连接异常，发送错误；
	//否则调用 s.SendToClient(...) 发送消息给客户端；
	//这是一个典型的 生产者-消费者模型，用于异步发送消息。
	go func() {
		for {
			select {
			case message, ok := <-s.Send:
				if !ok {
					errChan <- errors.New("write message error")
					return
				}
				s.SendToClient(string(message), SuccessCode)
			}
		}
	}()

	//read
	go func() {
		for {
			//从 WebSocket 读取消息；
			//如果出错（比如客户端断开连接），记录日志并通知主 goroutine；
			//如果收到 ping 消息（可能是心跳检测），更新最后收到时间，并回复 pong；
			//这是典型的 心跳响应机制，用于保持连接活跃。
			//第一个返回值：消息类型（text 或 binary），一般不关心，用 _ 忽略；
			//message []byte：读到的原始消息内容；
			//err error：如果有错误（比如客户端断开连接）；
			_, message, err := s.Socket.ReadMessage()
			if err != nil {
				log.Logger.Sugar().Error(s.Id+" ReadMessage err ", err)
				errChan <- err
				return
			}

			//update heartbeat time
			if string(message) == "ping" || string(message) == `"ping"` || string(message) == "'ping'" {
				s.LastTime = time.Now().Unix()
				s.SendToClient("pong", PongCode)
			}
			continue

		}
	}()

	//check heartbeat
	//主 goroutine 一直运行，监听两个事件：
	//每秒钟检查一次是否收到心跳；
	//如果超过 UserPingPongDurTime 没有收到心跳，发送超时消息并断开连接；
	//如果有任意子 goroutine 发生错误，也退出整个连接；
	//这是整个连接的 主控制循环。
	for {
		select {
		case <-time.After(time.Second):
			if time.Now().Unix()-s.LastTime >= UserPingPongDurTime {
				s.SendToClient("heartbeat timeout", ErrorCode)
				return
			}
		case err := <-errChan:
			log.Logger.Sugar().Error(s.Id, " ReadAndWrite returned ", err)
			return
		}
	}
}

// 为 StartServer 的函数，其作用是启动一个服务器，用于监听价格信息并通过 WebSocket 将其发送给所有连接的客户端。
func StartServer() {
	log.Logger.Info("WsServer start")
	//这里是试用for来循环执行,select多路复用来监听通道
	for {
		select {
		case price, ok := <-kucoin.PlgrPriceChan:
			if ok {
				//.Range() 是它的方法，用于 遍历所有键值对：
				//func (m *Map) Range(f func(key, value interface{}) bool)
				//f 是一个函数，对每个键值对都会调用一次；
				//如果返回 true，则继续遍历下一个元素。
				Manager.Servers.Range(func(key, value interface{}) bool {
					value.(*Server).SendToClient(price, SuccessCode)
					return true
				})
			}
		}
	}
}
