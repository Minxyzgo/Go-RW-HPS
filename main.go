package main

import (
	"encoding/binary"
	"fmt"
	"github.com/panjf2000/gnet"
	"go-rwhps/protocol"
	"go-rwhps/server"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

//
func main() {
	fmt.Println("start server...")
	proto := protocol.ServerProtocol{}
	server.GameServer = server.NewServer(10, &proto)
	proto.Server = server.GameServer
	codec := gnet.NewLengthFieldBasedFrameCodec(gnet.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                -4,
		LengthIncludesLengthFieldLength: false,
	}, gnet.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    4,
		InitialBytesToStrip: 4,
	})
	fmt.Println(proto)
	log.Fatal(gnet.Serve(
		server.GameServer,
		"tcp://:50000",
		gnet.WithMulticore(true),
		gnet.WithCodec(codec),
		gnet.WithTicker(true),
		gnet.WithLogLevel(zapcore.DebugLevel),
		gnet.WithNumEventLoop(3),
		gnet.WithTCPKeepAlive(time.Second*15)))
	//listen, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 50000})
	//if err != nil {
	//	fmt.Println("listen failed, err:", err)
	//	return
	//}
	//go func() {
	//	for {
	//		p := <-server.GameServer.Packets
	//		server.GameServer.Players.Each(func(player *core.Player) {
	//			err := p.Send(player)
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//		})
	//	}
	//}()
	//for {
	//	conn, err := listen.AcceptTCP()//监听是否有连接
	//	if err != nil {
	//		fmt.Println("accept failed, err:", err)
	//		continue
	//	}
	//
	//	go func(player *core.Player){
	//		defer fmt.Println(player.Close())
	//		for {
	//			p, err := rnet.Process(player)
	//			if err != nil {
	//				fmt.Println(err)
	//				_ = conn.Close()
	//				break
	//			}
	//			err = server.GameServer.ReceivePacket(*p, *player)
	//			if err != nil {
	//				fmt.Println(err)
	//				_ = conn.Close()
	//				break
	//			}
	//		}
	//	}(server.GameServer.Join(conn))
	//}
}

//func main() {
//	conn, err := net.Dial("tcp", "127.0.0.1:5123")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	_,_=conn.Write([]byte{0,0,0,45,0,0,0,160,0, 22, 99, 111, 109, 46, 99, 111, 114, 114, 111, 100, 105, 110, 103, 103, 97, 109, 101,115 ,46 ,114 ,116, 115 ,0, 0, 0, 3, 0, 0, 0, 151, 0, 0 ,0 ,2 ,0 ,0, 6, 115, 99, 101, 110, 100, 50})
//	buf:=make([]byte,1024)
//	_,_=conn.Read(buf)
//	fmt.Println(buf)
//
//}
//
