package linkscanner

import "testing"

func TestFixWWW(t *testing.T) {

	var testlink = "www.example.com"
	var resultlink = "http://www.example.com"

	linklist, errs := HTTPScanner(testlink)
	if len(errs) != 0 {
		t.Errorf("FAIL: Unexpected errors parsing WWW %v", errs)
	}

	if len(linklist) > 0 && len(linklist) == 1 {
		if linklist[0] != resultlink {
			t.Errorf("FAIL: WWW not fixed %s became %s", testlink, linklist[0])
		}
	}

	FixWWW(false)

	linklist, errs = HTTPScanner(testlink)
	if len(errs) != 0 {
		t.Errorf("FAIL: Unexpected errors parsing WWW %v", errs)
	}

	if len(linklist) > 0 && len(linklist) == 1 {
		if linklist[0] != testlink {
			t.Errorf("FAIL: WWW incorrectly changed %s became %s", testlink, linklist[0])
		}
	}
}

func TestIndexOutput(t *testing.T) {
	var testSentence = "this is a short www.example.com sentence."
	var pos = 5

	linklist, errs := HTTPScannerIndex(testSentence)
	if len(errs) != 0 {
		t.Errorf("FAIL: Unexpected errors parsing WWW %v", errs)
	}

	if len(linklist) > 0 && len(linklist) == 1 {
		for x := range linklist {
			for k, _ := range linklist[x] {
				if k != pos {
					t.Errorf("FAIL: Index returned is different than expected %d received, expected %d", k, pos)
				}
			}
		}
	}
}
