package p2p
// Peer is an interface that represents the remot node
type Peer interface {
    Close() error

}

//Transport  is anything that handles communicatio between the nodes in the network
//This can be of the form TCP or UDP
type Transport interface{
    ListenAndAccept()error
    Consume () <- chan RPC
}