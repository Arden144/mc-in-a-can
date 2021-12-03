package req

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetJson(url string, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("json request exception: %v", r)
		}
	}()

	r, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("http GET failed: %v", r.Status)
	}

	return json.NewDecoder(r.Body).Decode(target)
}
