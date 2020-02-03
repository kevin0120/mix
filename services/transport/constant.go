package transport

const (
	HEADER_SUBJECT   = "X-Subject"
	HEADER_REPLY     = "X-Reply-To"
	HEADER_MSG_ID    = "X-Message-ID"
	HEADER_CLIENT_ID = "X-Client-ID"
)

const (
	GRPCTransport   = "grpc"
	BrokerTransport = "broker"
	HttpTransport   = "http"
)
