package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

const FuncGetWebPageFromURL = "get_web_page_from_url"

func InitFuncGetWebPageFromURLFunction() Function {
	f := Function{
		Name:        FuncGetWebPageFromURL,
		Description: "Get the web page from the URL",
		Func:        GetWebPageFromURL,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"url": map[string]interface{}{
					"type":        "string",
					"description": fmt.Sprintf("The URL to get the Web page. More than %d characters are cut off", maxTextLength),
				},
			},
			"required":             []string{"url"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncGetWebPageFromURL] = f

	return f
}

type GetWebPageFromURLInput struct {
	URL string
}

func GetWebPageFromURL(input GetWebPageFromURLInput) (string, error) {
	u, err := url.Parse(input.URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

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

	t := extractText(doc)

	return t, nil
}

const maxTextLength = 80000

func extractText(node *html.Node) string {
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data)
	}

	// ignore script and style elements
	if node.Type == html.ElementNode &&
		(node.Data == "script" || node.Data == "style" || node.Data == "iframe") {
		return ""
	}

	var text strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		childText := extractText(child)
		if childText != "" {
			text.WriteString(childText + " ")
			if text.Len() > maxTextLength {
				return ""
			}
		}
	}

	return strings.TrimSpace(text.String())
}
