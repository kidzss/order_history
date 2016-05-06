package main

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"github.com/go-ini/ini"
	"log"
	"strings"
	"time"
)

type MongoConf struct {
	hosts      []string
	user       string
	passwd     string
	database   string
	collection string
}

type MongoClient struct {
	MongoConf
	session *mgo.Session
	db      *mgo.Database
	c       *mgo.Collection
}

func (c *MongoConf) Init(inifile string) {
	param, _ := ini.Load(inifile)
	hosts := param.Section("mongo").Key("hosts").String()
	sli_hosts := strings.Split(hosts, ",")
	c.hosts = sli_hosts
	c.user = param.Section("mongo").Key("user").String()
	c.passwd = param.Section("mongo").Key("passwd").String()
	c.database = param.Section("mongo").Key("database").String()
	c.collection = param.Section("mongo").Key("collection").String()
}

func (c *MongoClient) Execute(result *Result) {
	err := c.Connect()
	if err == nil {
		defer c.Close()
		c.OpenDB()
		c.OpenCollection()
		c.InsertCollection(result)
	}
}

func (c *MongoClient) Connect() error {
	info := &mgo.DialInfo{Addrs: c.hosts,
		Timeout:  60 * time.Second,
		Database: c.database,
		Username: c.user,
		Password: c.passwd,
	}
	ses, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Printf("[mongo] connect %+v fail.", c.hosts, err)
		return err
	}
	c.session = ses
	//log.Printf("[mongo] conect %+v create a session success.\n", c.hosts)
	// Optional. Switch the session to a monotonic behavior.
	c.session.SetMode(mgo.Monotonic, true)
	return err
}

func (c *MongoClient) OpenDB() {
	db := c.session.DB(c.database)
	c.db = db
}

func (c *MongoClient) OpenCollection() {
	c.c = c.db.C(c.collection)
}

func (c *MongoClient) InsertCollection(result *Result) error {
	err := c.c.Insert(result)
	return err
}

func (c *MongoClient) Close() {
	c.session.LogoutAll()
	c.session.Close()
	//log.Printf("[mongo] close mongo session (%+v) success.\n", c.hosts)
}
