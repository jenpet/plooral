package rest

import "encoding/json"

type ResponseBody struct {
	Data         interface{} `json:"data"`
	ErrorMessage *string     `json:"error"`
}

func (r *ResponseBody) SetData(d interface{}) *ResponseBody {
	r.Data = d
	return r
}

func (r *ResponseBody) SetErrorMessage(msg string) *ResponseBody {
	if msg == "" {
		r.ErrorMessage = nil
	} else {
		r.ErrorMessage = &msg
	}
	return r
}

func (r *ResponseBody) JSON() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic("could not marshal JSON response body")
	}
	return b
}