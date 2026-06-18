// Go NETCONF Client
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type nilCloser struct {
	io.Reader
	io.Writer
}

func newNilCloser(r io.Reader, w io.Writer) *nilCloser {
	return &nilCloser{r, w}
}

func (nc *nilCloser) Close() error {
	return nil
}

type transportTest struct {
	TransportBasicIO
}

func newTransportTest(input string) (*transportTest, *bytes.Buffer) {
	testReader := bytes.NewReader([]byte(input))
	testWriter := new(bytes.Buffer)

	var t transportTest
	t.ReadWriteCloser = newNilCloser(testReader, testWriter)
	return &t, testWriter
}

func TestReceiveHello(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected *HelloMessage
	}{
		{
			name: "juniperHello",
			input: `<!-- No zombies were killed during the creation of this user interface -->
<!-- user bbennett, class j-super-user -->
<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <capabilities>
    <capability>urn:ietf:params:xml:ns:netconf:base:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:candidate:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:confirmed-commit:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:validate:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:url:1.0?protocol=http,ftp,file</capability>
    <capability>http://xml.juniper.net/netconf/junos/1.0</capability>
    <capability>http://xml.juniper.net/dmi/system/1.0</capability>
  </capabilities>
  <session-id>19313</session-id>
</hello>
]]>]]>`,
			expected: &HelloMessage{
				XMLName:   xml.Name{Space: "urn:ietf:params:xml:ns:netconf:base:1.0", Local: "hello"},
				SessionID: 19313,
				Capabilities: []string{
					"urn:ietf:params:xml:ns:netconf:base:1.0",
					"urn:ietf:params:xml:ns:netconf:capability:candidate:1.0",
					"urn:ietf:params:xml:ns:netconf:capability:confirmed-commit:1.0",
					"urn:ietf:params:xml:ns:netconf:capability:validate:1.0",
					"urn:ietf:params:xml:ns:netconf:capability:url:1.0?protocol=http,ftp,file",
					"http://xml.juniper.net/netconf/junos/1.0",
					"http://xml.juniper.net/dmi/system/1.0",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			trans, _ := newTransportTest(tc.input)

			hello, err := trans.ReceiveHello()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !cmp.Equal(hello, tc.expected) {
				t.Errorf("unexpected hello message%s", cmp.Diff(hello, tc.expected))
			}
		})
	}
}

func TestSendHello(t *testing.T) {
	tt := []struct {
		name     string
		input    *HelloMessage
		expected string
	}{
		{
			name:  "default",
			input: &HelloMessage{Capabilities: DefaultCapabilities},
			expected: `<?xml version="1.0" encoding="UTF-8"?>
<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><capabilities><capability>urn:ietf:params:netconf:base:1.0</capability><capability>urn:ietf:params:netconf:base:1.1</capability></capabilities></hello>]]>]]>`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			trans, out := newTransportTest("")
			trans.SendHello(tc.input)
			rawHello := out.String()

			if rawHello != tc.expected {
				t.Errorf("unexpected result: (want %q, got %q)", tc.expected, rawHello)
			}
		})
	}
}

// Login test needs to be over 4096 bytes to fully test the function
var loginText = `
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

SRX240 (ttyp2)

Password:

--- JUNOS 12.1X45-D15.5 built 2013-09-19 07:42:15 UTC
bbennett@SRX240>

`

func TestWaitForRegexp(t *testing.T) {
	tt, _ := newTransportTest(loginText)

	var promptRE = regexp.MustCompile(`([>%])\s+`)
	output, matches, err := tt.WaitForRegexp(promptRE)

	if err != nil {
		t.Errorf("WaitForRegexp failed: %s", err)
	}

	if len(matches) != 1 {
		t.Errorf("WaitForRegexp Length of regexp matches is not equal:  Expecting '%d', got '%d", 1, len(matches))
	}

	if !bytes.Equal(matches[0], []byte(">")) {
		t.Errorf("WaitForRegexp #0 match not equal:  Expecting '%d', got '%d", '>', matches[0])
	}

	if !bytes.Equal(output, []byte(loginText)) {
		t.Errorf("WaitForRegexp output text not equal:  Expecting '%s', got '%s", loginText, output)
	}
}

