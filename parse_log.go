package main

import (
	"encoding/json"
	"log"
	"strings"
	"strconv"
)

//all->input->param->flights->adtCabin
type OrderLog struct {
	Client     string
	Code       string
	Cver       string
	Desc       string
	Input      string
	Ip         string
	La         string
	Level      string
	Linkmode   string
	Lo         string
	LogLevel   string
	LogName    string
	LogTime    float64
	Phoneid    string
	Pid        string
	System     string
	Systemtime float64
	Time       string
	Uid        string
	Uuid       string
}

func (e *OrderLog) Unmarshal(item []byte) error {
	err := json.Unmarshal(item, e)
	return err
}

type Input struct {
	Uid        string
	Lo         string
	Client     string
	Param      string
	Userid     string
	Beta       string
	Page       string
	Dver       string
	Iver       string
	Linkmode   string
	Platform   string
	Sid        string
	Imei       string
	Systemtime string
	Pid        string
	Linkcode   string
	Cver       string
	Source     string
	S          string
	P          string
	Pt         string
	La         string
	Uuid       string
	Totalprice float64
	Insureid   string
	Receipts   string
	Contact    string
	Flight     string
	PassInfo   string
}

func (i *Input) Unmarshal(item string) {
	map1 := make(map[string]string, 0)
	sli1 := strings.Split(item, "&")
	for _, v := range sli1 {
		sli2 := strings.SplitN(v, "=", 2)
		if len(sli2) == 2 {
			map1[sli2[0]] = sli2[1]
		}
	}
	for k, v := range map1 {
		switch k {
		case "uid":
			i.Uid = v
		case "lo":
			i.Lo = v
		case "client":
			i.Client = v
		case "param":
			i.Param = v
		case "userid":
			i.Userid = v
		case "beta":
			i.Beta = v
		case "page":
			i.Page = v
		case "dver":
			i.Dver = v
		case "iver":
			i.Iver = v
		case "linkmode":
			i.Linkmode = v
		case "platform":
			i.Platform = v
		case "sid":
			i.Sid = v
		case "imei":
			i.Imei = v
		case "systemtime":
			i.Systemtime = v
		case "pid":
			i.Pid = v
		case "linkcode":
			i.Linkcode = v
		case "cver":
			i.Cver = v
		case "source":
			i.Source = v
		case "s":
			i.S = v
		case "p":
			i.P = v
		case "pt":
			i.Pt = v
		case "la":
			i.La = v
		case "uuid":
			i.Uuid = v
		case "totalprice":
			val , _ := strconv.ParseFloat(v, 64)
			i.Totalprice = val
		case "insureid":
			i.Insureid = v
		case "receipts":
			i.Receipts = v
		case "contact":
			i.Contact = v
		case "flight":
			i.Flight = v
		case "passInfo":
			i.PassInfo = v
		default:
			//errlog := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
			log.Printf("unknown filed, %s=%s\n", k, v)
		}
	}
}

type Param struct {
	Arr             string
	Dep             string
	DepDate         string
	ReturnDate      string
	Flights         []Flight
	SearchCabinCode string
	SearchSrc       string
	SearchType      int
	TripType        string
	Type            string
}

func (p *Param) Unmarshal(item []byte) error {
	err := json.Unmarshal(item, p)
	return err
}

type Flight struct {
	AdtCabin    AdtCabin
	Arr         string
	ArrDate     string
	ArrTerminal string
	Arrtime     string
	Channel     string
	Dep         string
	DepDate     string
	DepTerminal string
	Deptime     string
	LinkNos     []string
	No          string
	Opid        string
	Type        string
}

type AdtCabin struct {
	BaseCode  string
	Code      string
	Fare      float64
	Saleprice float64
	Tax       float64
	Vote      int
}

type Contact struct {
	Name  string
	Phone string
}

func (c *Contact) Unmarshal(item []byte) error {
	err := json.Unmarshal(item, c)
	return err
}
