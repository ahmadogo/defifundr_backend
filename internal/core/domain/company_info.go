package domain

type CompanyInfo struct {
	CompanyName         *string `json:"company_name,omitempty"`
	CompanySize         *string `json:"company_size,omitempty"`
	CompanyIndustry     *string `json:"company_industry,omitempty"`
	CompanyDescription  *string `json:"company_description,omitempty"`
	CompanyHeadquarters *string `json:"company_headquarters,omitempty"`
}
