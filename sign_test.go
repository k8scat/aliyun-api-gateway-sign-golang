package sign

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestSign(t *testing.T) {
	appKey := ""
	appSecret := ""
	api, err := NewAPIGateway(appKey, appSecret)
	if err != nil {
		t.Error(err)
		return
	}

	rawurl := ""
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		t.Error(err)
		return
	}

	cookie := os.Getenv("COOKIE")
	req.Header.Set("Cookie", cookie)

	err = api.Sign(req)
	if err != nil {
		t.Error(err)
		return
	}

	if err == nil {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(string(b))
	}
}
