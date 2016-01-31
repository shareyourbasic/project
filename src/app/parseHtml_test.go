package main

import (
	"testing"
	"fmt"
	"encoding/json"
	"time"
)

func ATestFetch(t *testing.T) {
	data, _ := GetHtmlData("http://kampanya.kickoffpages.com/?kolid=87ERS")
	js, _ := getJSON(data)
	fmt.Println(js)

	lead := KolLead{}
	json.Unmarshal([]byte(js), &lead)
	fmt.Println("influence direct:", lead)
}

func ATestGetListFetch(t *testing.T) {
	kolids := []string{
		"http://kampanya.kickoffpages.com/?kolid=88S7M",
		"http://kampanya.kickoffpages.com/?kolid=88S8G",
		"http://kampanya.kickoffpages.com/?kolid=88SA5",
		"http://kampanya.kickoffpages.com/?kolid=88SAJ",
		"http://kampanya.kickoffpages.com/?kolid=88SAS",
		"http://kampanya.kickoffpages.com/?kolid=88SB6",
		"http://kampanya.kickoffpages.com/?kolid=88SBF",
		"http://kampanya.kickoffpages.com/?kolid=88SBY",
		"http://kampanya.kickoffpages.com/?kolid=88SC9",
		"http://kampanya.kickoffpages.com/?kolid=88SCJ",
	}
	gd := Data{Emails: make(map[string]*KolLead), TargetTshirt:120, Start:time.Now().Unix()}
	for _, k := range kolids {
		data, _ := GetKolLead(k)
		gd.Emails[data.Email] = data
	}
	b, _ := json.Marshal(gd)
	fmt.Println("OK:")
	fmt.Println(string(b))
}
func TestFloat(t *testing.T) {
	var t1 int
	var t2 int
	t1 = 120
	t2 = 8

	fmt.Println(int((t2 * 100.0) / t1))
}