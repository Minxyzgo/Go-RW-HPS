package core

import (
	"go-rwhps/rnet"
)

// Protocol All protocol should implement it.
type Protocol interface {
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
