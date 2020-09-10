package messages

import (
	"reflect"
	"testing"
)

func TestUintToBytes(t *testing.T) {
	type args struct {
		value uint
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Zeroes Input",
			args: args{
				value: 0,
			},
			want: []byte{0},
		},
		{
			name: "Non Zeroes Input",
			args: args{
				value: 0x0F0F0F0F0F0F0F0F,
			},
			want: []byte{0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F},
		},
		{
			name: "Non Zeroes Input",
			args: args{
				value: 0x04,
			},
			want: []byte{0x04},
		},
		{
			name: "Non Zeroes Input",
			args: args{
				value: 0x0102030405060708,
			},
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UintToBytes(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UintToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
