package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func usage(programName string) {
	fmt.Println()
	fmt.Printf("usage: %s conf/cf.ini\n", programName)
	fmt.Println()
	fmt.Println("conf/cf.ini      configure file")
	fmt.Println()
}

func main() {
	if len(os.Args) != 2 {
		usage(os.Args[0])
		os.Exit(-1)
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	inifile := os.Args[1]
	conf := &Configure{}
	if err := conf.LoadConfigure(inifile); err != nil {
		log.Println(err)
		os.Exit(-2)
	}
	//log.Printf("%+v\n", conf)
	//创建工作池
	jobQueue := make(JobQueue, conf.WorkPool.MaxQueue)
	d := NewDispatcher(conf.WorkPool.MaxWorker)
	d.Run(jobQueue)

	//取消gorountining
	gonum := 1
	done := make(chan struct{}, gonum)

	//mysql
	cmysql := MysqlGetInstance()
	cmysql.Init(inifile)

	//mongo
	mongo := &MongoConf{}
	mongo.Init(inifile)

	//kafka
	kafkaClient := NewKafkaClient()
	err := kafkaClient.NewConsumer(conf)
	if err != nil {
		os.Exit(-3)
	}
	defer kafkaClient.Close()
	//从kafka的topic里读数据
	kafkaClient.GetTopicMsg(conf.Kafka.Order_log_topic, conf.Kafka.Offset)
	go HandleOrderLog(jobQueue, kafkaClient.Topics[conf.Kafka.Order_log_topic], cmysql, mongo, done)

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)
ConsumerLoop:
	for {
		select {
		case <-signals:
			break ConsumerLoop
		}
	}
	//发送取消信号
	for i := 0; i < gonum; i++ {
		done <- struct{}{}
	}
	d.Stop()
}
