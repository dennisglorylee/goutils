package ws

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
)

// 连接动态参数选项
type ConnOption func(*Connection)

//通用option
func ConnectCbOption(connCallback IConnCallback) ConnOption {
	return func(conn *Connection) {
		conn.connCallback = connCallback
	}
}

func HeartbeatCbOption(heartbeatCallback IHeartbeatCallback) ConnOption {
	return func(conn *Connection) {
		conn.heartbeatCallback = heartbeatCallback
	}
}

func SendBufferOption(bufferSize int) ConnOption {
	return func(conn *Connection) {
		conn.sendBuffer = make(chan *Message, bufferSize)
	}
}

func defaultNetParamsOption() ConnOption {
	return func(conn *Connection) {
		conn.maxFailureRetry = 10                   //重试次数
		conn.readWait = 60 * time.Second            //读等待
		conn.writeWait = 60 * time.Second           //写等待
		conn.temporaryWait = 500 * time.Millisecond //网络抖动重试等待
	}
}

func NetMaxFailureRetryOption(maxFailureRetry int) ConnOption {
	return func(conn *Connection) {
		conn.maxFailureRetry = maxFailureRetry
	}
}

func NetReadWaitOption(readWait time.Duration) ConnOption {
	return func(conn *Connection) {
		conn.readWait = readWait
	}
}

func NetWriteWaitOption(writeWait time.Duration) ConnOption {
	return func(conn *Connection) {
		conn.writeWait = writeWait
	}
}

func NetTemporaryWaitOption(temporaryWait time.Duration) ConnOption {
	return func(conn *Connection) {
		conn.temporaryWait = temporaryWait
	}
}

//服务端特有
//upgrader定制
func SrvUpgraderOption(upgrader *websocket.Upgrader) ConnOption {
	return func(conn *Connection) {
		conn.upgrader = upgrader
	}
}

//为每种消息拉取逻辑分别注册不同的通道
func SrvPullChannelsOption(channels []int) ConnOption {
	return func(conn *Connection) {
		pullChannelMap := make(map[int]chan struct{})
		for _, channel := range channels {
			pullChannelMap[channel] = make(chan struct{}, 2)
		}

		conn.pullChannelMap = pullChannelMap
	}
}

func SrvUpgraderCompressOption(compress bool) ConnOption {
	return func(conn *Connection) {
		conn.upgrader.EnableCompression = compress
	}
}

func SrvCheckOriginOption(checkOrigin func(r *http.Request) bool) ConnOption {
	return func(conn *Connection) {
		conn.upgrader.CheckOrigin = checkOrigin
	}
}

// 客户端专用
// 默认使用时间戳来记录客户端所连服务器的id
func ClientIdOption(id string) ConnOption {
	return func(conn *Connection) {
		conn.id = id
	}
}

func ClientDialOption(dialer *websocket.Dialer) ConnOption {
	return func(conn *Connection) {
		conn.dialer = dialer
	}
}

func ClientDialWssOption(sUrl string, secureWss bool) ConnOption {
	u, err := url.Parse(sUrl)
	if err != nil {
		log.Error(context.Background(), "Parse url %s err:%v", sUrl, err)
	}

	return func(conn *Connection) {
		if u != nil && u.Scheme == "wss" && !secureWss {
			conn.dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
}

func ClientDialCompressOption(compress bool) ConnOption { //, compressLevel int
	return func(conn *Connection) {
		conn.dialer.EnableCompression = compress
	}
}

func ClientDialHandshakeTimeoutOption(handshakeTimeout time.Duration) ConnOption {
	return func(conn *Connection) {
		conn.dialer.HandshakeTimeout = handshakeTimeout
	}
}

func ClientDialRetryOption(retryNum int, retryInterval time.Duration) ConnOption {
	return func(conn *Connection) {
		conn.dialRetryNum = retryNum
		conn.dialRetryInterval = retryInterval
	}
}
