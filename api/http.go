package api

import "encoding/json"

type HttpResp struct {
	Code   int           `json:"code" binding:"required"`
	Msg    string        `json:"msg" binding:"required"`
	Data   interface{}   `json:"data"`
	Total  int           `json:"total,omitempty"` // 通过Data返回分页的结果时 必须填充 Total 表示分页前的总数
	Errors []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Msg   string `json:"msg"`
}

func (f *HttpResp) UnmarshalJSON(data []byte) error {
	type cloneType HttpResp

	rawMsg := json.RawMessage{}
	f.Data = &rawMsg

	if err := json.Unmarshal(data, (*cloneType)(f)); err != nil {
		return err
	}

	return nil
}
