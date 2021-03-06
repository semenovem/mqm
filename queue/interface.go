package queue

import (
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

type manager interface {
  RegisterConn() <-chan *ibmmq.MQQueueManager
  IsConfigured() bool
  Reconnect()
}

type provider interface {
  GetBaseCfg() *BaseConfig
}
