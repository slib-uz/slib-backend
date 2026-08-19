package response

type JobResponse struct {
	ID               uint    `json:"id"`
	CompanyINN       string  `json:"company_inn"`
	CompanyName      string  `json:"company_name"`
	ContractDate     string  `json:"contract_date"`
	ContractNumber   string  `json:"contract_number"`
	EndDate          *string `json:"end_date"`
	OrderDate        string  `json:"order_date"`
	OrderNumber      string  `json:"order_number"`
	PositionName     string  `json:"position_name"`
	StartDate        string  `json:"start_date"`
	StructureName    string  `json:"structure_name"`
	WorkplaceAddress string  `json:"workplace_address"`
	TransactionID    uint    `json:"transaction_id"`
	SoatoCode        string  `json:"soato_code"`
	WorkType         *int    `json:"work_type"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        *string `json:"updated_at"`
	Scientist        *uint   `json:"scientist"`
}
