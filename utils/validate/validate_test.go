package validate

import "testing"

func TestCheckPhone(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "12345678901", args: args{phone: "12345678901"}, want: false},
		{name: "133123", args: args{phone: "133123"}, want: false},
		{name: "13345678901", args: args{phone: "13345678901"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPhone(tt.args.phone); got != tt.want {
				t.Errorf("CheckPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckEmail(t *testing.T) {
	r := CheckEmail("123@x.y")
	t.Log(r)
	r = CheckEmail("123.x")
	t.Log(r)
	r = CheckEmail("123@x")
	t.Log(r)
}
