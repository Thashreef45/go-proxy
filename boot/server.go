package boot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/Thashreef45/proxy-server/app/handler"
	"github.com/Thashreef45/proxy-server/app/middleware"
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

func (s *Server) Start(cfg model.Config) {
	httpListen := ":" + strconv.Itoa(cfg.ListenPort)
	fmt.Println("Proxy server started running on :", cfg.ListenPort)

	//middleware wrappers
	wrappedHandler := middleware.CORSMiddleware(cfg.CORS, s.handler)

	err := http.ListenAndServe(httpListen, wrappedHandler)
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

	if err = validateConfig(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func validateConfig(c model.Config) error {
	if c.ListenPort <= 0 || c.ListenPort > 65535 {
		return fmt.Errorf("invalid ListenPort: %d", c.ListenPort)
	}
	if c.Backend == "" {
		return fmt.Errorf("backend URL cannot be empty")
	}
	if c.MaxIdleConns < 0 || c.MaxIdleConnsPerHost < 0 || c.IdleConnTimeoutSec < 0 {
		return fmt.Errorf("connection pool values cannot be negative")
	}
	return nil
}
