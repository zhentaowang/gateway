package conf

type T struct {
    // Maximum number of active connections which may be established to server
    MaxConns int `yaml:"max_conns"`
    // MaxIdle maximum number of idle connections
    MaxIdle int `yaml:"max_idle"`
    // MaxConnDuration Keep-alive connections are closed after this duration.
    MaxConnDuration int `yaml:"max_conn_duration"`
    // MaxIdleConnDuration Idle keep-alive connections are closed after this duration.
    MaxIdleConnDuration int `yaml:"max_idle_conn_duration"`
    // ReadBufferSize Per-connection buffer size for responses' reading.
    ReadBufferSize int `yaml:"read_buffer_size"`
    // WriteBufferSize Per-connection buffer size for requests' writing.
    WriteBufferSize int `yaml:"write_buffer_size"`
    // ReadTimeout Maximum duration for full response reading (including body).
    ReadTimeout int `yaml:"read_timeout"`
    // WriteTimeout Maximum duration for full request writing (including body).
    WriteTimeout int `yaml:"write_timeout"`
    // MaxResponseBodySize Maximum response body size.
    MaxResponseBodySize int `yaml:"max_response_body_size"`
}

var TConfig T