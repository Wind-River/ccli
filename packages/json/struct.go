package json

import "encoding/json"

type Profile struct {
	Profile   string  `json:"profile"`
	Label     string  `json:"label"`
	Format    float64 `json:"format"`
	Name      string  `json:"name,omitempty"`
	Version   string  `json:"version,omitempty"`
	FVC       string  `json:"fvc,omitempty"`
	Sha256    string  `json:"sha256"`
	CatalogID string  `json:"catalog_id,omitempty"`
}

type MainProfile struct {
	Profile
	InsertDate       string          `json:"insert_date,omitempty"`
	License          string          `json:"license,omitempty"`
	LicenseRationale json.RawMessage `json:"license_rationale,omitempty"`
	Size             int64           `json:"size,omitempty"`
	Aliases          []string        `json:"aliases,omitempty"`
	ComprisedOf      string          `json:"comprised_of,omitempty"`
	CompositeList    []string        `json:"composite_list,omitempty"`
}

type SecurityProfile struct {
	Profile
	CVEList []struct {
		ID          string   `json:"cve_id"`
		Description string   `json:"description"`
		Status      string   `json:"status"`
		Date        string   `json:"date"`
		Comments    []string `json:"comments,omitempty"`
		Link        string   `json:"link"`
	} `json:"cve_list"`
}

type QualityProfile struct {
	Profile
	BugList []struct {
		Name        string `json:"name"`
		ID          string `json:"id"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Level       string `json:"level"`
		Date        string `json:"date"`
		Link        string `json:"link"`
	} `json:"bug_list"`
}

type ProfileCollection struct {
	MainProfile
	SecurityProfile    *SecurityProfile
	QualityProfile     *QualityProfile
	UnexpectedProfiles map[string]json.RawMessage
}
