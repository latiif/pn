package pn

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	personnummer "github.com/personnummer/go"
)

type Info struct {
	Firstname string
	Lastname  string
	Address   string
}

func GetInfoWithClient(pn interface{}, client *http.Client) (Info, error) {
	nothing := Info{}
	switch v := pn.(type) {
	case int, int32, int64, uint, uint32, uint64:
		return getInfo(fmt.Sprint(v), client)
	case string:
		return getInfo(v, client)
	default:
		return nothing, invalidInputValue(v)
	}
}

func GetInfo(pn interface{}) (Info, error) {
	return GetInfoWithClient(pn, &http.Client{}) // Use default Go HTTP Client
}

func getInfo(pn string, client *http.Client) (Info, error) {
	nothing := Info{}

	if !personnummer.Valid(pn) {
		return nothing, invalidPn(pn)
	}
	resp, err := client.Get(fmt.Sprintf("https://www.merinfo.se/search?who=%s", pn))
	if err != nil {
		return nothing, cannotPerformRequest()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nothing, cannotPerformRequest()
	}
	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nothing, cannotParseResponse()
	}
	return Info{
		Firstname: htmlquery.InnerText(htmlquery.Find(doc, `/html/body/main/div[1]/div[3]/div[2]/div[3]/div[1]/div/div/div[1]/h2/a/text())`)[0]),
		Lastname:  htmlquery.InnerText(htmlquery.Find(doc, `/html/body/main/div[1]/div[3]/div[2]/div[3]/div[1]/div/div/div[1]/h2/a/span`)[0]),
		Address:   htmlquery.InnerText(htmlquery.Find(doc, `/html/body/main/div[1]/div[3]/div[2]/div[3]/div[1]/div/div/div[2]/p[2]`)[0]),
	}, nil
}

func invalidInputValue(value interface{}) error {
	return fmt.Errorf("Invalid input value: %v", value)
}

func invalidPn(pn string) error {
	return fmt.Errorf("Invalid personal number: %q", pn)
}

func cannotPerformRequest() error {
	return fmt.Errorf("Couldn't perform request")
}

func cannotParseResponse() error {
	return fmt.Errorf("Couldn't parse response")
}
