package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// LoadConfig loads filed config from config file.
func LoadConfig(fn string, ptr interface{}) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("failed to read config file <%s>, err: %v", fn, err)
	}
	return json.Unmarshal(data, ptr)
}
