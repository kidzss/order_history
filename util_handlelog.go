package main

import (
	"log"
	//"os"
	//"errors"
)

func HandleOrderLog(jobQueue JobQueue, tinfo *TopicInfo, c *CheckMysql, m *MongoConf, done chan struct{}) {
	for {
		select {
		case msg := <-tinfo.Msgs:
			//log.Println(string(msg.Value))
			orderLog := &OrderLog{}
			err := orderLog.Unmarshal(msg.Value)
			if err == nil {
				input := &Input{}
				input.Unmarshal(orderLog.Input)
				contact := &Contact{}
				if len(input.Contact) > 0 {
					contact.Unmarshal([]byte(input.Contact))
				}
				param := &Param{}
				var b []byte
				if len(input.Param) > 0 {
					b = []byte(input.Param)
				} else if len(input.Flight) > 0 {
					b = []byte(input.Flight)
				}
				err := param.Unmarshal(b)
				if len(b) > 0 && err == nil {
					//jobQueue <- Job{C: c, Query: "select count(*) from interdetail"}
					r := &Result{}
					r.Set(orderLog, input, param, contact)
					log.Printf("org=%s,dst=%s,date=%s,rdate=%s,status=%s\n", 
										r.Dep, r.Arr, r.Date, r.Rdate, r.Transstatus)
					jobQueue <- Job{C: m, Query: r}
				} else {
					log.Println("unmarshal param failed. err:", err)
				}
			} else {
				log.Println("unmarshal msg.Vaule failed. err:", err)
			}
		case <-done:
			return
		}
	}
}
