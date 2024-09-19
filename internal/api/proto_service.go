package api

import (
	"net/http"
	"sync"

	"connectrpc.com/grpcreflect"
	messagingv1connect "github.com/mcorrigan89/messaging/internal/api/serviceapis/messaging/v1/messagingv1connect"
	"github.com/mcorrigan89/messaging/internal/config"
	"github.com/mcorrigan89/messaging/internal/services"
	"github.com/rs/zerolog"
)

type ProtoServer struct {
	config            *config.Config
	wg                *sync.WaitGroup
	logger            *zerolog.Logger
	services          *services.Services
	messagingV1Server *MessagingServerV1
}

func NewProtoServer(cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup, services *services.Services) *ProtoServer {

	messagingV1Server := newMessagingProtoUrlServer(cfg, logger, wg, services)

	return &ProtoServer{
		config:            cfg,
		wg:                wg,
		logger:            logger,
		services:          services,
		messagingV1Server: messagingV1Server,
	}
}

func (s *ProtoServer) Handle(r *http.ServeMux) {

	reflector := grpcreflect.NewStaticReflector(
		"serviceapis.messaging.v1.EmailService",
	)

	reflectPath, reflectHandler := grpcreflect.NewHandlerV1(reflector)
	r.Handle(reflectPath, reflectHandler)
	reflectPathAlpha, reflectHandlerAlpha := grpcreflect.NewHandlerV1Alpha(reflector)
	r.Handle(reflectPathAlpha, reflectHandlerAlpha)

	messagingV1Path, messagingV1Handle := messagingv1connect.NewEmailServiceHandler(s.messagingV1Server)
	r.Handle(messagingV1Path, messagingV1Handle)
}
