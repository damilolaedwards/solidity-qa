package utils

import (
	"strings"

	"golang.org/x/net/html"
)

// CoverageReport represents the structured form of the test coverage analysis.
type CoverageReport struct {
	Files []File
}

// File represents a single contract's code coverage details.
type File struct {
	Path         string
	LinesCovered string
	Content      string
}

// FilterCoverageFiles removes contracts from the coverage report that are not in the includePaths or are in the excludePaths.
func FilterCoverageFiles(coverageReport *CoverageReport, includePaths []string, excludePaths []string) {
	indexesToExclude := make([]int, 0)
	for index, file := range coverageReport.Files {
		for _, includePath := range includePaths {
			if !strings.HasPrefix(file.Path, includePath) {
				// Ensure the path is not already in slices to exclude
				if !SliceContains(indexesToExclude, index) {
					indexesToExclude = append(indexesToExclude, index)
				}
				continue
			}
		}
	}

	coverageReport.Files = RemoveElementsFromSlice(coverageReport.Files, indexesToExclude)

	indexesToExclude = make([]int, 0)
	for index, file := range coverageReport.Files {
		for _, excludePath := range excludePaths {
			if strings.HasPrefix(file.Path, excludePath) {
				// Ensure the path is not already in slices to exclude
				if !SliceContains(indexesToExclude, index) {
					indexesToExclude = append(indexesToExclude, index)
				}
				continue
			}
		}
	}

	coverageReport.Files = RemoveElementsFromSlice(coverageReport.Files, indexesToExclude)
}

// ParseCoverageReportHTML parses the HTML content of the medusa fuzz test coverage analysis and parses it into CoverageReport.
func ParseCoverageReportHTML(htmlContent string) (CoverageReport, error) {
	report := CoverageReport{}

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return report, err
	}

	var contracts []File

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "button" && hasClass(n, "collapsible") {
			contract := File{}

			path := extractPath(n)
			contract.Path = path

			nextDiv := findNextElementSibling(n)
			if nextDiv != nil && hasClass(nextDiv, "collapsible-container") {
				linesCovered := findLinesCovered(nextDiv)
				contract.LinesCovered = linesCovered

				content := extractCodeContent(nextDiv)
				contract.Content = content

				contracts = append(contracts, contract)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	report.Files = contracts

	return report, nil
}

// hasClass returns true if the given node has the given class
func hasClass(n *html.Node, class string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" {
			classes := strings.Split(a.Val, " ")
			for _, c := range classes {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

// findNextElementSibling returns the next element sibling of the given node
func findNextElementSibling(n *html.Node) *html.Node {
	for s := n.NextSibling; s != nil; s = s.NextSibling {
		if s.Type == html.ElementNode {
			return s
		}
	}
	return nil
}

// findLinesCovered returns the number of lines covered in the given file in the given node
func findLinesCovered(n *html.Node) string {
	var linesCovered string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && strings.Contains(n.Data, "%)") {
			linesCovered = strings.TrimSpace(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return linesCovered
}

// extractCodeContent returns the code content of the given file in the given node
func extractCodeContent(n *html.Node) string {
	var content strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "pre" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					if hasClass(n, "row-line-covered") {
						content.WriteString("(LINE EXECUTED)")
					} else if hasClass(n, "row-line-uncovered") {
						content.WriteString("(LINE NOT EXECUTED)")
					} else {
						content.WriteString("(LINE NOT EXECUTABLE)")
					}

					content.WriteString(strings.TrimSpace(c.Data))
					content.WriteString("\n")
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return content.String()
}

// extractPath returns the path of the file in the given node
func extractPath(n *html.Node) string {
	var path string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			path = strings.TrimSpace(n.LastChild.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return path
}
