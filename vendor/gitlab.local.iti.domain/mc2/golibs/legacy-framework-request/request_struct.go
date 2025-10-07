package request

import (
	"encoding/json"

	comcontext "github.com/mc2soft/framework/communication/context"
)

type requestStruct struct {
	Headers        comcontext.Headers `json:"headers"`
	From           string             `json:"from"`
	To             string             `json:"to"`
	Method         string             `json:"method"`
	Path           string             `json:"path"`
	Data           json.RawMessage    `json:"data"`
	IsAsynchronous bool               `json:"is_asynchronous"`
}
