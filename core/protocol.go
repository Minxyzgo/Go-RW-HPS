package core

import (
	"fmt"
	"go-rwhps/io"
	"go-rwhps/rnet"
	"time"
)

// Protocol All protocol should implement it.
type (
	Protocol interface {
		ReceivePacket(packet rnet.Packet, sender *Player) error

		ServerInfo(player Player) *rnet.Packet

		AddChat(packet rnet.Packet) string

		SendSysChat(msg string) *rnet.Packet

		SendChat(msg, sender string, team int32) *rnet.Packet

		Kick(msg string) *rnet.Packet

		RegisterConnection() *rnet.Packet

		ErrorPasswd() *rnet.Packet

		TeamData(player Player) *rnet.Packet

		Ping(sender *Player) *rnet.Packet

		ServerID() string

		GameNetVersion() int32
	}

	ProtocolImpl struct{}
)

func (p *ProtocolImpl) ReceivePacket(packet rnet.Packet, sender *Player) error {
	return nil
}

func (p *ProtocolImpl) ServerInfo(player Player) *rnet.Packet {
	return nil
}

func (p *ProtocolImpl) AddChat(packet rnet.Packet) string {
	str, err := io.NewDataBuffer(packet.Data).ReadUTF()
	if err != nil {
		fmt.Println(err)
		return "read error"
	}

	return str
}

func (p *ProtocolImpl) SendSysChat(msg string) *rnet.Packet {
	return nil
}

func (p *ProtocolImpl) SendChat(msg, sender string, team int32) *rnet.Packet {
	return rnet.NewPacketData(rnet.PacketSendChat,
		msg,
		byte(3),
		true,
		sender,
		team,
		team,
	)
}

func (p *ProtocolImpl) Kick(msg string) *rnet.Packet {
	return rnet.NewPacketData(rnet.PacketKick, msg)
}

func (p *ProtocolImpl) RegisterConnection() *rnet.Packet {
	return nil
}

func (p *ProtocolImpl) ErrorPasswd() *rnet.Packet {
	return rnet.NewPacketData(rnet.PacketPasswdError, 0)
}

func (p *ProtocolImpl) TeamData(player Player) *rnet.Packet {
	return nil
}

func (p *ProtocolImpl) Ping(sender *Player) *rnet.Packet {
	sender.TimeTemp = time.Now().UnixMilli()
	return rnet.NewPacketData(rnet.PacketHeartBeat,
		int64(1000),
		byte(0),
	)
}

func (p *ProtocolImpl) ServerID() string {
	return ""
}

func (p *ProtocolImpl) GameNetVersion() int32 {
	return 0
}
