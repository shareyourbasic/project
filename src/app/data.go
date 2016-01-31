package main

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const (
	MaxUnverifiedInfluence = 60
	MaxDirectInfluence = 50
)

type Influence struct {
	Indirect int `json:"indirect"`
	Direct   int `json:"direct"`
}

type KolLead struct {
	ID                 int `json:"id"`
	Email              string `json:"email"`
	SocialID           string `json:"social_id"`
	SocialUrl          string `json:"social_url"`
	RedirectUrl        string `json:"redirect_url"`
	UnverifiedLeads    *Influence `json:"unverified_leads"`
	VerifiedLeads      *Influence `json:"verified_leads"`
	Influence          *Influence `json:"influence"`
	SubscriptionNumber int `json:"subscription_number"`
	Uid                string `json:"uid"`
	Rank               int `json:"rank"`
	Counter            int `json:"counter"`
	LeadCount          int `json:"lead_count"`
	ListId             int `json:"list_id"`
}

func (k *KolLead)IsDone() bool {
	return k.UnverifiedLeads.Direct > MaxUnverifiedInfluence || k.Influence.Direct >= MaxDirectInfluence
}

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return s[i] < s[j]
}

type Data struct {
	VerifiedEmail           int
	TotalEmailRegistered    int
	IndirectEmailRegistered int
	TshirtGot               int
	Start                   int64
	TargetTshirt            int
	Emails                  map[string]*KolLead
}

func LoadData() (*Data, error) {
	file, e := ioutil.ReadFile("data.json")
	if e != nil {
		return nil, e
	}
	var data Data
	err := json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GetTshirtNumber(directInfluence int) int {
	if (directInfluence >= 50) {
		return 12
	}else if (directInfluence >= 25) {
		return 4
	}else if (directInfluence >= 10) {
		return 1
	}
	return 0
}

func (d *Data)GetTemplate() map[string]interface{} {
	t := make(map[string]interface{})
	tn := 0
	link := ""

	ems := make([]string, 0)
	for k, _ := range d.Emails {
		ems = append(ems, k)
	}
	sort.Sort(ByLength(ems))
	td := 0
	ti := 0
	for _, kkk := range ems {
		v := d.Emails[kkk]
		if v.Influence != nil && v.UnverifiedLeads != nil {
			tn += GetTshirtNumber(v.Influence.Direct)
			td += v.Influence.Direct
			ti += v.UnverifiedLeads.Direct
			if !v.IsDone() && len(link) == 0 {
				link = v.SocialUrl
			}
		}
	}

	if len(link) == 0 {
		fmt.Println("INFO: UnverifiedLeads.Direct is full looking for verified users")
		for _, kkk := range ems {
			v := d.Emails[kkk]
			if v.Influence != nil {
				if v.Influence.Direct < MaxDirectInfluence {
					link = v.SocialUrl
					break
				}
			}
		}
	}

	t["TotalDirectInfluence"] = td
	t["TotalUnverified"] = ti
	t["TargetTshirt"] = d.TargetTshirt
	t["CurrentTshirt"] = tn
	t["CurrentLink"] = link

	if tn >= d.TargetTshirt {
		t["DisableClass"] = "disabled"
	}else {
		t["DisableClass"] = ""
	}
	t["ProgressValue"] = int((tn * 100) / d.TargetTshirt)

	return t
}

func (d *Data)UpdateDB(all bool) {
	test := true
	ems := make([]string, 0)
	for k, _ := range d.Emails {
		ems = append(ems, k)
	}
	sort.Sort(ByLength(ems))

	for _, kkk := range ems {
		v := d.Emails[kkk]
		if !test {
			break
		}

		if !v.IsDone() || all {
			fmt.Printf("updating %s to %d un:%d \n", v.Email, v.Influence.Direct, v.UnverifiedLeads.Direct)
			kl, e := GetKolLead(v.RedirectUrl)
			if e == nil {
				test = all || kl.IsDone()
				fmt.Printf("updated %s to %d un:%d \n", v.Email, kl.Influence.Direct, kl.UnverifiedLeads.Direct)
				d.Emails[kkk] = kl
			}
		}
	}
	b, err := json.Marshal(d)
	if err != nil {
		fmt.Println("error: json marshall error", err)
		return
	}
	err = ioutil.WriteFile("data.json", b, os.ModePerm)
	if err == nil {
		fmt.Println("data.json updated")
	}else {
		fmt.Println("error: while writing updated data.json", err)
	}
}

func (d *Data)Calculate() {
	td := 0
	ti := 0
	tn := 0
	ii := 0
	for _, v := range d.Emails {
		tn += GetTshirtNumber(v.Influence.Direct)
		td += v.Influence.Direct
		ti += v.UnverifiedLeads.Direct
		ii += v.UnverifiedLeads.Indirect
	}
	d.VerifiedEmail = td
	d.TotalEmailRegistered = ti
	d.TshirtGot = tn
	d.IndirectEmailRegistered = ii
}