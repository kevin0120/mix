package reader

const (
	WS_READER_UID = "WS_READER_UID"
	DISPATCHER_NFC_READER = "DISPATCHER_NFC_READER"
)

type ReaderUID struct {
	UID string `json:"uid"`
}
