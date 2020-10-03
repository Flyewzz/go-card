package controllers

import (
	"net/http"
	"testing"
)

func TestHandleError(t *testing.T) {
	type args struct {
		w    *http.ResponseWriter
		r    *http.Request
		code int
		text string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleError(tt.args.w, tt.args.r, tt.args.code, tt.args.text)
		})
	}
}
