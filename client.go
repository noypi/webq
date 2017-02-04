package webq

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/noypi/util"
)

type Client struct {
	ctx    context.Context
	PrivK  *util.PrivKey
	client *http.Client
	id     string
}

type ClientOption interface {
	Apply(o *Client)
}

func NewClient(ctx context.Context, opts ...ClientOption) (o *Client, err error) {
	o = new(Client)
	o.ctx = ctx
	for _, opt := range opts {
		opt.Apply(o)
	}
	if nil == o.client {
		var jar *cookiejar.Jar
		if jar, err = cookiejar.New(nil); nil != err {
			return
		}
		o.client = &http.Client{
			Jar: jar,
		}
	}
	if nil == o.PrivK {
		if o.PrivK, err = util.GenPrivKey(2048); nil != err {
			return
		}
	}

	o.id = o.PrivK.PubKey().DigiPrint()
	return
}

func (this *Client) Auth() (err error) {
	var msg AuthMessage
	msg.Content = "hello"
	this.postMsg("/auth", &msg)
	return
}

func (this *Client) Subscribe(topic string) (err error) {
	var msg SubscribeMessage
	this.postMsg("/subscribe", &msg)
	return
}

func (this *Client) Notify() {

}

func (this *Client) post(path string, content []byte) (err error) {
	s := fmt.Sprintf("%s?id=%s", path, this.id)
	buf := bytes.NewBuffer(content)
	_, err = this.client.Post(s, "application/octet-stream", buf)
	return
}

func (this *Client) postMsg(path string, msg interface{}) (err error) {
	bb, err := util.SerializeGob(msg)
	if nil != err {
		return
	}
	bbSignedmsg, err := this.PrivK.SignMessageAndMarshal(bb)
	if nil != err {
		return
	}

	err = this.post("/auth", bbSignedmsg)
	return
}

type withPrivKey struct{ opt *util.PrivKey }

func WithPrivKey(o *util.PrivKey) ClientOption {
	return withPrivKey{opt: o}
}

func (this withPrivKey) Apply(o *Client) {
	o.PrivK = this.opt
}

type withHttpClient struct{ opt *http.Client }

func WithHttpClient(o *http.Client) ClientOption {
	return withHttpClient{opt: o}
}
func (this withHttpClient) Apply(o *Client) {
	o.client = this.opt
}
