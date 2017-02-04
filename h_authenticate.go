package webq

import (
	"context"
	"fmt"
	"net/http"

	"github.com/noypi/util"
	"github.com/noypi/webutil"
)

func hValidate(nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := webutil.CurrentSession(ctx)

		bValid, _ := session.Values[KAuthenticated].(bool)
		if !bValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized."))
			return
		}

		nexth.ServeHTTP(w, r)
	})
}

func hMustHaveId(nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if 0 == len(id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
			return
		}
		nexth.ServeHTTP(w, r)
	})
}

func (this *Server) hMustHaveClient(nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.FormValue("id")
		client, err := this.GetClientInfo(id)
		if nil == client || nil != err {
			util.LogErr(ctx, "id=%s, not found, err=%v", id, err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			return
		}

		ctx = context.WithValue(ctx, KClient, client)
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}

func hVerifyMsgBody(nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client := ctx.Value(KClient).(*ClientInfo)
		bbMsg, err := client.VerifyMessage(r.Body)
		if nil != err {
			util.LogErr(ctx, "%s", err.Error())
			err = fmt.Errorf("bad request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, KBbMessage, bbMsg)
		nexth.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (this *Server) hAuthenticate(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.Write([]byte(err.Error()))
		}
	}()

	ctx := r.Context()
	if nil == ctx.Value(KBbMessage) {
		util.LogErr(ctx, "%s", err.Error())
		err = fmt.Errorf("bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// authenticated
	session := webutil.CurrentSession(ctx)
	session.Values[KAuthenticated] = true
	if err = session.Save(r, w); nil != err {
		util.LogErr(ctx, "%s", err.Error())
		err = fmt.Errorf("internal error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
