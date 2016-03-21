package rtorrent

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClientDownloadTotal(t *testing.T) {
	wantTotal := 1024

	c, done := testClient(t, "get_down_total", nil, wantTotal)
	defer done()

	total, err := c.DownloadTotal()
	if err != nil {
		t.Fatalf("failed call to Client.DownloadTotal: %v", err)
	}

	if want, got := wantTotal, total; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected total download bytes:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestClientUploadTotal(t *testing.T) {
	wantTotal := 1024

	c, done := testClient(t, "get_up_total", nil, wantTotal)
	defer done()

	total, err := c.UploadTotal()
	if err != nil {
		t.Fatalf("failed call to Client.UploadTotal: %v", err)
	}

	if want, got := wantTotal, total; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected total upload bytes:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestClientDownloadRate(t *testing.T) {
	wantRate := 1024

	c, done := testClient(t, "get_down_rate", nil, wantRate)
	defer done()

	rate, err := c.DownloadRate()
	if err != nil {
		t.Fatalf("failed call to Client.DownloadRate: %v", err)
	}

	if want, got := wantRate, rate; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected download rate:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestClientUploadRate(t *testing.T) {
	wantRate := 1024

	c, done := testClient(t, "get_up_rate", nil, wantRate)
	defer done()

	rate, err := c.UploadRate()
	if err != nil {
		t.Fatalf("failed call to Client.UploadRate: %v", err)
	}

	if want, got := wantRate, rate; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected upload rate:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func testClient(t *testing.T, method string, args []string, out interface{}) (*Client, func()) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var xr xmlrpcRequest
		if err := xml.NewDecoder(r.Body).Decode(&xr); err != nil {
			t.Fatalf("failed to decode XML-RPC body: %v", err)
		}

		if want, got := method, xr.MethodName; want != got {
			t.Fatalf("unexpected XML-RPC method name:\n- want: %q\n-  got: %q",
				want, got)
		}

		if want, got := len(args), len(xr.Params); want != got {
			t.Fatalf("unexpected number of XML-RPC parameters:\n- want: %v\n-  got: %v",
				want, got)
		}

		params := make([]string, 0, len(xr.Params))
		for _, p := range xr.Params {
			params = append(params, p.Param.Value.String)
		}

		if args == nil {
			args = make([]string, 0)
		}

		if want, got := args, params; !reflect.DeepEqual(want, got) {
			t.Fatalf("unexpected XML-RPC parameters:\n- want: %v\n-  got: %v",
				want, got)
		}

		if err := writeXMLRPC(w, out); err != nil {
			t.Fatalf("unexpected error encoding XML-RPC response: %v", err)
		}
	}))

	c, err := New(s.URL, nil)
	if err != nil {
		t.Fatalf("failed to create Client: %v", err)
	}

	done := func() {
		if err := c.Close(); err != nil {
			t.Fatalf("failed to clean up Client: %v", err)
		}

		s.Close()
	}

	return c, done
}

// XML-RPC helper routines and structures

func writeXMLRPC(w io.Writer, out interface{}) error {
	var xw xmlrpcResponse
	xw.Params = make([]xmlrpcParam, 1, 1)

	switch out := out.(type) {
	case int:
		xw.Params[0].Param.Value.Int = out
	case string:
		xw.Params[0].Param.Value.String = out
	case []string:
		xw.Params[0].Param.Value.Array = new(xmlrpcArray)
		xw.Params[0].Param.Value.Array.Data.Value = make([]xmlrpcArrayData, len(out))

		for i, s := range out {
			xw.Params[0].Param.Value.Array.Data.Value[i].String = s
		}
	}

	// Inspect buf for debugging if needed
	buf := bytes.NewBuffer(nil)
	mw := io.MultiWriter(w, buf)

	err := xml.NewEncoder(mw).Encode(xw)
	return err
}

type xmlrpcRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	MethodName string   `xml:"methodName"`

	Params []xmlrpcParam `xml:"params"`
}

type xmlrpcResponse struct {
	XMLName xml.Name `xml:"methodResponse"`

	Params []xmlrpcParam `xml:"params"`
}

type xmlrpcParam struct {
	Param struct {
		Value struct {
			Array  *xmlrpcArray `xml:"array,omitempty"`
			Int    int          `xml:"i8,omitempty"`
			String string       `xml:"string,omitempty"`
		} `xml:"value"`
	} `xml:"param"`
}

type xmlrpcArray struct {
	Data struct {
		Value []xmlrpcArrayData `xml:"value"`
	} `xml:"data"`
}

type xmlrpcArrayData struct {
	String string `xml:"string"`
}
