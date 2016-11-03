package receiver

import (
	"github.com/gaobrian/open-falcon-backend/modules/transfer/receiver/rpc"
	"github.com/gaobrian/open-falcon-backend/modules/transfer/receiver/socket"
)

func Start() {
	go rpc.StartRpc()
	go socket.StartSocket()
}
