package protocol

import (
	"fmt"
	"github.com/panjf2000/gnet/logging"
	"go-rwhps/core"
	"go-rwhps/io"
	"go-rwhps/rcmd"
	"go-rwhps/rnet"
	"go-rwhps/server"
	"math/rand"
	"time"
)

type ServerProtocol struct {
	*server.Server
	*core.ProtocolImpl
}

func (p *ServerProtocol) Init(server *server.Server) {
	p.Server = server
}

//for _, player := range server.Players.members {

func (p *ServerProtocol) ServerInfo(player core.Player) *rnet.Packet {
	return rnet.NewPacketData(rnet.PacketServerInfo,
		p.ServerID(),
		p.GameNetVersion(),
		1, //TODO map
		"[p10] test",
		p.Credits,
		p.Mist,
		true,
		1,
		7,
		false,
		player.Admin,
		p.MaxUnit,
		p.MaxUnit,
		p.InitUnit,
		p.Income,
		p.NoNukes,
		false,
		false, // TODO util data
		p.SharedControl,
		false,
		false,
		true,
		false,
	)
}

func (p *ServerProtocol) SendSysChat(msg string) *rnet.Packet {
	return p.SendChat(msg, "SERVER", 5)
}

func (p *ServerProtocol) ReceivePacket(packet rnet.Packet, sender *core.Player) error {
	data := io.NewDataBuffer(packet.Data)
	switch packet.Head {
	case rnet.PacketAddChat:
		str, err := data.ReadUTF()
		if err != nil {
			return err
		}
		isCmd, _ := rcmd.QcCmd.ParseChat(sender, str)
		if isCmd {
			err = p.EachPacket(*p.TeamData(*sender))
			if err != nil {
				return err
			}
		} else {
			err = p.EachPacket(*p.SendChat(str, sender.Name, int32(sender.Team)))
			if err != nil {
				return err
			}
		}
		//p.Packets<-*p.SendChat(str, sender.Name, int32(sender.Team))

	case rnet.PacketPlayerInfo:
		info := &struct {
			_              string
			_, _, _        int32
			Name           string
			B              bool `ignore:"Passwd"`
			Passwd, _, Uid string
			_              int32
			_              string
		}{}
		err := data.ReadData(info)
		if err != nil {
			return err
		}
		sender.Name = info.Name
		sender.Uid = info.Uid

		err = p.EachPacketPlayer(func(player *core.Player) *rnet.Packet {
			return p.TeamData(*player)
		})
		if err != nil {
			return err
		}

		err = p.ServerInfo(*sender).Send(sender)
		if err != nil {
			return err
		}

		err = p.EachPacket(*p.SendSysChat("welcome to server # go-rwhps"))
		if err != nil {
			return err
		}
	case rnet.PacketDisconnect:
		err := p.Exit(sender)
		if err != nil {
			return err
		}
	case rnet.PacketHeartBeatResponse:
		sender.Ping = int(time.Now().UnixMilli()-sender.TimeTemp) >> 1
	}

	return nil
}

func (p *ServerProtocol) RegisterConnection() *rnet.Packet {
	return rnet.NewPacketData(rnet.PacketRegisterConnection,
		p.ServerID(),
		1,
		p.GameNetVersion(),
		p.GameNetVersion(),
		"com.corrodinggames.rts.server",
		p.ConnectUuid.String(),
		100000+rand.Intn(899999),
	)
}

func (p *ServerProtocol) TeamData(player core.Player) *rnet.Packet {
	d := io.NewDataBuffer([]byte{})
	//_ = d.WriteData(true)
	fmt.Println(server.GameServer.Players.Get())
	for _, player := range p.Players.Get() {
		if player == nil {
			_ = d.WriteData(false)
			continue
		} else {
			_ = d.WriteData(
				true,
				0,
			)
			if p.StartGame {
				_ = d.WriteData(
					player.Site,
					player.Ping,
					p.SharedControl,
					player.SharedControl,
				)
				continue
			} else {
				fmt.Println(player)
				_ = d.WriteData(
					player.Site,
					p.Credits,
					player.Team,
					true,
					player.Name,
					false,
					player.Ping,
					time.Now().UnixMilli(),
					false,
					0,
					int32(player.Site),
					byte(0),
					p.SharedControl,
					player.SharedControl,
					false,
					false,
					-9999,
					false,
				)
				if player.Admin {
					_ = d.WriteData(1)
				} else {
					_ = d.WriteData(0)
				}
			}
		}
	}
	buf := d.Bytes()
	fmt.Println(buf)
	d.Reset()
	err := d.WriteData(
		int32(player.Site),
		p.StartGame,
		p.Players.Cap(),
	)
	if err != nil {
		logging.LogErr(err)
		return nil
	}
	//[]byte{0, 0, 0, 0, 1, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 6, 115, 99, 101 ,110, 100, 50, 0, 0, 0, 0, 0, 0, 0, 1 ,125, 29, 42, 66, 233 ,0 ,0 ,0 ,0, 0, 0, 0, 0, 4, 0, 0 ,0 ,0 ,0 ,255, 255, 216, 241, 0 ,0, 0, 0, 0, 0, 0, 0 ,0 ,0})
	//
	//b_d := io.NewDataBuffer([]byte{})
	//b_d.WriteGzipData("teams", buf)
	//head, d_debug, err := b_d.ReadGzipData()
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	//fmt.Println("decode", head, d_debug)

	err = d.WriteGzipData("teams", buf)
	if err != nil {
		logging.LogErr(err)
		return nil
	}
	err = d.WriteData(
		p.Mist,
		p.Credits,
		true,
		1,
		byte(5),
		p.MaxUnit,
		p.MaxUnit,
		p.InitUnit,
		p.Income,
		p.NoNukes,
		false,
		false,
		p.SharedControl,
		false,
	)
	if err != nil {
		logging.LogErr(err)
		return nil
	}
	buf = d.Bytes()
	//logging.Infof("team data:", buf)
	return rnet.NewPacket(rnet.PacketTeamList, buf)
}

func (p *ServerProtocol) ServerID() string {
	return "com.github.go.rwhps"
}

func (p *ServerProtocol) GameNetVersion() int32 {
	return 157
}
