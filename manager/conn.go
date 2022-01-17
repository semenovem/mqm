package manager

import (
  "context"
  "fmt"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "time"
)

func (m *Mqpro) Connect() error {
  m.log.Trace("Request to establish connection to manager ibmmq...")

  m.mx.Lock()
  defer m.mx.Unlock()

  if !m.IsConfigured() {
    m.log.Warn(ErrNotConfigured)
    return ErrNotConfigured
  }

  if !m.IsDisconn() {
    m.log.Warn(ErrAlreadyConnected)
    return ErrAlreadyConnected
  }

  m.ctx, m.ctxCanc = context.WithCancel(m.rootCtx)

  m.stateConn()

  select {
  case <-m.ctx.Done():
    return nil
  case <-m.registerConn():
  }
  m.log.Trace("Запущено подключение к менеджеру")

  return nil
}

// Вызов из горутины изменения состояния
func (m *Mqpro) connect() error {
  cd := ibmmq.NewMQCD()
  cno := ibmmq.NewMQCNO()
  csp := ibmmq.NewMQCSP()

  cd.ChannelName = m.channel
  cd.ConnectionName = m.endpoint()
  cd.MaxMsgLength = m.maxMsgLen

  // TODO попробовать mutual authentication
  //cd.CertificateLabel

  cno.SecurityParms = csp
  cno.ClientConn = cd
  cno.Options = ibmmq.MQCNO_CLIENT_BINDING
  cno.ApplName = m.app

  if m.tls {
    sco := ibmmq.NewMQSCO()
    sco.KeyRepository = m.keyRepository

    cno.SSLConfig = sco

    cd.SSLCipherSpec = "ANY_TLS12"
    cd.SSLClientAuth = ibmmq.MQSCA_OPTIONAL
  }

  if m.user == "" {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_NONE
  } else {
    csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
    csp.UserId = m.user
    csp.Password = m.pass
  }

  mgr, err := ibmmq.Connx(m.manager, cno)
  if err != nil {
    return err
  }
  m.mgr = &mgr

  m.log.WithFields(map[string]interface{}{
    "endpoint": cd.ConnectionName,
    "manager":  m.manager,
  }).Debug("Connected to ibmmq manager")

  return nil
}

func (m *Mqpro) endpoint() string {
  return fmt.Sprintf("%s(%d)", m.host, m.port)
}

func (m *Mqpro) Disconnect() error {
  m.log.Trace("Request to disconnect from IBM MQ...")

  if m.IsDisconn() {
    m.log.Warn(ErrNoEstablishedConnection)
    return ErrNoEstablishedConnection
  }

  m.ctxCanc()
  m.stateDisconn()

  m.mx.Lock()
  defer m.mx.Unlock()

  select {
  case <-m.rootCtx.Done():
  case <-time.After(m.disconnDelay):
  }

  m.log.Info("Connection dropped")

  return nil
}

// Вызов из горутины изменения состояния
func (m *Mqpro) disconnect() {
  mgr := m.mgr
  if mgr != nil {
    m.mgr = nil
    err := mgr.Disc()
    if err != nil {
      m.log.WithField("mod", "disconnect").Warn(err)
    }
  }
}