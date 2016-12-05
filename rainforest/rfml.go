package rainforest

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RFTest is a struct representing the Rainforest Test with its settings and steps
type RFTest struct {
	RFMLID      string   `json:"rfml_id"`
	Title       string   `json:"title"`
	StartURI    string   `json:"start_uri"`
	SiteID      int      `json:"site_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	BrowsersMap string   `json:"browser_json"`
	Browsers    []string
	Steps       []interface{}
}

// RFTestStep contains single Rainforest step
type RFTestStep struct {
	Action   string
	Response string
	Redirect bool
}

// RFEmbeddedTest contains an embedded test details
type RFEmbeddedTest struct {
	RFMLID   string
	Redirect bool
}

// RFMLReader reads form RFML formatted file.
// It exports some settings that can be set before parsing.
type RFMLReader struct {
	r *bufio.Reader
	// Version sets the RFML spec version, it's set by NewRFMLReader to the newest one.
	Version int
}

// parseError is a custom error implementing error interface for reporting RFML parsing errors.
type parseError struct {
	line   int
	reason string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("RFML parsing error in line %v: %v", e.line, e.reason)
}

// NewRFMLReader returns RFML parser based on passed io.Reader - typically a RFML file.
func NewRFMLReader(r io.Reader) *RFMLReader {
	return &RFMLReader{
		r:       bufio.NewReader(r),
		Version: 1,
	}
}

// ReadAll parses whole RFML file using RFML version specified by Version parameter of reader
// and returns reulting RFTest
func (r *RFMLReader) ReadAll() (*RFTest, error) {
	parsedRFTest := &RFTest{}
	// Set up a new scanner to read in data line by line
	scanner := bufio.NewScanner(r.r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.HasPrefix(line, "#!") {
			// Handle shebang
			parsedRFTest.RFMLID = line[2:]
		} else if strings.HasPrefix(line, "#") {
			// Handle hashed lines
			content := line[1:]
			if strings.Contains(content, ":") {
				// Handle the key value pair
				split := strings.SplitN(content, ":", 2)
				key := strings.TrimSpace(split[0])
				value := strings.TrimSpace(split[1])
				switch key {
				case "title":
					parsedRFTest.Title = value
				case "start_uri":
					parsedRFTest.StartURI = value
				case "site_id":
					siteID, err := strconv.Atoi(value)
					if err != nil {
						return &RFTest{}, &parseError{lineNum, "Site ID must be a valid integer."}
					}
					parsedRFTest.SiteID = siteID
				case "tags":
					splitTags := strings.Split(value, ",")
					strippedTags := make([]string, len(splitTags))
					for i, tag := range splitTags {
						strippedTags[i] = strings.TrimSpace(tag)
					}
					parsedRFTest.Tags = strippedTags
				case "browsers":
					splitBrowsers := strings.Split(value, ",")
					strippedBrowsers := make([]string, len(splitBrowsers))
					for i, tag := range splitBrowsers {
						strippedBrowsers[i] = strings.TrimSpace(tag)
					}
					parsedRFTest.Browsers = strippedBrowsers
				case "redirect":
					// HANDLE REDIRECT FOR STEP
				default:
					// If it doesn't match known key add it to description
					parsedRFTest.Description += content + "\n"
				}
			} else {
				// If it'a a hashed line without key-value pair add it as a comment
				parsedRFTest.Description += content + "\n"
			}
		} else {
			// Handle non prefixed lines
		}
		fmt.Println(line)
	}
	return parsedRFTest, nil
}
