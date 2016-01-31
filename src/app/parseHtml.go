package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
	"errors"
)



func GetHtmlData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(bytes), nil
}

func getJSON(raw string) (string, error) {
	str := "window.kol_lead ="
	i1 := strings.Index(raw, str)
	if i1 < 0 {
		return "", errors.New("window.kol_lead = not found")
	}
	i2 := strings.Index(raw[i1:], "</script>")
	if i2 < 0 {
		return "", errors.New("</script> not found")
	}
	js := raw[i1 + len(str):i1 + i2]
	if strings.HasSuffix(js, ";") {
		return js[:len(js) - 1], nil
	}
	return js, nil
}

func GetKolLead(redirectURL string) (*KolLead, error) {
	GetHtmlData(redirectURL)
	data, err := GetHtmlData(redirectURL)
	if err != nil {
		return nil, err
	}
	js, err := getJSON(data)
	if err != nil {
		return nil, err
	}
	if len(js) == 0 {
		return nil, errors.New("invalid json data")
	}
	lead := KolLead{}
	err = json.Unmarshal([]byte(js), &lead)
	if err != nil {
		return nil, errors.New("invalid json data")
	}
	return &lead, nil
}


