package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"go-rwhps/core"
	"go-rwhps/io"
	"go-rwhps/rnet"
	"go-rwhps/structs"
	"time"
)

var GameServer *Server

type Server struct {
	Players *structs.Group
	//Packets chan rnet.Packet
	Pool        *goroutine.Pool
	ConnectUuid uuid.UUID
	StartGame   bool
	core.Protocol
	*core.Rules
	*gnet.EventServer
}

func NewServer(maxPlayer int, protocol core.Protocol) *Server {
	server := &Server{
		Players:  &structs.Group{},
		Protocol: protocol,
		Rules:    &core.Rules{},
		//Packets: make(chan rnet.Packet),
		Pool:        goroutine.Default(),
		ConnectUuid: uuid.New(),
	}
	server.Players.Init(maxPlayer)

	return server
}

func (server *Server) React(frame []byte, conn gnet.Conn) (out []byte, action gnet.Action) {

	// Use ants pool to unblock the event-loop.
	fmt.Println(frame)
	_ = server.Pool.Submit(func() {
		p, err := rnet.Process(frame)
		//fmt.Println(p)
		if p.Head == rnet.PacketPreregisterConnection || conn.Context() == nil {
			server.Join(conn, *p)
			return
		}
		if err != nil {
			_ = server.Exit(conn.Context().(*core.Player))
			if err != nil {
				return
			}
		}

		err = server.ReceivePacket(*p, conn.Context().(*core.Player))
		//err = conn.AsyncWrite(p.Data)

		if err != nil {
			_ = server.Exit(conn.Context().(*core.Player))
			if err != nil {
				return
			}
		}
	})

	//var buf []byte
	//data := io.NewDataBuffer(buf)
	//_ = data.WriteData(/*packet*/rnet.FormatByteToPacket("0 0 0 107 0 0 0 161 0 22 99 111 109 46 99 111 114 114 111 100 105 110 103 103 97 109 101 115 46 114 116 115 0 0 0 1 0 0 0 151 0 0 0 151 0 27 99 111 109 46 99 111 114 114 111 100 105 110 103 103 97 109 101 115 46 114 116 115 46 106 97 118 97 0 36 49 48 49 50 100 97 102 55 45 98 53 102 55 45 52 53 55 97 45 97 52 55 53 45 51 50 50 56 101 56 54 101 99 49 49 51 0 2 184 181 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0"))
	//defer logging.Infof("Send Successfully")
	//return buf, 0
	return
}

func (server *Server) OnClosed(conn gnet.Conn, _ error) (action gnet.Action) {
	player := conn.Context()

	if player != nil {
		_ = server.Players.Exit(player.(*core.Player))
	}
	return
}

func (server *Server) Tick() (delay time.Duration, action gnet.Action) {
	_ = server.Pool.Submit(func() {
		err := server.EachPacketPlayer(func(player *core.Player) *rnet.Packet {
			return server.Ping(player)
		})
		if err != nil {
			return
		}
	})
	_ = server.Pool.Submit(func() {
		err := server.EachPacketPlayer(func(player *core.Player) *rnet.Packet {
			return server.TeamData(*player)
		})
		if err != nil {
			return
		}
	})

	return time.Second * 2, gnet.None
}

// Join Make a connection join the server
func (server *Server) Join(conn gnet.Conn, packet rnet.Packet) {
	if server.Players.Full() {
		_ = server.Kick("No free slots on server").Send(conn)
		fmt.Println(conn.Close())
		return
	}

	data := io.NewDataBuffer(packet.Data)
	pData := &struct {
		_            string
		_, _, _      int32
		HasPasswd    bool `ignore:"Passwd"`
		Passwd, Name string
	}{}
	err := data.ReadData(pData)
	if err != nil {
		fmt.Println(err)
		_ = conn.Close()
		return
	}

	if pData.Passwd != server.Passwd {
		err = server.ErrorPasswd().Send(conn)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	player := &core.Player{
		Name: pData.Name,
		Conn: conn,
	}

	err = server.Players.Join(player)
	if err != nil {
		fmt.Println(err)
		_ = conn.Close()
		return
	}

	err = server.RegisterConnection().Send(player)

	if err != nil {
		fmt.Println(server.Exit(player))
		return
	}
	conn.SetContext(player)
	//server.SendServerInfo()
	return
}

// Exit Make a player exit from the group. It will close the connection at the same time
func (server *Server) Exit(player *core.Player) error {
	err := server.Players.Exit(player)
	if err != nil {
		fmt.Println(err)
		return err
	}
	player.SetContext(nil)
	return player.Close()
}

// EachPacket Send the designated package to each player. If you still want to specify a player， please use EachPacketPlayer
func (server *Server) EachPacket(packet rnet.Packet) error {
	var err error
	server.Players.Each(func(player *core.Player) {
		if err != nil {
			return
		}
		err = packet.Send(player)
	})

	return err
}

// EachPacketPlayer Send the designated package to each player, the param p should return the real packet.
func (server *Server) EachPacketPlayer(p func(player *core.Player) *rnet.Packet) error {
	var err error
	server.Players.Each(func(player *core.Player) {
		if err != nil {
			return
		}
		err = p(player).Send(player)
	})

	return err
}
