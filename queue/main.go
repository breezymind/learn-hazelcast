package main

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/config/property"
	"github.com/hazelcast/hazelcast-go-client/core/logger"
	"github.com/sirupsen/logrus"
	"sync"
)

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

	que, _ := cli.GetQueue("test-q")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < lmax; i++ {
			que.Offer(fmt.Sprintf("msg_%d", i))
		}

	}()

	go func() {
		defer wg.Done()

		for i := 0; i < lmax; i++ {
			item, _ := que.Take()
			if item == nil {
				break
			}
			logrus.Info(item)
		}

	}()

	wg.Wait()

	que.Destroy()
	cli.Shutdown()

}
