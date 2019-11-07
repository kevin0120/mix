package scanner

const (
	WS_SCANNER_STATUS = "WS_SCANNER_STATUS"
	WS_SCANNER_READ   = "WS_SCANNER_READ"
)

type ScannerRead struct {
	Src     string `json:"src"`
	SN      string `json:"sn"`
	Barcode string `json:"barcode"`
}

type ScannerStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
