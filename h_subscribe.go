package webq

import (
	"fmt"
	"net/http"

	"github.com/noypi/util"
)

func (this *Server) hSubscribe(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.Write([]byte(err.Error()))
		}
	}()

	ctx := r.Context()
	bbMsg := ctx.Value(KBbMessage).([]byte)

	var msg SubscribeMessage
	if err = util.DeserializeGob(&msg, bbMsg); nil != err {
		util.LogErr(ctx, "%s", err.Error())
		err = fmt.Errorf("bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	this.Subscribe(r.FormValue("id"), &msg)
}
