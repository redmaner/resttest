package resttest

import "net/http"

type Transport struct {
	http.RoundTripper
	headers []HttpHeader
}

func NewTransport(headers []HttpHeader) (*Transport, error) {
	t := http.DefaultTransport.(*http.Transport).Clone()

	return &Transport{
		RoundTripper: t,
		headers:      headers,
	}, nil
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	for i := range t.headers {
		r.Header.Set(t.headers[i].Key, t.headers[i].Value)
	}

	return t.RoundTripper.RoundTrip(r)
}
