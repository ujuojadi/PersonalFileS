package main

import (
     "github.com/ujuojadi/PersonalFileS/p2p"
	 "log"
	 "time"
	 "bytes"
	 "strings"
	 "fmt"
	 
	)

func makeServer(listenAddr string, nodes ...string ) *FileServer {
	tcptransportOpts :=p2p.TCPTransportOpts {
		ListenAddr: listenAddr,
		HandshakeFunc :p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		
	}
	tcpTransport :=p2p.NewTCPTransport(tcptransportOpts)
    

	safeRoot := strings.ReplaceAll(listenAddr, ":", "_")
	fileServerOpts :=ServerOpts{
		StorageRoot : safeRoot + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport: tcpTransport, 
		BootstrapNodes: nodes,

	}
	s:= NewFileServer(fileServerOpts)
    tcpTransport.OnPeer = s.OnPeer

return s
 
	
}


func main (){
	s1 :=makeServer(":3000", "")
	s2 :=makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())

	}()

	time.Sleep( 1  * time.Second)

	go s2.Start()
	time.Sleep(1 * time.Second)
	for i :=0; i < 10; i++ {
		data :=bytes.NewReader([]byte("my big data file here"))
		s2.Store(fmt.Sprintf("myprivatedata_%d", i), data)
		time.Sleep(5* time.Millisecond)
	}

   

	data :=bytes.NewReader([]byte("my big data file"))
	s2.Store("myprivatedata", data)

	// r, err :=s2.Get("myprivatedata")
	// if err!= nil{
	// 	log.Fatal(err)
	// }

	// b, err :=ioutil.ReadAll(r)
	// if err!= nil{
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

    select {}
	
	
}