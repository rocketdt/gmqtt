package gmqtt

import (
	"net"

	// "go.uber.org/zap"
	"github.com/rocketdt/gmqtt/zap"
)

type Options func(srv *server)

// WithConfig set the config of the server
func WithConfig(config Config) Options {
	return func(srv *server) {
		srv.config = config
	}
}

// WithTCPListener set  tcp listener(s) of the server. Default listen on  :1883.
func WithTCPListener(lns ...net.Listener) Options {
	return func(srv *server) {
		srv.tcpListener = append(srv.tcpListener, lns...)
	}
}

// WithWebsocketServer set  websocket server(s) of the server.
func WithWebsocketServer(ws ...*WsServer) Options {
	return func(srv *server) {
		srv.websocketServer = ws
	}
}

// WithPlugin set plugin(s) of the server.
func WithPlugin(plugin ...Plugable) Options {
	return func(srv *server) {
		srv.plugins = append(srv.plugins, plugin...)
	}
}

// WithHook set hooks of the server. Notice: WithPlugin() will overwrite hooks.
func WithHook(hooks Hooks) Options {
	return func(srv *server) {
		srv.hooks = hooks
	}
}

func WithLogger(logger *zap.Logger) Options {
	return func(srv *server) {
		zaplog = logger
	}
}
