package p2p
import "net"


//MESSage represents any arbitrary data that is being sent over each transport betwenn two nodes in the network
type RPC struct {
    From    net.Addr
    Payload []byte
}