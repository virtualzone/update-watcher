package main

import (
	"strconv"
	"strings"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"regexp"
)

type PackageVersion struct {
	Name		string
	Local		string
	Remote		string
}

func loadURL(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("HTTP response code was %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return bodyBytes, nil
}

func applyRegexp(input string, expression string) (string, error) {
	rx, err := regexp.Compile(expression)
	if (err != nil) {
		return "", fmt.Errorf("Couldn't parse regular expression: '%s'", err)
	}
	match := rx.FindStringSubmatch(input)
	if (len(match) < 2) {
		return "", fmt.Errorf("expression '%s' not found in input string", expression)
	}
	return match[1], nil
}

func getVersion(action CheckAction) string {
	json, err := loadURL(action.URL)
	if err != nil {
		return "ERROR"
	}
	res, err := getDocumentPath(action.Type, json, action.Path)
	if err != nil {
		return "ERROR"
	}
	res = strings.ReplaceAll(res, "\n", " ")
	res, err = applyRegexp(res, action.Regex)
	if err != nil {
		return "ERROR"
	}
	return strings.Trim(res, " ")
}

func getMaxLenName(packages []PackageVersion) int {
	max := 0
	for _, pkg := range packages {
		if (len(pkg.Name) > max) { max = len(pkg.Name) }
	}
	return max
}

func getMaxLenLocal(packages []PackageVersion) int {
	max := 0
	for _, pkg := range packages {
		if (len(pkg.Local) > max) { max = len(pkg.Local) }
	}
	return max
}

func getMaxLenRemote(packages []PackageVersion) int {
	max := 0
	for _, pkg := range packages {
		if (len(pkg.Remote) > max) { max = len(pkg.Remote) }
	}
	return max
}

func printUpdates(packages []PackageVersion) {
	l1, l2, l3 := getMaxLenName(packages), getMaxLenLocal(packages), getMaxLenRemote(packages)
	for _, pkg := range packages {
		if (pkg.Local != pkg.Remote) {
			fmt.Printf("! %"+strconv.Itoa(l1)+"s : %"+strconv.Itoa(l2)+"s -> %"+strconv.Itoa(l3)+"s\n", pkg.Name, pkg.Local, pkg.Remote)
		}
	}
}

func printNoUpdates(packages []PackageVersion) {
	l1, l2 := getMaxLenName(packages), getMaxLenLocal(packages)
	for _, pkg := range packages {
		if (pkg.Local == pkg.Remote) {
			fmt.Printf("✓ %"+strconv.Itoa(l1)+"s : %"+strconv.Itoa(l2)+"s\n", pkg.Name, pkg.Local)
		}
	}
}

func main() {
	fmt.Println("Checking versions, please wait...")
	config := readConfig()
	packages := make([]PackageVersion, len(config.Checks))
	for i, checkItem := range config.Checks {
		var pkg = PackageVersion{
			Name: checkItem.Name,
			Local: getVersion(checkItem.Local),
			Remote: getVersion(checkItem.Remote),
		}
		packages[i] = pkg
	}
	printUpdates(packages)
	printNoUpdates(packages)
}
