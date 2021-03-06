package queue

import (
  "context"
)

func (q *Queue) Browse(ctx context.Context) (<-chan *Msg, error) {
  ch, err := q.browse(ctx)
  if err != nil {
    q.errorHandler(err)
  }

  return ch, err
}

func (q *Queue) browse(ctx context.Context) (<-chan *Msg, error) {
  var (
    l = q.log.WithField("method", "BrowseOpen")
    ch   = make(chan *Msg)
    wait = make(chan struct{})
    err  error
  )

  if q.IsClosed() {
    l.Error(ErrNotOpen)
    return nil, ErrNotOpen
  }

  if q.ctlo != nil {
    return nil, ErrBusySubsc
  }

  go func(w chan struct{}) {
    var (
      msg        = &Msg{}
      ll         = l.WithField("method", "BrowseGet")
      oper       = operBrowseFirst
      cx, cancel = context.WithCancel(ctx)
    )
    defer cancel()

    for ctx.Err() == nil {
      err = q.get(cx, oper, msg, ll)
      if err != nil || msg.MsgId == nil {
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
