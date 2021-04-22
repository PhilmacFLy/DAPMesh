package gossiper

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
	"github.com/philmacfly/DAPMesh/pkg/config"
)

var (
	broadcasts *memberlist.TransmitLimitedQueue
	mtx        sync.RWMutex
)

type call struct {
	Version    int
	Message    string
	Subscriber []string
	Groups     []string
	Emergency  bool
}

type delegate struct{}

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}

	switch b[0] {
	case 'd': // data
		var calls []*call
		if err := json.Unmarshal(b[1:], &calls); err != nil {
			return
		}
		mtx.Lock()
		for _, c := range calls {
			fmt.Println("Call", c)
		}
		mtx.Unlock()
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	return nil
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {

}

type eventDelegate struct{}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Println("A node has joined: " + node.String())
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Println("A node has left: " + node.String())
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Println("A node was updated: " + node.String())
}

type Gossiper struct {
	memberlist *memberlist.Memberlist
}

func StartGossiper(config config.Config) error {
	var g Gossiper
	hostname, _ := os.Hostname()
	c := memberlist.DefaultLocalConfig()
	c.Events = &eventDelegate{}
	c.Delegate = &delegate{}
	c.BindPort = config.BindingPort
	c.Name = hostname + "-" + uuid.NewUUID().String()
	var err error
	g.memberlist, err = memberlist.Create(c)
	if err != nil {
		return errors.New("Error creating Gossiper:" + err.Error())
	}
	g.memberlist.Join(config.Members)
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return g.memberlist.NumMembers()
		},
		RetransmitMult: c.RetransmitMult,
	}
	node := g.memberlist.LocalNode()
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)
	return nil
}
