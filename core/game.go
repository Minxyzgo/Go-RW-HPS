package core

import (
	"github.com/panjf2000/gnet"
)

type Player struct {
	gnet.Conn
	Site                 byte
	Team                 int
	Admin, SharedControl bool
	Name, Uid            string
	Ping                 int
	TimeTemp             int64
}

type Rules struct {
	Credits, Mist, MaxUnit, InitUnit int32

	Income float32

	NoNukes, SharedControl bool

	Passwd string
}
