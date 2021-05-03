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
	"github.com/philmacfly/DAPMesh/pkg/dapmesh"
)

var (
	mtx sync.RWMutex
)

type delegate struct {
	gossiper *Gossiper
}

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
	case 'm': // DAPMeshMessages
		var message []*dapmesh.DAPMeshMessage
		if err := json.Unmarshal(b[1:], &message); err != nil {
			return
		}
		mtx.Lock()
		for _, m := range message {
			fmt.Println("Message", m)
		}
		mtx.Unlock()
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return d.gossiper.broadcasts.GetBroadcasts(overhead, limit)
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
	broadcasts *memberlist.TransmitLimitedQueue
}

func StartGossiper(config config.GossiperConf) (*Gossiper, error) {
	var g Gossiper
	hostname, _ := os.Hostname()
	c := memberlist.DefaultWANConfig()
	c.Events = &eventDelegate{}
	c.Delegate = &delegate{gossiper: &g}
	c.BindPort = config.BindingPort
	c.Name = hostname + "-" + uuid.NewUUID().String()
	var err error
	g.memberlist, err = memberlist.Create(c)
	if err != nil {
		return nil, errors.New("Error creating Gossiper:" + err.Error())
	}
	g.memberlist.Join(config.Members)
	g.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return g.memberlist.NumMembers()
		},
		RetransmitMult: c.RetransmitMult,
	}
	node := g.memberlist.LocalNode()
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)
	return &g, nil
}
