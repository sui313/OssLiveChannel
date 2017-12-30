package OssLiveChannel

import (
	"io"
	//"log"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	//"github.com/denverdino/aliyungo/util"
)

const HeaderOSSPrefix = "x-oss-"

type request struct {
	method   string
	bucket   string
	path     string
	params   url.Values
	headers  http.Header
	baseurl  string
	payload  io.Reader
	prepared bool
	timeout  time.Duration
}

var ossParamsToSign = map[string]bool{
	"acl":                          true,
	"delete":                       true,
	"location":                     true,
	"logging":                      true,
	"notification":                 true,
	"partNumber":                   true,
	"policy":                       true,
	"requestPayment":               true,
	"torrent":                      true,
	"uploadId":                     true,
	"uploads":                      true,
	"versionId":                    true,
	"versioning":                   true,
	"versions":                     true,
	"response-content-type":        true,
	"response-content-language":    true,
	"response-expires":             true,
	"response-cache-control":       true,
	"response-content-disposition": true,
	"response-content-encoding":    true,
	"bucketInfo":                   true,
	"endTime":                      true,
	"startTime":                    true,
	"vod":                          true,
	"comp":                         true,
	"live":                         true,
	"status":                       true,
	"marker":                       true,
	"prefix":                       true,
	"max-keys":                     true,
}

func signRequest(request *request) {

	query := request.params

	urlSignature := query.Get("OSSAccessKeyId") != ""

	headers := request.headers
	contentMd5 := headers.Get("Content-Md5")
	contentType := headers.Get("Content-Type")
	date := ""
	if urlSignature {
		date = query.Get("Expires")
	} else {
		date = headers.Get("Date")
	}

	resource := request.path

	if request.bucket != "" {
		resource = "/" + request.bucket + request.path
	}

	params := make(url.Values)
	for k, v := range query {
		if ossParamsToSign[k] {
			params[k] = v
		}
	}

	if len(params) > 0 {
		resource = resource + "?" + Encode(params)
	}

	canonicalizedResource := resource

	_, canonicalizedHeader := canonicalizeHeader(headers)

	stringToSign := request.method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedHeader + canonicalizedResource

	fmt.Println("stringToSign: ", stringToSign)

	hmacSha1 := hmac.New(sha1.New, []byte(AccessKeySecret))
	hmacSha1.Write([]byte(stringToSign))
	sign := hmacSha1.Sum(nil)

	// Encode to Base64
	signature := base64.StdEncoding.EncodeToString(sign)

	request.headers.Set("Authorization", "OSS "+AccessKeyId+":"+signature)

}

func Encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			if v != "" {
				buf.WriteString("=")
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	return buf.String()
}

//Have to break the abstraction to append keys with lower case.
func canonicalizeHeader(headers http.Header) (newHeaders http.Header, result string) {
	var canonicalizedHeaders []string
	newHeaders = http.Header{}

	for k, v := range headers {
		if lower := strings.ToLower(k); strings.HasPrefix(lower, HeaderOSSPrefix) {
			newHeaders[lower] = v
			canonicalizedHeaders = append(canonicalizedHeaders, lower)
		} else {
			newHeaders[k] = v
		}
	}

	sort.Strings(canonicalizedHeaders)

	var canonicalizedHeader string

	for _, k := range canonicalizedHeaders {
		canonicalizedHeader += k + ":" + headers.Get(k) + "\n"
	}

	return newHeaders, canonicalizedHeader
}
