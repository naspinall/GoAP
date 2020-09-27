package coding

import (
	"reflect"
	"testing"
)

func TestEncodeUint16(t *testing.T) {
	type args struct {
		value uint16
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Encoding 2",
			args: args{
				value: 0x0002,
			},
			want: []byte{0x00, 0x02},
		},
		{
			name: "Encoding 0x2002",
			args: args{
				value: 0x2002,
			},
			want: []byte{0x20, 0x02},
		},
		{
			name: "Encoding 0x0202",
			args: args{
				value: 0x0202,
			},
			want: []byte{0x02, 0x02},
		},
		{
			name: "Encoding 0x2000",
			args: args{
				value: 0x2000,
			},
			want: []byte{0x20, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeUint16(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeUint32(t *testing.T) {
	type args struct {
		value uint32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Encoding 2",
			args: args{
				value: 0x0002,
			},
			want: []byte{0x02, 0x00, 0x00, 0x00},
		},
		{
			name: "Encoding 0x0202",
			args: args{
				value: 0x0202,
			},
			want: []byte{0x02, 0x02, 0x00, 0x00},
		},
		{
			name: "Encoding 0x22334466",
			args: args{
				value: 0x22334466,
			},
			want: []byte{0x66, 0x44, 0x33, 0x22},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeUint32(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeUint64(t *testing.T) {
	type args struct {
		value uint64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Encoding 2",
			args: args{
				value: 0x0002,
			},
			want: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "Encoding 0x0202",
			args: args{
				value: 0x0202,
			},
			want: []byte{0x02, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "Encoding 0x22334466",
			args: args{
				value: 0x22334466,
			},
			want: []byte{0x66, 0x44, 0x33, 0x22, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "Encoding 0x22334466",
			args: args{
				value: 0x2233446677889911,
			},
			want: []byte{0x11, 0x99, 0x88, 0x77, 0x66, 0x44, 0x33, 0x22},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeUint64(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeUint(t *testing.T) {
	type args struct {
		value uint
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Encoding 0x02",
			args: args{
				value: 0x01,
			},
			want: []byte{0x01},
		},
		{
			name: "Encoding 0x0201",
			args: args{
				value: 0x0201,
			},
			want: []byte{0x01, 0x02},
		},
		{
			name: "Encoding 0x030201",
			args: args{
				value: 0x030201,
			},
			want: []byte{0x01, 0x02, 0x03},
		},
		{
			name: "Encoding 0x04030201",
			args: args{
				value: 0x030201,
			},
			want: []byte{0x01, 0x02, 0x03},
		},
		{
			name: "Encoding 0x04030201",
			args: args{
				value: 0x0504030201,
			},
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		},
		{
			name: "Encoding 0x04030201",
			args: args{
				value: 0x060504030201,
			},
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		},
		{
			name: "Encoding 0x04030201",
			args: args{
				value: 0x07060504030201,
			},
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		},
		{
			name: "Encoding 0x04030201",
			args: args{
				value: 0x0807060504030201,
			},
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeUint(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeUint16(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "2",
			args: args{
				input: []byte{0x02},
			},
			want: 2,
		},
		{
			name: "0x0202",
			args: args{
				input: []byte{0x02, 0x02},
			},
			want: 0x0202,
		},
		{
			name: "0x0220",
			args: args{
				input: []byte{0x20, 0x02},
			},
			want: 0x0220,
		},
		{
			name: "0x3210",
			args: args{
				input: []byte{0x10, 0x32},
			},
			want: 0x3210,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeUint16(tt.args.input); got != tt.want {
				t.Errorf("DecodeUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeUint32(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "0x3210",
			args: args{
				input: []byte{0x10, 0x32},
			},
			want: 0x3210,
		},
		{
			name: "0x76543210",
			args: args{
				input: []byte{0x10, 0x32, 0x54, 0x76},
			},
			want: 0x76543210,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeUint32(tt.args.input); got != tt.want {
				t.Errorf("DecodeUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeUint64(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "0x3210",
			args: args{
				input: []byte{0x10, 0x32},
			},
			want: 0x3210,
		},
		{
			name: "0x76543210",
			args: args{
				input: []byte{0x10, 0x32, 0x54, 0x76},
			},
			want: 0x76543210,
		},
		{
			name: "0x76543210",
			args: args{
				input: []byte{0x10, 0x32, 0x54, 0x76, 0x98},
			},
			want: 0x9876543210,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeUint64(tt.args.input); got != tt.want {
				t.Errorf("DecodeUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeUint(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			name: "0x3210",
			args: args{
				input: []byte{0x10, 0x32},
			},
			want: 0x3210,
		},
		{
			name: "0x76543210",
			args: args{
				input: []byte{0x10, 0x32, 0x54, 0x76},
			},
			want: 0x76543210,
		},
		{
			name: "0x76543210",
			args: args{
				input: []byte{0x10, 0x32, 0x54, 0x76, 0x98},
			},
			want: 0x9876543210,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeUint(tt.args.input); got != tt.want {
				t.Errorf("DecodeUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroPad(t *testing.T) {
	type args struct {
		input  []byte
		length int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Base Case",
			args: args{
				input:  []byte{1, 2, 3, 4},
				length: 0,
			},
			want: []byte{1, 2, 3, 4},
		},
		{
			name: "6 -> 8",
			args: args{
				input:  []byte{1, 2, 3, 4, 5, 6},
				length: 8,
			},
			want: []byte{1, 2, 3, 4, 5, 6, 0, 0},
		},
		{
			name: "4 -> 8",
			args: args{
				input:  []byte{1, 2, 3, 4},
				length: 8,
			},
			want: []byte{1, 2, 3, 4, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZeroPad(tt.args.input, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZeroPad() = %v, want %v", got, tt.want)
			}
		})
	}
}
