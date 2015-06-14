package searcher

import (
	"encoding/json"
)


func toPrettyJson(d interface{}) string {
	b,_ := json.MarshalIndent(d, "", "   ")
	s := string(b)
	return s
}




