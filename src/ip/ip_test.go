package ip

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type testReadSeeker struct {
	s string
	i int64 // current reading index
}

// Read implements the [io.Reader] interface.
func (r *testReadSeeker) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// Seek implements the [io.Seeker] interface.
func (r *testReadSeeker) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

func TestMonitorAndUpdate(t *testing.T) {
	type args struct {
		c *Config
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MonitorAndUpdate(tt.args.c)
		})
	}
}

func Test_handleIPCheck(t *testing.T) {
	type args struct {
		c *config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handleIPCheck(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("handleIPCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ipHasChanged(t *testing.T) {
	type args struct {
		newIp string
		r     io.ReadSeeker
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ip_has_indeed_changed",
			args: args{
				newIp: "192.168.169.169",
				r: &testReadSeeker{
					s: "100.10.1.4",
					i: 0,
				},
			},
			want: true,
		},
		{
			name: "ip_hasn't_changed",
			args: args{
				newIp: "192.168.169.169",
				r: &testReadSeeker{
					s: "192.168.169.169",
					i: 0,
				},
			},
			want: false,
		},
		{
			name: "ip_has_indeed_changed2",
			args: args{
				newIp: "1.1.1.1",
				r: &testReadSeeker{
					s: "100.10.1.4",
					i: 123,
				},
			},
			want: true,
		},
		{
			name: "ip_hasn't_changed2",
			args: args{
				newIp: "8.8.8.8",
				r: &testReadSeeker{
					s: "8.8.8.8",
					i: 456,
				},
			},
			want: false,
		},
		{
			name: "empty",
			args: args{
				newIp: "99.99.99.99",
				r: &testReadSeeker{
					s: "",
					i: 0,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipHasChanged(tt.args.newIp, tt.args.r); got != tt.want {
				t.Errorf("ipHasChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateIP(t *testing.T) {
	type args struct {
		f  *os.File
		ip string
	}
	tests := []struct {
		name         string
		args         args
		want         string
		wantErr      bool
		fileContents []byte
	}{
		{
			name:         "overwriting_IP",
			args:         args{f: nil, ip: "192.168.169.169"},
			want:         "192.168.169.169",
			wantErr:      false,
			fileContents: []byte("100.010.1.04"),
		},
		{
			name:         "empty_file",
			args:         args{f: nil, ip: "192.168.169.169"},
			want:         "192.168.169.169",
			wantErr:      false,
			fileContents: []byte(""),
		},
		{
			name:    "overwriting_chunk_of_text",
			args:    args{f: nil, ip: "192.168.169.169"},
			want:    "192.168.169.169",
			wantErr: false,
			fileContents: []byte(
				"What the fuck did you just fucking say about me, you little bitch? I’ll have you know I graduated top of my class in the Navy Seals, and I’ve been involved in numerous secret raids on Al-Quaeda, and I have over 300 confirmed kills. I am trained in gorilla warfare and I’m the top sniper in the entire US armed forces. You are nothing to me but just another target. I will wipe you the fuck out with precision the likes of which has never been seen before on this Earth, mark my fucking words. You think you can get away with saying that shit to me over the Internet? Think again, fucker. As we speak I am contacting my secret network of spies across the USA and your IP is being traced right now so you better prepare for the storm, maggot. The storm that wipes out the pathetic little thing you call your life. You’re fucking dead, kid. I can be anywhere, anytime, and I can kill you in over seven hundred ways, and that’s just with my bare hands. Not only am I extensively trained in unarmed combat, but I have access to the entire arsenal of the United States Marine Corps and I will use it to its full extent to wipe your miserable ass off the face of the continent, you little shit. If only you could have known what unholy retribution your little “clever” comment was about to bring down upon you, maybe you would have held your fucking tongue. But you couldn’t, you didn’t, and now you’re paying the price, you goddamn idiot. I will shit fury all over you and you will drown in it. You’re fucking dead, kiddo. ",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const filename string = "test.txt"

			if tt.fileContents != nil {
				f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				tt.args.f = f

				_, err = f.Write(tt.fileContents)
				if err != nil {
					t.Fatal(err)
				}

			}
			if err := updateIP(tt.args.f, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("updateIP() error = %v, wantErr %v", err, tt.wantErr)
			}

			b, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.EqualFold(string(b), tt.want) {
				t.Errorf("got bytes <%v>, wanted <%v>\n", string(b), tt.want)
			}

			t.Cleanup(func() {
				os.Remove(filename)
			})
		})
	}
}
