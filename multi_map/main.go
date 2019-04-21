package main

import (
	// "fmt"
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
	mb, e := cli.GetMultiMap("mmap-basic")
	if e != nil {
		logrus.Error(e)
		return
	}

	// get map size
	len, _ := mb.Size()

	logrus.Infof("M-Map %v Size : %v", mb.Name(), len)

	// put
	for i := 0; i < 20; i++ {
		mb.Put("twenty", i)
	}

	for i := 0; i < 50; i++ {
		mb.Put("fifty", i)
	}

	len, _ = mb.Size()

	logrus.Infof("M-Map %v Size : %v", mb.Name(), len)

	// range
	keys, _ := mb.KeySet()
	for idx, kname := range keys {
		// get
		val, _ := mb.Get(kname)
		logrus.Infof("[%v] %v => %v", idx, kname, val)
	}

	if ok, _ := mb.ContainsKey("twenty"); !ok {
		logrus.Error("not found key")
	}
	if ok, _ := mb.ContainsValue(60); !ok {
		logrus.Error("not found value")
	}

	mb.Destroy()

	// close client
	cli.Shutdown()

}
