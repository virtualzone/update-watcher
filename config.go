package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Checks		[]CheckItem `json:"checks"`
}

type CheckItem struct {
	Name		string `json:"name"`
	Remote		CheckAction
	Local		CheckAction
}

type CheckAction struct {
	Type		string
	URL			string
	Path		string
	Regex		string
}

func readConfig() Config {
	data, err := ioutil.ReadFile("./config.json")
	if (err != nil) {
		panic("Could not read config.json: " + err.Error())
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(err)
		if e, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("syntax error at byte offset %d\n", e.Offset)
		}
	}
	return config
}