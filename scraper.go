package scraper

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

var properties = []map[string]string{
	{"name": "Janka", "unit": "N"},
	{"name": "Average Dried Weight", "unit": "kg/m"},
	{"name": "Modulus of Rupture", "unit": "MPa"},
	{"name": "Elastic Modulus", "unit": "GPa"},
	{"name": "Crushing Strength", "unit": "MPa"},
}

func extractProperty(data string, property map[string]string) []string {

	propertyPos := strings.Index(data, property["name"])
	unitPos := propertyPos + strings.Index(data[propertyPos:], property["unit"])

	result := strings.TrimSpace(data[propertyPos:unitPos])
	formatedText := removeHTMLTags(result)
	return formatedText
}

func getProperties(data string, name string) map[string]interface{} {

	output := make(map[string]interface{})
	output["name"] = name
	var props []string
	for _, prop := range properties {
		props = append(props, extractProperty(data, prop)[0])
	}
	output["properties"] = props
	return output
}

func removeHTMLTags(input string) []string {
	// Regular expression to match HTML tags
	htmlTagRegex := regexp.MustCompile("<[^>]*>")

	// var output []string
	// Remove HTML tags from the input string
	result := htmlTagRegex.ReplaceAllString(input, "")
	formatedString := strings.Replace(result, "&nbsp;", " ", -1)
	splitedText := strings.Split(formatedString, "\n")

	return splitedText
}

func LinkToWood(link string) map[string]interface{} {

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return getProperties(string(body), strings.Split(link, "/")[len(strings.Split(link, "/"))-2])
}
