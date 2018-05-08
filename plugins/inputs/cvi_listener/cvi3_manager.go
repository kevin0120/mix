package cvi_listener

type CVI3Manager struct {
	cvi3_clients   map[string]*CVI3Client
}

func (cm *CVI3Manager) StartService(configs	[]*CVIConfig) {
	cm.cvi3_clients = map[string]*CVI3Client{}
	for _, cvi3 := range configs {
		client := CVI3Client{}
		client.Config = cvi3

		cm.cvi3_clients[cvi3.SN] = &client
		go client.Start()
	}
}

func (cm *CVI3Manager) Pset(sn string, pset int, workorder_id int) {
	cm.cvi3_clients[sn].PSet(pset, workorder_id)
}
