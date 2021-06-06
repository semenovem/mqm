package mqpro

import (
  "context"
  "errors"
  "github.com/sirupsen/logrus"
  "sync"
  "time"
)

type Mqpro struct {
  rootCtx               context.Context
  ctx                   context.Context
  ctxCancel             context.CancelFunc
  conns                 []*Mqconn
  connGet               []*Mqconn
  connPut               []*Mqconn
  connBrowse            []*Mqconn
  fnEventInMsg          func(*Msg)    // Обработчик входящих сообщений
  mx                    sync.Mutex    // подключение / отключение
  delayBeforeDisconnect time.Duration // Задержка перед разрывом соединения
  reconnDelay           time.Duration // Задержка при повторных попытках подключения к MQ
  log                   *logrus.Entry
}

const (
  defDisconnDelay = time.Millisecond * 500 // По умолчанию задержка перед разрывом соединения
  defReconnDelay  = time.Second * 3        // По умолчанию задержка повторных попыткок соединения
)

var (
  ErrNoEstablishedConnection = errors.New("ibm mq: no established connections")
  ErrNoConnection            = errors.New("ibm mq: no connections")
  ErrNoData                  = errors.New("ibm mq: no data to connect to IBM MQ")
  ErrConnBroken              = errors.New("ibm mq conn: connection broken")
  ErrPutMsg                  = errors.New("ibm mq: failed to put message")
  ErrGetMsg                  = errors.New("ibm mq: failed to get message")
  ErrBrowseMsg               = errors.New("ibm mq: failed to browse message")
)

func New(rootCtx context.Context) *Mqpro {
  l := logrus.New()
  l.SetLevel(logrus.TraceLevel)

  return &Mqpro{
    rootCtx:               rootCtx,
    delayBeforeDisconnect: defDisconnDelay,
    reconnDelay:           defReconnDelay,
    log:                   logrus.NewEntry(l).WithField("pkg", "mqpro"),
  }
}

func (p *Mqpro) SetConn(connLi ...*Mqconn) {
  for _, conn := range connLi {

    switch conn.Type() {
    case TypeGet:
      p.connGet = append(p.connGet, conn)
    case TypePut:
      p.connPut = append(p.connPut, conn)
    case TypeBrowse:
      p.connBrowse = append(p.connBrowse, conn)

    default:
      p.log.Panic("Unknown connection type")
    }

    p.conns = append(p.conns, conn)

    if p.fnEventInMsg != nil {
      conn.RegisterEventInMsg(p.fnEventInMsg)
    }
  }
}

func (p *Mqpro) SetLogger(l *logrus.Entry) {
  p.log = l
}
