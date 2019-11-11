package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func unmarshalRequestBody(r *http.Request, intf interface{}) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	// unmarshal request payload
	if err = json.Unmarshal(bodyBytes, intf); err != nil {
		return err
	}
	return nil
}