var waitforStringResponse = `
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

SRX240 (ttyp2)

`

func TestWaitForString(t *testing.T) {
	tt, _ := newTransportTest(loginText)

	output, err := tt.WaitForString("Password:")

	if err != nil {
		t.Errorf("WaitForString failed: %s", err)
	}

	if output != waitforStringResponse {
		t.Errorf("WaitForRegexp output text not equal:  Expecting '%s', got '%s", waitforStringResponse, output)
	}
}

func TestWaitForBytesEmpty(t *testing.T) {
	tt, _ := newTransportTest("")

	_, err := tt.WaitForBytes([]byte("Test"))
	if err == nil {
		t.Errorf("WaitForBytes should error on empty input!")
	}
}

// chunkReader は固定サイズのチャンクでデータを返し、最後に finalErr を返す io.Reader。
// デリミタ欠損シナリオ（サーバがボディのみ送信して接続を切断）の再現に使用する。
type chunkReader struct {
	data     []byte
	chunkSz  int
	offset   int
	finalErr error
}

func (r *chunkReader) Read(b []byte) (int, error) {
	if r.offset >= len(r.data) {
		return 0, r.finalErr
	}
	end := r.offset + r.chunkSz
	if end > len(r.data) {
		end = len(r.data)
	}
	n := copy(b, r.data[r.offset:end])
	r.offset += n
	return n, nil
}

// multiChunkReader は複数のチャンクを順番に返す io.Reader。
// デリミタ分割シナリオ（TCPセグメント境界でデリミタが分断）の再現に使用する。
type multiChunkReader struct {
	chunks [][]byte
	idx    int
}

func (r *multiChunkReader) Read(b []byte) (int, error) {
	if r.idx >= len(r.chunks) {
		return 0, io.EOF
	}
	n := copy(b, r.chunks[r.idx])
	r.idx++
	return n, nil
}

// TestWaitForFuncDelimiterMissing はデリミタなしでボディのみ送信された場合、
// WaitForBytes がエラーを返すことを確認する（ループが無限継続しないことの検証）。
func TestWaitForFuncDelimiterMissing(t *testing.T) {
	body := `<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="1">
<ok/>
</rpc-reply>
`
	r := &chunkReader{
		data:     []byte(body),
		chunkSz:  32,
		finalErr: io.ErrUnexpectedEOF,
	}

	var trans transportTest
	trans.ReadWriteCloser = newNilCloser(r, new(bytes.Buffer))

	_, err := trans.WaitForBytes([]byte(msgSeperator))
	if err == nil {
		t.Error("WaitForBytes should return error when delimiter is missing")
	}
}

// TestWaitForFuncDelimiterSplit はデリミタ ]]>]]> がTCPセグメント境界で
// ]]> + ]]> に分割されても正しく検出・組み立てできることを確認する。
func TestWaitForFuncDelimiterSplit(t *testing.T) {
	prefix := `<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="1">
<ok/>
</rpc-reply>
`
	delim := msgSeperator
	splitAt := len(delim) / 2 // ]]> と ]]> に分割

	chunks := [][]byte{
		[]byte(prefix + delim[:splitAt]),
		[]byte(delim[splitAt:]),
	}

	r := &multiChunkReader{chunks: chunks}
	var trans transportTest
	trans.ReadWriteCloser = newNilCloser(r, new(bytes.Buffer))

	result, err := trans.WaitForBytes([]byte(delim))
	if err != nil {
		t.Fatalf("WaitForBytes failed unexpectedly: %v", err)
	}

	expected := []byte(prefix)
	if !bytes.Equal(result, expected) {
		t.Errorf("unexpected result:\nwant: %q\ngot:  %q", expected, result)
	}
}
