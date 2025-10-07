package context

import (
	"net/http"
	"strings"

	"github.com/IBM/sarama"
)

// Headers is a map with headers that should be passed to request or will be passed to
// communication handler when incoming request is received. It's implementation
// similar to http.Header but assumes that it can be transformed to a different type
// of headers, including http.Header. It's API pretty similar to http.Header one.
type Headers map[string][]string

// NewHeadersFromHTTP creates new Headers type from http.Header. It will copy all
// passed headers, so original http.Header will be unaffected. Also all headers names
// will be converted to lowercase.
func NewHeadersFromHTTP(headers http.Header) Headers {
	hdrs := make(Headers)

	for key, values := range headers {
		k := strings.ToLower(key)
		hdrs[k] = make([]string, len(values))
		copy(hdrs[k], values)
	}

	return hdrs
}

// NewHeadersFromInterfaceMapWithDelimiter creates new Headers type from passed map. It
// will split header value by delimiter for getting multiple values. Also all headers names
// will be converted to lowercase.
func NewHeadersFromInterfaceMapWithDelimiter(headers map[string]interface{}, delimiter string) Headers {
	hdrs := make(Headers)

	for key, value := range headers {
		// Delimiter presence always assumes that we have a string behind interface{}.
		rawValue, ok := value.(string)
		if !ok {
			continue
		}

		if delimiter == "" {
			hdrs[strings.ToLower(key)] = []string{rawValue}
		} else {
			hdrs[strings.ToLower(key)] = strings.Split(rawValue, delimiter)
		}
	}

	return hdrs
}

// Add adds header. It appends passed value for key.
func (h Headers) Add(key, value string) {
	realKey := strings.ToLower(key)

	if h[realKey] == nil {
		h[realKey] = make([]string, 0)
	}

	h[realKey] = append(h[realKey], value)
}

// Get returns first value of header. If header isn't found then empty string is returned.
func (h Headers) Get(key string) string {
	if h == nil {
		return ""
	}

	value := h[strings.ToLower(key)]

	if len(value) == 0 {
		return ""
	}

	return value[0]
}

// Set sets one header value. It replaces any existing values.
func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = []string{value}
}

// Put replace header value. It replaces any existing values.
func (h Headers) Put(key string, value []string) {
	h[strings.ToLower(key)] = value
}

// ToHTTPHeader returns http.Header composed from our headers.
func (h Headers) ToHTTPHeader() http.Header {
	httpHeaders := make(http.Header)

	for key, values := range h {
		for _, value := range values {
			httpHeaders.Add(key, value)
		}
	}

	return httpHeaders
}

// ToInterfaceMap returns map[string]interface{} with strings as values. If header
// contains more than one value it will be joined using passed delimiter.
func (h Headers) ToInterfaceMap(delimiter string) map[string]interface{} {
	headersToReturn := make(map[string]interface{})

	for key, values := range h {
		headersToReturn[key] = strings.Join(values, delimiter)
	}

	return headersToReturn
}

// ToSaramaHeaders returns []sarama.RecordHeader - headers for Kafka. If header
// contains more than one value it will be joined using passed delimiter.
func (h Headers) ToSaramaHeaders(delimiterKey, delimiter string) []sarama.RecordHeader {
	headersToReturn := make([]sarama.RecordHeader, 0, len(h))

	for key, values := range h {
		headersToReturn = append(headersToReturn, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(strings.Join(values, delimiter)),
		})
	}

	if len(headersToReturn) > 0 {
		headersToReturn = append(headersToReturn, sarama.RecordHeader{
			Key:   []byte(delimiterKey),
			Value: []byte(delimiter),
		})
	}

	return headersToReturn
}

// Values returns all values for passed header. If header wasn't found - nil is returned.
// Returned slice is a copy that does not affect original headers.
func (h Headers) Values(key string) []string {
	header, found := h[strings.ToLower(key)]
	if !found {
		return nil
	}

	headerDataToReturn := make([]string, len(header))
	copy(headerDataToReturn, header)

	return headerDataToReturn
}
