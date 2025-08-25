package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/xfp-881643/gmqtt/config"
	"github.com/xfp-881643/gmqtt/server"
)

import (
	_ "github.com/xfp-881643/gmqtt/plugin/admin"
	_ "github.com/xfp-881643/gmqtt/plugin/auth"
	_ "github.com/xfp-881643/gmqtt/plugin/federation"
	_ "github.com/xfp-881643/gmqtt/plugin/prometheus"
	_ "github.com/xfp-881643/gmqtt/persistence"
	_ "github.com/xfp-881643/gmqtt/topicalias/fifo"
)


func GetListeners(c config.Config) (tcpListeners []net.Listener, websockets []*server.WsServer, err error) {
	for _, v := range c.Listeners {
		var ln net.Listener
		if v.Websocket != nil {
			ws := &server.WsServer{
				Server: &http.Server{Addr: v.Address},
				Path:   v.Websocket.Path,
			}
			if v.TLSOptions != nil {
				ws.KeyFile = v.Key
				ws.CertFile = v.Cert
			}
			websockets = append(websockets, ws)
			continue
		}
		if v.TLSOptions != nil {
			var cert tls.Certificate
			cert, err = tls.LoadX509KeyPair(v.Cert, v.Key)
			if err != nil {
				return
			}
			ln, err = tls.Listen("tcp", v.Address, &tls.Config{
				Certificates: []tls.Certificate{cert},
			})
		} else {
			ln, err = net.Listen("tcp", v.Address)
		}
		tcpListeners = append(tcpListeners, ln)
	}
	return
}


func Test_gmqtt(t *testing.T) {
	hk	:= server.Hooks{
		OnAccept: func(ctx context.Context, conn net.Conn) bool {
			log.Println("[hook] OnAccept: new connection")
			return true
		},
		OnConnected: func(ctx context.Context, client server.Client) {
			log.Printf("[hook] OnConnected: clientID=%s", client.ClientOptions().ClientID)
		},
		OnMsgArrived: func(ctx context.Context, client server.Client, req *server.MsgArrivedRequest) error {
			log.Printf("[hook] OnMsgArrived: clientID=%s topic=%s payload=%s",
				client.ClientOptions().ClientID, req.Message.Topic, string(req.Message.Payload))
			return nil
		},
		OnStop: func(ctx context.Context) {
			log.Println("[hook] OnStop: broker stopped")
		},
	}



	c, _ := config.ParseConfig("/home/xfp/Desktop/mqtt/gmqtt/cmd/gmqttd/default_config.yml")
	tcpListeners, websockets, _ := GetListeners(c)
	l, _ := c.GetLogger(c.Log)


	srv	:= server.New(
		server.WithConfig(c),
		server.WithTCPListener(tcpListeners...),
		server.WithWebsocketServer(websockets...),
		server.WithLogger(l),
	)
	srv.Init(
		server.WithHook(hk),
	)

	// 3) TCP 리스너 실행
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("broker error: %v", err)
		}
	}()

	// 4) 종료 대기
	stop := make(chan struct{})
	<-stop

	// 5) 브로커 종료
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Stop(ctx)
}