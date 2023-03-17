package xhttp

import (
	"totem/pkg/xtls"
)

/**
 * @Author: jweny
 * @Author: https://github.com/jweny
 * @Date: 2021/12/2 16:36
 * @Desc:
 */

var HTTPOptions *ClientOptions

// ClientOptions http client options
type ClientOptions struct {
	Proxy string `mapstructure:"proxy" json:"proxy" yaml:"proxy" #:"漏洞扫描时使用的代理, 如: http://127.0.0.1:8080. 如需设置多个代理, 请使用 proxy_rule 或自行创建上层代理"`
	//ProxyRule           []Rule       `mapstructure:"proxy_rule" json:"proxy_rule" yaml:"proxy_rule" #:"漏洞扫描使用多个代理的配置规则, 具体请参照文档"`
	DialTimeout int `mapstructure:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout" #:"建立 tcp 连接的超时时间"`
	ReadTimeout int `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout" #:"读取 http 响应的超时时间, 不可太小, 否则会影响到 sql 时间盲注的判断"`
	// MaxRequestTimeout is the maximum time to wait for the request
	MaxRequestTimeout   int  `mapstructure:"max_request_timeout" json:"max_request_timeout" yaml:"max_request_timeout" #:"等待响应的最大时间(包含跳转)"`
	MaxConnsPerHost     int  `mapstructure:"max_conns_per_host" json:"max_conns_per_host" yaml:"max_conns_per_host" #:"同一 host 最大允许的连接数, 可以根据目标主机性能适当增大"`
	EnableHTTP2         bool `mapstructure:"enable_http2" json:"enable_http2" yaml:"enable_http2" #:"是否启用 http2, 开启可以提升部分网站的速度, 但目前不稳定有崩溃的风险"`
	IdleConnTimeout     int  `mapstructure:"idle_conn_timeout" json:"idle_conn_timeout" yaml:"idle_conn_timeout"`
	MaxIdleConns        int  `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	TLSHandshakeTimeout int  `mapstructure:"tls_handshake_timeout" json:"tls_handshake_timeout" yaml:"tls_handshake_timeout"`

	FailRetries                      int               `mapstructure:"fail_retries" json:"fail_retries" yaml:"fail_retries" #:"请求失败的重试次数, 0 则不重试"`
	MaxRedirect                      int               `mapstructure:"max_redirect" json:"max_redirect" yaml:"max_redirect" #:"单个请求最大允许的跳转数"`
	MaxRespBodySize                  int64             `mapstructure:"max_resp_body_size" json:"max_resp_body_size" yaml:"max_resp_body_size" #:"最大允许的响应大小, 默认 4M"`
	AllowMethods                     []string          `mapstructure:"allow_methods" json:"allow_methods" yaml:"allow_methods" #:"允许的请求方法"`
	Headers                          map[string]string `mapstructure:"headers" json:"headers" yaml:"headers" #:"自定义 headers"`
	Cookies                          map[string]string `mapstructure:"cookies" json:"cookies" yaml:"cookies" #:"自定义 cookies, 参考 headers 格式， key: value"`
	Debug                            bool              `mapstructure:"http_debug" json:"http_debug" yaml:"http_debug" #:"是否启用 debug 模式, 开启 request trace"`
	DisableKeepAlives                bool              `mapstructure:"disable_keep_alives" json:"disable_keep_alives" yaml:"disable_keep_alives" #:"是否禁用 keepalives"`
	CacheMaxNumber                   int               `mapstructure:"cache_max_number" json:"cache_max_number" yaml:"cache_max_number" #:"请求最大缓存个数"`
	CacheSingleRequestExpirationTime int               `mapstructure:"cache_single_request_expiration_time" json:"cache_single_request_expiration_time" yaml:"cache_single_request_expiration_time" #:"单个请求缓存过期时间"`

	TlsOptions *xtls.ClientOptions `mapstructure:"-" json:"-" yaml:"-"`
}

//type Server struct {
//	Options   string `json:"addr" yaml:"addr"`
//	Weight int    `json:"weight" yaml:"weight"`
//}
//
//type Rule struct {
//	Match   string   `json:"match" yaml:"match"`
//	Servers []Server `json:"servers" yaml:"servers"`
//}

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"

	// MethodConnect HTTP method
	MethodConnect = "CONNECT"

	// MethodTrace HTTP method
	MethodTrace = "TRACE"

	// MethodMove HTTP method
	MethodMove = "MOVE"

	// MethodPURGE MethodMove HTTP method
	MethodPURGE = "PURGE"

	//MethodPropFind = ""
)

func DefaultClientOptions() *ClientOptions {
	defaultHeaders := make(map[string]string)
	// 扫描引擎特征
	defaultHeaders["User-Agent"] = "threatbook-security-vulnerability-detection"
	defaultHeaders["Threatbook-Aspect"] = "threatbook-security-vulnerability-detection"

	return &ClientOptions{
		MaxRequestTimeout:   15,
		DialTimeout:         3,
		ReadTimeout:         30,
		IdleConnTimeout:     60,
		FailRetries:         0, // 默认改为0，否则如果配置文件指定了0，会不生效。 "nil value" 的问题
		MaxConnsPerHost:     20,
		MaxIdleConns:        50,
		TLSHandshakeTimeout: 5,
		MaxRedirect:         10,
		MaxRespBodySize:     2 << 20, // 4M
		//MaxQPS:              20,
		Headers: defaultHeaders,
		AllowMethods: []string{
			MethodHead,
			MethodGet,
			MethodPost,
			MethodPut,
			MethodPatch,
			MethodDelete,
			MethodOptions,
			MethodConnect,
			MethodTrace,
			MethodMove,
			MethodPURGE,
		},
		Cookies:                          make(map[string]string),
		EnableHTTP2:                      false,
		TlsOptions:                       xtls.DefaultClientOptions(),
		Debug:                            false,
		DisableKeepAlives:                true,
		CacheMaxNumber:                   2000,
		CacheSingleRequestExpirationTime: 180,
	}
}

func GetHTTPOptions() *ClientOptions {
	if HTTPOptions != nil {
		return HTTPOptions
	} else {
		return DefaultClientOptions()
	}
}
