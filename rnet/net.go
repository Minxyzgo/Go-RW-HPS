package rnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/panjf2000/gnet"
	"go-rwhps/io"
	"strconv"
	"strings"
)

const PacketServerDebug = 2000

//server packet
const (
	PacketRegisterConnection = 161
	PacketTeamList           = 115
	PacketHeartBeat          = 108
	PacketSendChat           = 141
	PacketServerInfo         = 106
	PacketKick               = 150
	PacketSyncChecksumStatus = 31
	PacketA                  = 30
)

//client packet
const (
	PacketPreregisterConnection = 160
	PacketHeartBeatResponse     = 109
	PacketAddChat               = 140
	PacketPlayerInfo            = 110
	PacketDisconnect            = 111
	PacketAcceptStartGame       = 112
	PacketAcceptButtonGame      = 20
)

//game commands
const (
	PacketAddGameCommand = 20
	PacketTick           = 10
	PacketSync           = 35
	PacketStartGame      = 120
	PacketPasswdError    = 113
)

type Packet struct {
	Head int32
	Data []byte
}

func (packet *Packet) Send(conn gnet.Conn) (err error) {
	var buf []byte
	data := io.NewDataBuffer(buf)
	err = data.WriteData(packet)
	if err != nil {
		return
	}
	buf = data.Bytes()
	fmt.Println(buf)
	err = conn.AsyncWrite(buf)

	//0 0 0 107 0 0 0 161 0 22 99 111 109 46 99 111 114 114 111 100 105 110 103 103 97 109 101 115 46 114 116 115 0 0 0 1 0 0 0 151 0 0 0 151 0 27 99 111 109 46 99 111 114 114 111 100 105 110 103 103 97 109 101 115 46 114 116 115 46 106 97 118 97 0 36 49 48 49 50 100 97 102 55 45 98 53 102 55 45 52 53 55 97 45 97 52 55 53 45 51 50 50 56 101 56 54 101 99 49 49 51 0 2 184 181 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	return
}

func NewPacket(head int32, data []byte) *Packet {
	return &Packet{
		Head: head,
		Data: data,
	}
}

func NewPacketData(head int32, data ...interface{}) *Packet {
	dt := io.NewDataBuffer([]byte{})
	err := dt.WriteData(data...)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &Packet{
		Head: head,
		Data: dt.Bytes(),
	}
}

func FormatByteToBytes(msg string) []byte {
	var data []byte
	msg = strings.Trim(msg, "[]")
	for _, s := range strings.Split(msg, " ") {
		i, _ := strconv.ParseUint(s, 10, 32)
		data = append(data, byte(i))
	}
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, int32(len(data)-4))
	if err != nil {
		return nil
	}
	return append(buf.Bytes(), data...)
}

// Process Decode Packet
func Process(data []byte) (*Packet, error) {
	var head int32
	dt := io.NewDataBuffer(data[:4])
	err := dt.ReadData(&head)
	//i := 0
	//if err != nil {
	//	fmt.Println("read err:", err)
	//	return nil, err
	//}
	//fmt.Println("length: head:", length, head)
	//var data []byte
	//for i < int(length) {
	//	loss := int(length) - i
	//	size, tmp := conn.ReadN(loss)
	//	conn.ShiftN(size)
	//	if err != nil {
	//		fmt.Println("read err:", err)
	//		return nil, err
	//	}
	//
	//	data = append(data, tmp...)
	//	i += size
	//}
	//
	//fmt.Println(data)

	return &Packet{
		Head: head,
		Data: data[4:],
	}, err
}
