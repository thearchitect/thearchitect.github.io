package pty

//import (
//	"context"
//	"github.com/gorilla/websocket"
//	"io"
//)
//
//func (p *PTY) AttachWebsocket(ctx context.Context, conn *websocket.Conn) {
//	p.Attach(ctx, &wsio{conn: conn})
//}
//
//////////////////////////////////////////////////////////////////
////// WSIO
//////
//
//var _ io.ReadWriter = new(wsio)
//
//type wsio struct {
//	conn *websocket.Conn
//}
//
//func (io *wsio) Read(p []byte) (n int, err error) {
//	_, data, err := io.conn.ReadMessage()
//	if err != nil {
//		return 0, err
//	}
//	n = copy(p, data)
//	return
//}
//
//func (io *wsio) Write(p []byte) (n int, err error) {
//	err = io.conn.WriteMessage(websocket.BinaryMessage, p)
//	n = len(p)
//	return
//}
