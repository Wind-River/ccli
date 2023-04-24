package yaml

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
		Comments    string   `yaml:"comments,omitempty"`
		Link        string   `yaml:"link"`
		References  []string `yaml:"references"`
	} `yaml:"cve_list"`
}
