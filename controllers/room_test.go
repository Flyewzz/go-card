package controllers

import (
	"net/http"
	"testing"
)

func TestHandlerData_ConnectToRoom(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		hd   *HandlerData
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hd.ConnectToRoom(tt.args.w, tt.args.r)
		})
	}
}
