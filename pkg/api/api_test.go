package api

import (
	"gonews/pkg/storage"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestAPI_posts(t *testing.T) {
	type fields struct {
		db storage.Interface
		r  *mux.Router
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				db: tt.fields.db,
				r:  tt.fields.r,
			}
			api.posts(tt.args.w, tt.args.r)
		})
	}
}
