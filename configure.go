package main

import (
	//"bufio"
	//"errors"
	"github.com/go-ini/ini"
	//"os"
	"strings"
)

type Kafka struct {
	Hosts           []string
	Order_log_topic string
	Offset          string
}

type Debug struct {
	Debug bool
}

type MySql struct {
	Host   string
	Port   string
	User   string
	Passwd string
	Db     string
	Table  string
}

type WorkPool struct {
	MaxWorker int
	MaxQueue  int
}

type Configure struct {
	Debug
	Kafka
	MySql
	WorkPool
}

func (conf *Configure) LoadConfigure(inifile string) error {

	param, err := ini.Load(inifile)
	if err != nil {
		return err
	}

	enable := param.Section("debug").Key("enable").MustInt()
	debug := true
	if enable == 0 {
		debug = false
	}
	conf.Debug = Debug{Debug: debug}

	//载入kafka配置
	kafka_hosts := param.Section("kafka").Key("hosts").String()
	order_log_topic := param.Section("kafka").Key("order_log_topic").String()
	offset := param.Section("kafka").Key("offset").String()

	sli_khosts := strings.Split(kafka_hosts, ",")
	for index, _ := range sli_khosts {
		sli_khosts[index] = strings.TrimSpace(sli_khosts[index])
	}
	conf.Kafka = Kafka{Hosts: sli_khosts, Order_log_topic: order_log_topic, Offset: offset}
	/*
		mysql_host := param.Section("mysql").Key("host").String()
		mysql_port := param.Section("mysql").Key("port").String()
		mysql_user := param.Section("mysql").Key("username").String()
		mysql_passwd := param.Section("mysql").Key("passwd").String()
		mysql_db := param.Section("mysql").Key("db").String()
		mysql_table := param.Section("mysql").Key("table").String()

		conf.MySql = MySql{Host: mysql_host,
			Port: mysql_port, User: mysql_user, Passwd: mysql_passwd,
			Db: mysql_db, Table: mysql_table}
	*/
	maxWorker := param.Section("work_pool").Key("workers").MustInt()
	maxQueue := param.Section("work_pool").Key("queue").MustInt()
	conf.WorkPool = WorkPool{MaxWorker: maxWorker, MaxQueue: maxQueue}

	return nil
}
