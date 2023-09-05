package servers

import (
	"api-service/config"
	"api-service/routers"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type GinServer struct {
	server *http.Server
}

func NewGinServer() *GinServer {
	return &GinServer{}
}

func (s *GinServer) Serve(group *errgroup.Group) error {
	if group == nil {
		return nil
	}

	group.Go(func() error {
		svc := gin.New()
		routers.InitRoute(svc)

		s.server = &http.Server{
			Addr:    config.GinAddr,
			Handler: svc,
		}
		if err := s.server.ListenAndServe(); err != nil {
			return err
		}
		return nil
	})

	return nil
}
