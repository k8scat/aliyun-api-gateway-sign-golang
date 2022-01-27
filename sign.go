package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	HTTPHeaderAccept      = "Accept"
	HTTPHeaderContentMD5  = "Content-MD5"
	HTTPHeaderContentType = "Content-Type"
	HTTPHeaderDate        = "Date"

	HTTPHeaderCAPrefix           = "X-Ca-"
	HTTPHeaderCASignature        = "X-Ca-Signature"
	HTTPHeaderCATimestamp        = "X-Ca-Timestamp"
	HTTPHeaderCANonce            = "X-Ca-Nonce"
	HTTPHeaderCAKey              = "X-Ca-Key"
	HTTPHeaderCASignatureHeaders = "X-Ca-Signature-Headers"

	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeStream = "application/octet-stream"
	ContentTypeJSON   = "application/json"
	ContentTypeXML    = "application/xml"
	ContentTypeText   = "application/text"

	LF = "\n"
)

type APIGateway struct {
	Key    string
	Secret string
}

func NewAPIGateway(key, secret string) *APIGateway {
	return &APIGateway{
		Key:    key,
		Secret: secret,
	}
}

func (api *APIGateway) Sign(req *http.Request) error {
	t := time.Now().UnixNano()
	req.Header.Set(HTTPHeaderCATimestamp, strconv.FormatInt(t/1000000, 10))
	req.Header.Set(HTTPHeaderCANonce, strconv.FormatInt(t, 10))
	req.Header.Set(HTTPHeaderCAKey, api.Key)

	s, signatureHeaders, err := api.buildSignature(req)
	if err != nil {
		return err
	}

	req.Header.Set(HTTPHeaderCASignature, api.sha256Hmac([]byte(s)))
	req.Header.Set(HTTPHeaderCASignatureHeaders, strings.Join(signatureHeaders, ","))
	return nil
}

func (api *APIGateway) sha256Hmac(b []byte) string {
	h := hmac.New(sha256.New, []byte(api.Secret))
	h.Write(b)
	return base64.StdEncoding.EncodeToString(h.Sum([]byte{}))
}

func (api *APIGateway) buildSignature(req *http.Request) (string, []string, error) {
	var buf strings.Builder

	buf.WriteString(strings.ToUpper(req.Method))
	buf.WriteString(LF)
	buf.WriteString(req.Header.Get(HTTPHeaderAccept))
	buf.WriteString(LF)
	buf.WriteString(req.Header.Get(HTTPHeaderContentMD5))
	buf.WriteString(LF)
	buf.WriteString(req.Header.Get(HTTPHeaderContentType))
	buf.WriteString(LF)
	buf.WriteString(req.Header.Get(HTTPHeaderDate))
	buf.WriteString(LF)

	s, signatureHeaders, err := api.buildHeader(req.Header)
	if err != nil {
		return "", nil, err
	}
	buf.WriteString(s)

	path, err := api.buildPath(req)
	if err != nil {
		return "", nil, err
	}
	buf.WriteString(path)

	return buf.String(), signatureHeaders, nil
}

func (api *APIGateway) buildHeader(header http.Header) (string, []string, error) {
	signatureHeaders := make([]string, 0)
	var buf strings.Builder
	for _, k := range sortedKeys(url.Values(header)) {
		if strings.HasPrefix(k, HTTPHeaderCAPrefix) {
			buf.WriteString(k)
			buf.WriteString(":")
			buf.WriteString(header.Get(k))
			buf.WriteString(LF)
			signatureHeaders = append(signatureHeaders, k)
		}
	}
	return buf.String(), signatureHeaders, nil
}

func (api *APIGateway) buildPath(req *http.Request) (string, error) {
	data := make(url.Values)
	for k, v := range req.URL.Query() {
		data[k] = make([]string, 0)
		if len(v) > 0 {
			data[k] = append(data[k], v[0])
		}
	}
	if strings.ToUpper(req.Method) == http.MethodPost &&
		req.Header.Get(HTTPHeaderContentType) == ContentTypeForm &&
		req.Body != nil {
		err := req.ParseForm()
		if err != nil {
			return "", err
		}
		for k, v := range req.Form {
			copy(data[k], v)
		}
	}

	var buf strings.Builder
	buf.WriteString(req.URL.Path)
	if len(data) > 0 {
		buf.WriteString("?")
	}
	for i, k := range sortedKeys(data) {
		buf.WriteString(k)
		v := data.Get(k)
		if v != "" {
			buf.WriteString("=")
			buf.WriteString(v)
		}
		if i != len(data)-1 {
			buf.WriteString("&")
		}
	}
	return buf.String(), nil
}

func sortedKeys(values url.Values) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
