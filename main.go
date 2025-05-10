package main

import (
	"SuperAgent/global"
	"SuperAgent/initialize"
	"SuperAgent/router"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig(".")
	initialize.InitClient()
	initialize.InitEmbedder()
	initialize.InitIndexer()
	initialize.InitRetriever()
	initialize.InitChatModel()
	iddr := fmt.Sprintf("%s:%s", global.ServerConfig.HOST, global.ServerConfig.PORT)
	h := server.Default(server.WithHostPorts(iddr))
	router.InitRouter(h)
	h.Spin()

}
