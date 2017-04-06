package linkscanner

import "testing"

func TestFixWWW(t *testing.T) {

	var testlink = "www.example.com"
	var resultlink = "http://www.example.com"

	linklist, errs := HttpScanner(testlink)
	if len(errs) != 0 {
		t.Errorf("FAIL: Unexpected errors parsing WWW %v", errs)
	}

	if len(linklist) > 0 && len(linklist) == 1 {
		if linklist[0] != resultlink {
			t.Errorf("FAIL: WWW not fixed %s became %s", testlink, linklist[0])
		}
	}

	FixWWW(false)

	linklist, errs = HttpScanner(testlink)
	if len(errs) != 0 {
		t.Errorf("FAIL: Unexpected errors parsing WWW %v", errs)
	}

	if len(linklist) > 0 && len(linklist) == 1 {
		if linklist[0] != testlink {
			t.Errorf("FAIL: WWW incorrectly changed %s became %s", testlink, linklist[0])
		}
	}
}
