package yaml

type Part struct {
	Format      float64 `yaml:"format"`
	FVC         string  `yaml:"fvc"`
	Sha256      string  `yaml:"sha256"`
	CatalogID   string  `yaml:"catalog_id"`
	Name        string  `yaml:"name"`
	Version     string  `yaml:"version"`
	Label       string  `yaml:"label"`
	Description string  `yaml:"description"`
	License     struct {
		LicenseExpression string `yaml:"license_expression"`
		AnalysisType      string `yaml:"analysis_type"`
	} `yaml:"license"`
	Size          string   `yaml:"size"`
	Aliases       []string `yaml:"aliases"`
	ComprisedOf   string   `yaml:"comprised_of"`
	CompositeList []string `yaml:"composite_list"`
}

type Profile struct {
	Profile   string  `yaml:"profile"`
	Format    float64 `yaml:"format"`
	Name      string  `yaml:"name"`
	Version   string  `yaml:"version"`
	FVC       string  `yaml:"fvc"`
	Sha256    string  `yaml:"sha256"`
	CatalogID string  `yaml:"catalog_id"`
}

type SecurityProfile struct {
	CVEList []struct {
		ID          string   `yaml:"cve_id"`
		Description string   `yaml:"description"`
		Status      string   `yaml:"status"`
		Date        string   `yaml:"date"`
		Comments    string   `yaml:"comments"`
		Link        string   `yaml:"link"`
		References  []string `yaml:"references"`
	} `yaml:"cve_list"`
}

type QualityProfile struct {
	BugList []struct {
		Name        string   `yaml:"name"`
		ID          string   `yaml:"id"`
		Description string   `yaml:"description"`
		Status      string   `yaml:"status"`
		Level       string   `yaml:"level"`
		Date        string   `yaml:"date"`
		Link        string   `yaml:"link"`
		Comments    string   `yaml:"comments"`
		References  []string `yaml:"references"`
	} `yaml:"bug_list"`
}

type LicensingProfile struct {
	LicenseAnalysis []License `yaml:"license_analysis"`
	Copyrights      []string  `yaml:"copyrights"`
	LegalNotice     string    `yaml:"legal_notice"`
}

type License struct {
	LicenseExpression string   `yaml:"license_expression"`
	AnalysisType      string   `yaml:"analysis_type"`
	LegalNotice       string   `yaml:"legal_notice"`
	OtherLegalNotices []string `yaml:"other_legal_notices"`
}
