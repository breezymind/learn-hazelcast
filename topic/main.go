package main

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/config/property"
	"github.com/hazelcast/hazelcast-go-client/core"
	"github.com/hazelcast/hazelcast-go-client/core/logger"
	"github.com/sirupsen/logrus"
	"sync"
)

type Listener struct {
	wg *sync.WaitGroup
}

func (t *Listener) OnMessage(msg core.Message) error {
	logrus.Info("Got message: ", msg.MessageObject())
	t.wg.Done()
	return nil
}

func main() {

	// init
	cfg := hazelcast.NewConfig()
	cfg.SetProperty(property.LoggingLevel.Name(), logger.ErrorLevel)

	cfg.NetworkConfig().AddAddress("127.0.0.1:5701")
	cfg.GroupConfig().SetName("hz-compose")
	cfg.GroupConfig().SetPassword("s3crEt")

	// new client
	cli, e := hazelcast.NewClientWithConfig(cfg)
	if e != nil {
		logrus.Error(e)
		return
	}

	lmax := 100000

	topic, _ := cli.GetTopic("test-topic")

	wg := sync.WaitGroup{}
	wg.Add(lmax)

	topic.AddMessageListener(&Listener{
		wg: &wg,
	})
	go func() {
		for i := 0; i < lmax; i++ {
			topic.Publish(fmt.Sprintf("msg_%d", i))
		}
	}()
	wg.Wait()

	topic.Destroy()
	cli.Shutdown()

}
