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
//	_, _ = conn.Write([]byte{0, 0, 0, 45, 0, 0, 0, 160, 0, 22, 99, 111, 109, 46, 99, 111, 114, 114, 111, 100, 105, 110, 103, 103, 97, 109, 101, 115, 46, 114, 116, 115, 0, 0, 0, 3, 0, 0, 0, 151, 0, 0, 0, 2, 0, 0, 6, 115, 99, 101, 110, 100, 50})
//
//
//	for {
//		time.Sleep(time.Second * 2)
//		buf := make([]byte, 1024)
//		_, _ = conn.Read(buf)
//		fmt.Println(string(buf))
//		p := rnet.FormatByteToBytes("[0 0 0 110 0 22 99 111 109 46 99 111 114 114 111 100 105 110 103 103 97 109 101 115 46 114 116 115 0 0 0 4 0 0 0 151 0 0 0 151 0 6 115 99 101 110 100 50 0 0 27 99 111 109 46 99 111 114 114 111 100 105 110 103 103 97 109 101 115 46 114 116 115 46 106 97 118 97 0 64 51 54 54 52 48 67 66 70 53 69 49 50 54 65 56 49 66 51 52 56 69 56 67 70 56 66 55 52 66 51 65 51 48 53 68 70 52 69 69 70 69 65 67 52 55 70 55 70 68 70 49 66 50 66 54 56 54 52 57 49 67 65 50 70 71 110 161 90 0 109 99 58 55 57 57 53 50 51 109 58 54 57 53 53 56 53 50 53 48 58 56 49 57 50 55 51 54 51 50 49 58 55 57 57 53 50 51 50 58 49 56 48 51 56 54 52 52 48 56 51 58 56 50 55 53 50 51 52 58 45 49 54 53 51 49 55 49 52 52 53 58 57 53 57 53 50 51 54 58 57 56 57 55 49 55 50 51 50 116 49 58 56 49 57 50 55 51 54 51 50 100 58 51 57 57 55 54 49 53]")
//		fmt.Printf("r: %v\n", p)
//		_, _ = conn.Write(p)
//	}
//
//}
//
