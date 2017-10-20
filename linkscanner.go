package linkscanner

import (
	"bufio"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

// Package linkscanner scans an abitrary piece of text and extracts a URL; currently
// one of HTTP:// HTTPS:// and FTP://. Because a length of text can contain more than
// one URL we return both a list of URLs, and a list of URL parsing errors
// encountered by way of ensuring the link returned to the calling code is
// as valid as possible before further use of it.

var fixProtocol = false

// strings to look for that indicate a web resource
var (
	protoHTTPS = "https://"
	protoHTTP  = "http://"
	protoWww   = "www." // technically not a protocol
	protoFtp   = "ftp://"
	protoMailto = "mailto:"
)

//common line endings that shouldn't be in URL
var common = []string{"�", "\"", "'", ":", ";", ".", "`", ",", "*", ">", ")", "]"}

func cleanLink(link string, www bool) string {
	if www && fixProtocol {
		link = protoHTTP + link
	}

	//utf-8 replacement code character
	//https://codingrigour.wordpress.com/2011/02/17/the-case-of-the-mysterious-characters/
	link = strings.Replace(link, "\xEF\xBF\xBD", "", 1)

	// replace common invalid line-endings
	for _, x := range common {
		if x == link[len(link)-1:] {
			substring := link[0 : len(link)-1]
			return cleanLink(substring, false)
		}
	}
	return strings.TrimSuffix(link, "/")
}

// FixWWW enables the override of the default setting in this package to
// fix wwww links where there isn't a protocol specificed, e.g. http://
func FixWWW(f bool) {
	fixProtocol = f
}

func retrieveLink(literal string) (string, error) {
	literal = strings.ToLower(literal)
	var link string
	if strings.Contains(literal, protoHTTPS) {
		literal = literal[strings.Index(literal, protoHTTPS):]
		link = cleanLink(literal, false)
	} else if strings.Contains(literal, protoHTTP) {
		literal = literal[strings.Index(literal, protoHTTP):]
		link = cleanLink(literal, false)
	} else if strings.Contains(literal, protoFtp) {
		literal = literal[strings.Index(literal, protoFtp):]
		link = cleanLink(literal, false)
	} else if strings.Contains(literal, protoWww) {
		literal = literal[strings.Index(literal, protoWww):]
		link = cleanLink(literal, true)
	} else if strings.Contains(literal, protoMailto) {
		literal = literal[strings.Index(literal, protoMailto):]
		link = cleanLink(literal, false)
	}

	if link != "" {
		_, err := url.Parse(link)
		if err != nil {
			err = errors.Wrapf(err, "Excluding URL after failure to parse: "+link)
			return "", err
		}
	}

	return link, nil
}

// HTTPScanner expects a length of text as input and returns
// two slices dependant on what it discovers. First a unique list of
// URLs parsed successfully by net/url. Second a list of errors
// that were encountered trying to parse the URL found in the text.
func HTTPScanner(content string) ([]string, []error) {

	var hyperlinkList []string
	var errorsList []error

	reader := bufio.NewReader(strings.NewReader(content))
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		link, err := retrieveLink(scanner.Text())
		if err != nil {
			errorsList = append(errorsList, err)
		}
		if err == nil && link != "" {
			seen := false
			for _, x := range hyperlinkList {
				if x == link {
					seen = true
					break
				}
			}
			if !seen {
				hyperlinkList = append(hyperlinkList, link)
			}
		}
	}

	return hyperlinkList, errorsList
}

// HTTPScannerIndex prvides the same basic functionality of HTTPScanner.
// The number of words scanned is monitored. This count becomes an position
// integer providing an approximate index in the text where the hyperlink
// was found. The returned value is not a zero-based index.
func HTTPScannerIndex(content string) ([]map[int]string, []error) {

	var hyperlinkList []map[int]string
	var errorsList []error

	reader := bufio.NewReader(strings.NewReader(content))
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanWords)

	var pos int
	for scanner.Scan() {
		pos++
		link, err := retrieveLink(scanner.Text())
		if err != nil {
			errorsList = append(errorsList, err)
		}
		if err == nil && link != "" {
			tmp := make(map[int]string)
			tmp[pos] = link
			hyperlinkList = append(hyperlinkList, tmp)
		}
	}

	return hyperlinkList, errorsList
}
