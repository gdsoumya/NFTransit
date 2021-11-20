package server

import (
	"fmt"
	"time"

	"github.com/gdsoumya/nftransit/relay/pkg/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RelayServer struct {
	Logger     *zap.Logger
	TimeFormat string
	Debug      bool
	StackTrace bool
	ServerPort uint64
	Handlers   handlers.Handler
}

// NewDefaultRelayServer creates an instance of RelayServer with default settings
func NewDefaultRelayServer(logger *zap.Logger, handlers handlers.Handler, debug bool, port uint64) *RelayServer {
	return &RelayServer{
		Debug:      debug,
		Logger:     logger,
		TimeFormat: time.RFC3339,
		StackTrace: true,
		ServerPort: port,
		Handlers:   handlers,
	}
}

// Init starts listening for requests to the relay server
func (ls *RelayServer) Init(corsEnable bool) {
	if !ls.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(Ginzap(ls.Logger, ls.TimeFormat, true), RecoveryWithZap(ls.Logger, ls.StackTrace))

	if corsEnable {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
		}))
	}

	r.POST("/queue_mint", ls.Handlers.QueueMint)
	r.POST("/query_mint", ls.Handlers.QueryMintStatus)
	r.POST("/get_burn", ls.Handlers.GetBurnTx)
	r.POST("/verify_burn", ls.Handlers.VerifyBurnTx)
	r.POST("/user_tokens", ls.Handlers.GetUserTokens)

	ls.Logger.Info("starting server", zap.Uint64("port", ls.ServerPort))
	r.Run(fmt.Sprintf(":%v", ls.ServerPort))
}
