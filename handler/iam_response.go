package handler

import (
	"encoding/xml"
	"net/http"
	"strconv"
)

type IAMResponser interface {
	Send(http.ResponseWriter)
}

type IAMErrorResponse struct {
	httpStatus int
	XMLName    xml.Name `xml:"Error"`
	Code       string   `xml:"Code"`
	Message    string   `xml:"Message"`
	Resource   string   `xml:"Resource"`
	RequestId  string   `xml:"RequestId"`
}

func NewIAMErrorResponse(status int, code, message, resource, requestId string) IAMResponser {
	return &IAMErrorResponse{
		httpStatus: status,
		Code:       code,
		Message:    message,
		Resource:   resource,
		RequestId:  requestId,
	}
}

func FormatIAMResponse(resp IAMResponser, typ string) ([]byte, error) {
	if typ == "application/xml" {
		return xml.MarshalIndent(resp, "", "  ")
	} else {
		return nil, nil
	}
}

func (resp *IAMErrorResponse) Send(w http.ResponseWriter) {
	body, _ := xml.MarshalIndent(resp, "", " ")
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(resp.httpStatus)
	w.Write(body)
}

type IAMNilResponse struct {
	status int
}

func NewIAMNilResponse(status int) IAMResponser {
	return &IAMResponse{status}
}

func (resp *IAMNilResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(resp.status)
}
