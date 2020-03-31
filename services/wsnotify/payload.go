package wsnotify

type WSRegist struct {
	HMI_SN string `json:"hmi_sn"`
}

type WSRegistMsg struct {
	Msg string `json:"msg"`
}
