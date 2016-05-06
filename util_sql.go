package main

import (
	"database/sql"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type CheckMysql struct {
	Host   string
	Port   string
	User   string
	Passwd string
	DB     string
	Table  string
}

func MysqlGetInstance() *CheckMysql {
	return &CheckMysql{}
}

func (c *CheckMysql) Init(inifile string) {
	param, _ := ini.Load(inifile)
	c.Host = param.Section("mysql").Key("host").String()
	c.Port = param.Section("mysql").Key("port").String()
	c.User = param.Section("mysql").Key("username").String()
	c.Passwd = param.Section("mysql").Key("passwd").String()
	c.DB = param.Section("mysql").Key("db").String()
	c.Table = param.Section("mysql").Key("table").String()
}

func (c *CheckMysql) Insert(query string) error {
	hostinfo := c.User + ":" + c.Passwd + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.DB + "?charset=utf8"
	db, err := sql.Open("mysql", hostinfo)
	if err != nil {
		//errlog := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
		log.Println("[msyql] connect mysql(", hostinfo, ") error:", err.Error())
		return err
	}
	defer db.Close()

	results, err1 := db.Query(query)
	if err1 != nil {
		//errlog := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
		log.Println("[msyql] execute ", query, " error, ", err1)
		return err1
	}
	defer results.Close()

	for results.Next() {
		var count int
		if err := results.Scan(&count); err != nil {
			log.Fatal(err)
		}
		log.Printf("count=%d\n", count)
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}
