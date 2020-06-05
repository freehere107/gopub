package validator

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/url"
)

var validate = validator.New()

func Validate(data interface{}, model interface{}) error {
	var b []byte
	switch data.(type) {
	case []byte:
		b = data.([]byte)
	case url.Values:
		t := make(map[string]string)
		for k := range data.(url.Values) {
			t[k] = data.(url.Values).Get(k)
		}
		b, _ = json.Marshal(t)
	default:
		b, _ = json.Marshal(data)
	}
	fmt.Println(string(b))
	err := json.Unmarshal(b, model)

	if err != nil {
		return err
	}
	return validate.Struct(model)
}
