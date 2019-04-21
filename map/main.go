package main

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/config/property"
	"github.com/hazelcast/hazelcast-go-client/core/logger"
	"github.com/sirupsen/logrus"
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
	mb, e := cli.GetMap("map-basic")
	if e != nil {
		logrus.Error(e)
		return
	}
	// get map size
	len, _ := mb.Size()

	logrus.Infof("Map %v Size : %v", mb.Name(), len)

	// put
	for i := 0; i < 50; i++ {
		k := fmt.Sprintf("s%d", i)
		logrus.Info(k)
		mb.Put(k, i)
	}

	len, _ = mb.Size()

	logrus.Infof("Map %v Size : %v", mb.Name(), len)

	// range
	keys, _ := mb.KeySet()
	for idx, kname := range keys {
		// get
		val, _ := mb.Get(kname)
		logrus.Infof("[%v] %v => %v", idx, kname, val)
	}

	mb.Destroy()

	// close client
	cli.Shutdown()
}
