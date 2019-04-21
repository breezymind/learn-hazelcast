package main

import (
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

	lb, e := cli.GetList("list-basic")
	if e != nil {
		logrus.Error(e)
		return
	}
	// get map size
	len, _ := lb.Size()

	logrus.Infof("Map %v Size : %v", lb.Name(), len)

	// add
	for i := 0; i < 10; i++ {
		// 	k := fmt.Sprintf("s%d", i)
		// 	logrus.Info(k)
		// 	lb.Put(k, i)
		lb.Add(i)
	}

	len, _ = lb.Size()

	logrus.Infof("Map %v Size : %v", lb.Name(), len)

	// range
	list, _ := lb.ToSlice()
	logrus.Info(list)

	for k, v := range list {
		logrus.Infof("[%v] => %v", k, v)
	}

	lb.AddAt(4, 50)

	logrus.Info(lb.ToSlice())

	lb.RemoveAt(4)

	logrus.Info(lb.ToSlice())
	logrus.Info(lb.Contains(10))

	lb.Destroy()

	// close client
	cli.Shutdown()
}
