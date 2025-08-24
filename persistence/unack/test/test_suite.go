package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xfp-881643/gmqtt/config"
	"github.com/xfp-881643/gmqtt/persistence/unack"
	"github.com/xfp-881643/gmqtt/pkg/packets"
)

var (
	TestServerConfig = config.Config{}
	cid              = "cid"
	TestClientID     = cid
)

func TestSuite(t *testing.T, store unack.Store) {
	a := assert.New(t)
	a.Nil(store.Init(false))
	for i := packets.PacketID(1); i < 10; i++ {
		rs, err := store.Set(i)
		a.Nil(err)
		a.False(rs)
		rs, err = store.Set(i)
		a.Nil(err)
		a.True(rs)
		err = store.Remove(i)
		a.Nil(err)
		rs, err = store.Set(i)
		a.Nil(err)
		a.False(rs)

	}
	a.Nil(store.Init(false))
	for i := packets.PacketID(1); i < 10; i++ {
		rs, err := store.Set(i)
		a.Nil(err)
		a.True(rs)
		err = store.Remove(i)
		a.Nil(err)
		rs, err = store.Set(i)
		a.Nil(err)
		a.False(rs)
	}
	a.Nil(store.Init(true))
	for i := packets.PacketID(1); i < 10; i++ {
		rs, err := store.Set(i)
		a.Nil(err)
		a.False(rs)
	}

}
