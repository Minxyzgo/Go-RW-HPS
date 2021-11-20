package structs

import (
	"errors"
	"go-rwhps/core"
	"strconv"
	"sync"
)

type Group struct {
	members []*core.Player
}

// Join Add a player to the group
func (group *Group) Join(player *core.Player) error {
	if group.Full() {
		return errors.New("out of max length! len:" + strconv.Itoa(group.Len()) + "cap:" + strconv.Itoa(group.Cap()))
	}
	for i := range group.members {
		if group.members[i] == nil {
			group.members[i] = player
			player.Site = byte(i)
			if n := i + 1; n > 0 && n&1 == 0 {
				player.Team = 1
			} else {
				player.Team = 0
			}
			break
		}
	}

	return nil
}

// Exit Remove this player from the group. Note that this will not disconnect.
// A better choice is server.Server.Exit
func (group *Group) Exit(player *core.Player) error {
	var index = -1
	for i, member := range group.members {
		if member == player {
			index = i
			break
		}
	}
	if index == -1 {
		return errors.New("the player isn't in this group")
	}
	//group.members = append(group.members[:index], group.members[index+1:]...)
	group.members[index] = nil
	return nil
}

// Init Initialization the group
func (group *Group) Init(maxPlayer int) {
	group.members = make([]*core.Player, maxPlayer)
}

// Each Iterate every player in the group
func (group *Group) Each(cons func(player *core.Player)) {
	for _, p := range group.Real() {
		cons(p)
	}
}

// Len Group length. This return the length of all the players that are not nil
func (group *Group) Len() int {
	return len(group.Real())
}

// Cap The capacity of this group.
func (group *Group) Cap() int {
	return cap(group.members)
}

// Full Whether this group is full. Len will be real length
func (group *Group) Full() bool {
	return group.Cap() <= group.Len()
}

// Get Original group. It includes nil. If you don't want to do so, please use Real
func (group *Group) Get() []*core.Player {
	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()
	return group.members
}

// Real This returns a group which is not include nil
func (group *Group) Real() []*core.Player {
	tmp := make([]*core.Player, group.Cap())
	x := 0
	for _, member := range group.members {
		if member != nil {
			tmp[x] = member
			x++
		}
	}

	tmp = tmp[:x]
	return tmp
}

func (group *Group) Move(player *core.Player, index int, force bool) bool {
	if force {
		group.Get()[index], group.Get()[player.Site] = group.Get()[player.Site], group.Get()[index]
		group.Get()[player.Site].Site = player.Site
		player.Site = byte(index)
		return true
	} else {
		if group.Get()[index] != nil {
			return false
		} else {
			group.Get()[index] = player
			group.Get()[player.Site] = nil
			player.Site = byte(index)
			return true
		}
	}
}
