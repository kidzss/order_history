package main

import (
	"strconv"
	"time"
)

type Client struct {
	Client string
	Cver   string
	Dver   string
	Iver   string
}

type Result struct {
	Dep             string
	Arr             string
	Date            string
	Rdate           string
	Triptype        string
	Transstatus     string
	Desc            string
	Searchtype      int
	Searchsrc       string
	Searchcabincode string
	Ordertime       string
	Inserttime      string
	Updatetime      int64
	Contact         Contact
	Userid          string
	Uid             string
	Uuid            string
	Phoneid         string
	Totalprice      float64
	Flights         []Flight
	Client          Client
	Loc             []float64
	Ip              string
	Linkmode        string
	Linkcode        string

	System   string
	Pid      string
	Platform string
	Sid      string
	Source   string
	S        string
	Page     string
	P        string
	Pt       string
}

func (r *Result) Set(o *OrderLog, i *Input, p *Param, c *Contact) {
	r.Dep = p.Dep
	r.Arr = p.Arr
	r.Date = p.DepDate
	r.Rdate = p.ReturnDate
	r.Triptype = p.TripType
	r.Transstatus = o.Code
	r.Desc = o.Desc
	r.Searchtype = p.SearchType
	r.Searchsrc = p.SearchSrc
	r.Searchcabincode = p.SearchCabinCode
	stime, _ := strconv.ParseInt(i.Systemtime, 10, 64)
	r.Ordertime = time.Unix(stime/1000, 0).String()
	now := time.Now()
	r.Updatetime = now.UnixNano()/1000/1000
	r.Inserttime = now.String()

	r.Contact = *c
	r.Userid = i.Userid
	r.Uid = i.Uid
	r.Uuid = i.Uuid
	r.Phoneid = o.Phoneid
	r.Totalprice = i.Totalprice
	r.Flights = p.Flights

	r.Client.Client = i.Client
	r.Client.Dver = i.Dver
	r.Client.Iver = i.Iver
	r.Client.Cver = i.Cver

	lg, _ := strconv.ParseFloat(o.Lo, 64)
	la, _ := strconv.ParseFloat(o.La, 64)
	r.Loc = append(r.Loc, lg)
	r.Loc = append(r.Loc, la)
	r.Ip = o.Ip
	r.Linkmode = i.Linkmode
	r.Linkcode = i.Linkcode

	r.System = o.System
	r.Pid = i.Pid
	r.Platform = i.Platform
	r.Sid = i.Sid
	r.Source = i.Source
	r.S = i.S
	r.Page = i.Page
	r.P = i.P
	r.Pt = i.Pt
}
