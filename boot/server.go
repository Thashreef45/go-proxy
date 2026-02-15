package boot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/Thashreef45/proxy-server/app/handler"
	"github.com/Thashreef45/proxy-server/internal/model"
)

type Server struct {
	handler *httputil.ReverseProxy
}

func NewServer(cfg model.Config) (*Server, error) {

	proxyHandler, err := handler.NewProxyHandler(cfg)
	if err != nil {
		return &Server{}, err
	}

	return &Server{
		handler: proxyHandler,
	}, nil
}

func (s *Server) Start(listenPort int) {
	httpListen := ":" + strconv.Itoa(listenPort)
	fmt.Println("Proxy server started running on :", listenPort)
	err := http.ListenAndServe(httpListen, s.handler)
	if err != nil {
		fmt.Println("Error starting server :", err)
	}
}

func InitConfig() (model.Config, error) {

	cfg := model.Config{}

	data, err := os.ReadFile("./config.json")
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal([]byte(data), &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
