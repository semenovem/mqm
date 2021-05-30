package main

import (
  "fmt"
  mqpro "github.com/semenovem/gopkg_mqpro"
  "net/http"
)

// Подписка на входящие сообщения
//
// curl localhost:8080/sub
func onRegisterInMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Включено получение сообщений из очереди")
  ibmmq.RegisterEvenInMsg(handlerInMsg)

  fmt.Fprintf(w, "[subcribe] Ok\n")
}

// Отписаться
// curl localhost:8080/unsub
func offRegisterInMsg(w http.ResponseWriter, _ *http.Request) {
  fmt.Println("Отключено получение входящих сообщений")
  ibmmq.UnregisterEvenInMsg()

  fmt.Fprintf(w, "[unsubcribe] Ok\n")
}

// Обработчик входящих сообщений
func handlerInMsg(msg *mqpro.Msg) {
  logMsg(msg)
}