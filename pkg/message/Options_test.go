package messages

import (
	"reflect"
	"testing"
)

func TestEncodeSingleOption(t *testing.T) {
	type args struct {
		delta uint
		b     []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Small Delta and Small Length",
			args: args{
				delta: 0x2,
				b:     []byte{1},
			},
			want:    []byte{33, 1},
			wantErr: false,
		},
		{
			name: "Large Delta and Small Length",
			args: args{
				delta: 0x0D,
				b:     []byte{0x1},
			},
			want:    []byte{0xD1, 0x00, 0x1},
			wantErr: false,
		},
		{
			name: "Large Delta With Difference and Small Length",
			args: args{
				delta: 0x0F,
				b:     []byte{0x1},
			},
			want:    []byte{0xD1, 0x02, 0x1},
			wantErr: false,
		},
		{
			name: "Large Delta With Difference and Small Length",
			args: args{
				delta: 0x10E,
				b:     []byte{0x1},
			},
			want:    []byte{0xE1, 0x01, 0x1},
			wantErr: false,
		},
		{
			name: "Small Delta and Large Length",
			args: args{
				delta: 0x01,
				b:     []byte{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			},
			want:    []byte{0x1D, 0x01, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeSingleOption(tt.args.delta, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeSingleOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeSingleOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testHost = "coap://test.com"

func TestOptions_SetURI(t *testing.T) {
	type fields struct {
		ContentFormat uint
		ETag          [][]byte
		LocationPath  []string
		LocationQuery []string
		MaxAge        uint
		ProxyURI      *string
		ProxyScheme   *string
		URIHost       *string
		URIPath       []string
		URIPort       uint
		URIQuery      []string
		Accept        uint
		IfMatch       [][]byte
		IfNoneMatch   bool
		Size1         uint
	}
	type args struct {
		rawurl string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		wantURL  string
		wantPort uint
		wantPath []string
	}{
		{
			name: "Setting Normal URL",
			fields: fields{
				URIHost: &testHost,
			},
			args: args{
				rawurl: "coap://test.com",
			},
			wantErr:  false,
			wantURL:  "test.com",
			wantPath: []string{},
			wantPort: 5683,
		},
		{
			name: "Changing port number",
			fields: fields{
				URIHost: &testHost,
			},
			args: args{
				rawurl: "coap://test.com:80",
			},
			wantErr:  false,
			wantURL:  "test.com",
			wantPath: []string{},
			wantPort: 80,
		},
		{
			name: "Setting Normal URL With Path",
			fields: fields{
				URIHost: &testHost,
			},
			args: args{
				rawurl: "coap://test.com/a/path",
			},
			wantErr:  false,
			wantURL:  "test.com",
			wantPath: []string{"a", "path"},
			wantPort: 5683,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				ContentFormat: tt.fields.ContentFormat,
				ETag:          tt.fields.ETag,
				LocationPath:  tt.fields.LocationPath,
				LocationQuery: tt.fields.LocationQuery,
				MaxAge:        tt.fields.MaxAge,
				ProxyURI:      tt.fields.ProxyURI,
				ProxyScheme:   tt.fields.ProxyScheme,
				URIHost:       tt.fields.URIHost,
				URIPath:       tt.fields.URIPath,
				URIPort:       tt.fields.URIPort,
				URIQuery:      tt.fields.URIQuery,
				Accept:        tt.fields.Accept,
				IfMatch:       tt.fields.IfMatch,
				IfNoneMatch:   tt.fields.IfNoneMatch,
				Size1:         tt.fields.Size1,
			}
			if err := o.SetURI(tt.args.rawurl); (err != nil) != tt.wantErr {
				t.Errorf("Options.SetURI() error = %v, wantErr %v", err, tt.wantErr)
			}

			if *o.URIHost != tt.wantURL {
				t.Errorf("Options.SetURI() HostURI = %v, wantHostURI %v", *o.URIHost, tt.wantURL)
			}

			if o.URIPort != tt.wantPort {
				t.Errorf("Options.SetURI() HostPort = %v, wantPort %v", o.URIPort, tt.wantPort)
			}

			if !reflect.DeepEqual(o.URIPath, tt.wantPath) {
				t.Errorf("Options.SetURI() URIPath = %v, wantPath %v", o.URIPath, tt.wantPath)
			}

		})
	}
}
