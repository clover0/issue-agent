package functions

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

const FuncGetLatestVersionSearchResult = "get_latest_version_search_result"

const ddgBaseURL = "https://html.duckduckgo.com/html"

func InitFuncGetLatestVersionSearchResultFunction() Function {
	f := Function{
		Name: FuncGetLatestVersionSearchResult,
		Description: `List of summaries of the latest versions of libraries, packages, and other software from the Internet.
The format is 'title:' and 'url:' and 'snippet'. You should get page by the url next`,
		Func:     GetLatestVersionSearchResult,
		FuncType: reflect.TypeOf(GetLatestVersionSearchResult),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"software_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of software for which you want to get the latest version summary list",
				},
			},
			"required":             []string{"software_name"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncGetLatestVersionSearchResult] = f

	return f
}

type GetLatestVersionSearchResultInput struct {
	SoftwareName string
}

func GetLatestVersionSearchResult(input GetLatestVersionSearchResultInput) (_ string, err error) {
	u, err := url.Parse(ddgBaseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	param := url.Values{}
	param.Set("q", fmt.Sprintf("latest version of %s", input.SoftwareName))
	payload := bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest(http.MethodPost, u.String(), payload)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	results := findByClassName(doc, func(className string) bool {
		return className == "results"
	})

	if len(results) == 0 {
		return "", fmt.Errorf("failed to find results")
	}

	return buildTexts(results[0]), nil
}

func buildTexts(node *html.Node) string {
	var texts []string
	results := findByClassName(node, func(className string) bool {
		return className == "result"
	})
	if results == nil {
		return ""
	}

	targets := results
	if len(results) > 5 {
		targets = results[:5]
	}

	for _, result := range targets {
		t := ""

		// get a title
		title := findByClassName(result, func(className string) bool {
			return className == "result__a"
		})
		if title == nil {
			return "error"
		}
		if title[0].FirstChild == nil {
			return "error"
		}
		t += fmt.Sprintf("title:%s\n", title[0].FirstChild.Data)

		// get a URL
		for _, attr := range title[0].Attr {
			if attr.Key == "href" {
				t += fmt.Sprintf("url:%s\n", attr.Val)
				break
			}
		}

		// get a body
		body := findByClassName(result, func(className string) bool {
			return className == "result__snippet"
		})
		if len(body) == 0 {
			return ""
		}

		part := ""
		for e := body[0].FirstChild; e != nil; e = e.NextSibling {
			if e.Type == html.TextNode {
				part += e.Data
			} else {
				// e.g)  b tag
				part += e.FirstChild.Data
			}
		}
		t += fmt.Sprintf("snippet:%s\n", part)

		texts = append(texts, t)
	}

	return strings.Join(texts, "---\n")
}

func findByClassName(node *html.Node, classNameMatcher func(className string) bool) []*html.Node {
	var nodes []*html.Node

	if node.Type == html.ElementNode {
		touchedClass := false
		for _, attr := range node.Attr {
			if attr.Key == "class" {
				classNames := strings.Split(attr.Val, " ")
				for _, name := range classNames {
					if classNameMatcher(name) {
						nodes = append(nodes, node)
					}
				}
				touchedClass = true
			}

			if touchedClass {
				break
			}
		}
	}

	for e := node.FirstChild; e != nil; e = e.NextSibling {
		n := findByClassName(e, classNameMatcher)
		nodes = append(nodes, n...)
	}

	return nodes
}
