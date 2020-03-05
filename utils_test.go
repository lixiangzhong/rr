package rr

import "testing"

func Test_isHidden(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{".git/xxx"}, true},
		{"", args{".git/.xxx"}, true},
		{"", args{"git/.xxx"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHidden(tt.args.path); got != tt.want {
				t.Errorf("isHidden() = %v, want %v", got, tt.want)
			}
		})
	}
}
