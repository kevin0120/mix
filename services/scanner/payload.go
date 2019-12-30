package scanner

const (
	WS_SCANNER_STATUS = "WS_SCANNER_STATUS"
)

type ScannerRead struct {
	Src     string `json:"src"`
	SN      string `json:"sn"`
	Barcode string `json:"barcode"`
}
