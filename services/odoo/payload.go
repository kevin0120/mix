package odoo

type ODOOPR struct {
	Pr_group string `json:"pr_group"`
	Pr_value string `json:"pr_value"`
}

type ODOOMO struct {
	Pin_check_code int    `json:"pin_check_code"`
	Equipment_name string `json:"equipment_name"`
	Factory_name   string `json:"factory_name"`
	Pin            int    `json:"pin"`
	Year           int    `json:"year"`
	Assembly_line  string `json:"assembly_line"`
	Model          string `json:"model"`
	Vin            string `json:"vin"`
	Lnr            string `json:"lnr"`
	//Date_planned_start string   `json:"date_planned_start"`
	Prs []ODOOPR `json:"prs"`
}
