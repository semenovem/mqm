package main

import (
  "context"
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
  "time"
)

// Включение зеркального ответа на входящие сообщения
// curl host:port/on-mirror
func onMirror(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Включена отправка ответов на входящие сообщения")
  cfg.Mirror = true
  subscr()

  printCfg()
  _, _ = fmt.Fprintf(w, "[on-mirror] Ok\n")
}

// Выключение зеркального ответа на входящие сообщения
// curl host:port/off-mirror
func offMirror(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отключена отправка ответов на входящие сообщения")
  cfg.Mirror = false
  unsubscr()

  printCfg()
  _, _ = fmt.Fprintf(w, "[off-mirror] Ok\n")
}


// Отправляет зеркальный ответ
func mirror(msg *mqpro.Msg) {
  fmt.Println()
  fmt.Println("Отправляем ответ: ")
  reply := &mqpro.Msg{
    CorrelId: msg.MsgId,
    Payload:  msg.Payload,
    Props:    msg.Props,
  }

  ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
  defer cancel()

  id, err := ibmmq.Put(ctx, reply)
  if err != nil {
    fmt.Println(">>>>> [ERROR]: Ошибка при отправке ответа")
  } else {
    reply.MsgId = id
    logMsgOut(msg)
  }
}
