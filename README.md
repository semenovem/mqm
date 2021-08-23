### imbmq provider

Обертка над [mq-golang](https://github.com/ibm-messaging/mq-golang)


#### Описание переменных окружения в файле `env.go`


#### что может:
- поддерживать одновременную работу с любом кол-во очередей
- отправлять сообщения              `mq.Put(ctx, msg)`
- получать очередное сообщение      `mq.Get(ctx, msg)`
- получать сообщение по `CorrelId`  `mq.GetByCorrelId(ctx, msg)`
- получать сообщение по `MsglId`    `mq.GetByMsgId(ctx, msg)`
- просматривать сообщения           `mq.Browse(ctx)`
- подписаться на сообщения          `mq.RegisterEvenInMsg(func)`


# Примеры в sample


#### Создание / подключение
```
// ctx - контекст приложения (context.Context)
var mq = mqpro.New(ctx)

// Установить данные для подключения к менеджерам MQ
// можно использовать стандартные env переменные (смотри в env.go)
mq.UseDefEnv()

TODO - добавить метод для добавления данных подключения


// Подключение
// вернет ошибку только если не указанны данные для подключения
// если в данный момент менеджер ibm mq не доступен - будет повторять попытки подключения
// вернет `nil` после установки соединения с одним из менеджеров в каждой группе (get/put/browse)
// запускать нужно в отдельной горутине, что бы не ждать
//
// если при создании `mqpro.New(ctx)` был передан контекст, который при завершении работы приложения
// будет закрыт cancel() - метод закрытия соединения с IMB MQ `mq.Disconnect()`
// будет вызван автоматически
// в противном случае в `func main` нужно добавить вызов `defer mq.Disconnect()`
err := mq.Connect()

```
