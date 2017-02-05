package webq

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/noypi/util"
	"github.com/noypi/webutil"
)

const (
	KClient        = "$client"
	KAuthenticated = "$authenticated"
	KBbMessage     = "$bbmsg"
	SessionName    = "noypi/webq"
)

type Server struct {
	*util.PrivKey
	mux    webutil.Mux
	sstore *webutil.SessionStore
}

type ServerOption interface {
	Apply(*Server)
}

func NewServer(ctx context.Context, opts ...ServerOption) (o *Server, err error) {
	o = new(Server)
	for _, opt := range opts {
		opt.Apply(o)
	}
	if nil == o.mux {
		o.mux = webutil.NewPreferredMux()
	}
	o.sstore = webutil.NewCookieSession()
	o.initHandlers()
	return
}

func (this *Server) Subscribe(clientId string, msg *SubscribeMessage) {
	if 0 == len(msg.Topic) {
		return
	}
}

func (this *Server) Publish(clientId string, msg *PublishMessage) {
	if 0 == len(msg.Topic) {
		return
	}
}

func (this *Server) initHandlers() {
	mfn := webutil.MidFn

	mdef := func(ms ...*webutil.MidInfo) []*webutil.MidInfo {
		as := []*webutil.MidInfo{
			mfn(hMustHaveId),
			mfn(this.hMustHaveClient),
			mfn(hVerifyMsgBody),
			mfn(this.sstore.AddSessionHandler, SessionName),
		}
		if 0 < len(ms) {
			as = append(as, ms...)
		}
		return as
	}
	this.mux.Handle("/login/auth",
		webutil.MidSeqFunc(this.hAuthenticate,
			mdef()...,
		),
	)

	this.mux.Handle("/login/authpass",
		webutil.MidSeqFunc(this.hAuthenticate,
			mdef()...,
		),
	)

	this.mux.Handle("/info/client/subscriptions",
		webutil.MidSeqFunc(this.hAuthenticate,
			mdef()...,
		),
	)
	this.mux.Handle("/info/server/pubk",
		webutil.MidSeqFunc(this.hAuthenticate,
			mdef()...,
		),
	)

	this.mux.Handle("/msg/subscribe",
		webutil.MidSeqFunc(this.hSubscribe,
			mdef(mfn(hValidate))...,
		))

	this.mux.Handle("/msg/register",
		webutil.MidSeqFunc(hRegister,
			mdef(mfn(hValidate))...,
		))

}

type ClientInfo struct {
	PubK *util.PubKey
}

func (this *Server) GetClientInfo(id string) (client *ClientInfo, err error) {
	return
}

func (this ClientInfo) VerifyMessage(rdr io.Reader) (msg []byte, err error) {
	bb, err := ioutil.ReadAll(rdr)
	if nil != err {
		return
	}
	msg, err = this.PubK.VerifyMessageRaw(bb)
	return
}
