package rayengine

import (
	"reflect"
	"testing"
)

func Test_readObj(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *aabb
		wantErr bool
	}{
		{"dragon",
			args{"C:\\Users\\lredynk\\gopath\\src\\github.com\\Erdk\\gort\\static\\dragon.obj"},
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readObj(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("readObj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readObj() = %v, want %v", got, tt.want)
			}
		})
	}
}
