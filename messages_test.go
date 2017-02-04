package webq_test

import (
	"testing"

	"github.com/noypi/util"
	"github.com/noypi/webq"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestAuthMessage(t *testing.T) {
	assert := assertpkg.New(t)
	var msg webq.AuthMessage

	msg.Id = []byte("davidking")
	msg.Signature = []byte("migraso")

	bb, err := util.SerializeGob(&msg)
	assert.Nil(err)

	var msg2 webq.AuthMessage
	err = util.DeserializeGob(&msg2, bb)
	assert.Nil(err)
	assert.Equal(msg, msg2)

}
