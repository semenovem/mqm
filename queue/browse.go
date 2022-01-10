package queue

import (
  "context"
)

func (q *Queue) Browse(ctx context.Context) (<-chan *Msg, error) {
  ch, err := q.browse(ctx)
  if err == ErrConnBroken {
    q.stateError()
  }

  return ch, err
}

func (q *Queue) browse(ctx context.Context) (<-chan *Msg, error) {
  l := q.log.WithField("method", "BrowseOpen")

  q.mxMsg.Lock()
  defer q.mxMsg.Unlock()

  if !q.IsConnected() {
    q.log.Error(ErrNoConnection)
    return nil, ErrNoConnection
  }

  l.Trace("Start open BROWSE")

  var (
    ch   = make(chan *Msg)
    wait = make(chan struct{})
    err  error
    ok   bool
  )

  go func(w chan struct{}) {
    var msg *Msg
    cx, cancel := context.WithCancel(ctx)
    cancel()
    ll := l.WithField("method", "BrowseGet")
    oper := operBrowseFirst

    for ctx.Err() == nil {
      msg, ok, err = q.get(cx, oper, nil, ll)
      if err != nil || !ok {
        break
      }

      if w != nil {
        close(w)
        w = nil
      }
      ch <- msg
      oper = operBrowseNext
    }
    if w != nil {
      close(w)
    }
    close(ch)
    l.Debug("Закрытие канала обзора сообщений BROWSE")
  }(wait)

  select {
  case <-ctx.Done():
  case <-wait:
  }

  if err != nil {
    return nil, err
  }

  l.Debug("Success open for BROWSE")

  return ch, nil
}
