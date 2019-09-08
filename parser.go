package main

import (
	"github.com/antchfx/xmlquery"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"encoding/json"
	"fmt"

	"k8s.io/client-go/util/jsonpath"
)

func getDocumentPath(sourceType string, input []byte, selector string) (string, error) {
	switch (strings.ToLower(sourceType)) {
	case "json": return getJSONPath(input, selector)
	case "html": return getHTMLPath(input, selector)
	case "xml": return getXMLPath(input, selector)
	default: return "", fmt.Errorf("Unknown type %s", sourceType)
	}
}

func getJSONPath(jsonBytes []byte, jsonPath string) (string, error) {
	var jsonStruct interface{}
	if err := json.Unmarshal(jsonBytes, &jsonStruct); err != nil {
		return "", err
	}
	parser := jsonpath.New("").AllowMissingKeys(true)
	if err := parser.Parse(jsonPath); err != nil {
		return "", fmt.Errorf("Couldn't parse jsonpath expression: '%s'", err)
	}
	buf := new(bytes.Buffer)
	if err := parser.Execute(buf, jsonStruct); err != nil {
		return "", fmt.Errorf("Couldn't parse jsonpath expression: '%s'", err)
	}
	res := buf.Bytes()
	return string(res[:]), nil
}

func getXMLPath(xmlBytes []byte, xpath string) (string, error) {
	doc, err := xmlquery.Parse(bytes.NewReader(xmlBytes))
	if err != nil {
		return "", err
	}
	item := xmlquery.FindOne(doc, xpath)
	if item == nil {
		return "", err
	}
	return item.InnerText(), nil
}

func getHTMLPath(htmlBytes []byte, selector string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if (err != nil) {
		return "", err
	}
	return doc.Find(selector).First().Text(), nil
}