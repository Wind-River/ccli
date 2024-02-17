// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.
package yaml

// struct for storing part data
type Part struct {
	Format      float64 `yaml:"format"`
	FVC         string  `yaml:"fvc"`
	Sha256      string  `yaml:"sha256"`
	CatalogID   string  `yaml:"catalog_id"`
	Name        string  `yaml:"name"`
	Version     string  `yaml:"version"`
	Type        string  `yaml:"type"`
	ContentType string  `yaml:"content_type"`
	FamilyName  string  `yaml:"family_name"`
	Label       string  `yaml:"label"`
	Description string  `yaml:"description"`
	HomePage    string  `yaml:"home_page"`
	License     struct {
		LicenseExpression string `yaml:"license_expression"`
		AnalysisType      string `yaml:"analysis_type"`
	} `yaml:"license"`
	Size          string   `yaml:"size"`
	Aliases       []string `yaml:"aliases"`
	ComprisedOf   string   `yaml:"comprised_of"`
	CompositeList []string `yaml:"composite_list"`
}

// struct for storing profile data
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
	CVEList []CVE `yaml:"cve_list" json:"cve_list"`
}

type QualityProfile struct {
	BugList []Bug `yaml:"bug_list" json:"bug_list"`
}

type LicensingProfile struct {
	LicenseAnalysis   []License `yaml:"license_analysis" json:"license_analysis"`
	Copyrights        []string  `yaml:"copyrights" json:"copyrights"`
	LegalNotice       string    `yaml:"legal_notice" json:"legal_notices"`
	OtherLegalNotices []string  `yaml:"other_legal_notices" json:"other_legal_notices"`
}

type License struct {
	LicenseExpression string `yaml:"license_expression" json:"license_expression"`
	AnalysisType      string `yaml:"analysis_type" json:"analysis_type"`
	Comments          string `yaml:"comments" json:"comments"`
}

type CVE struct {
	ID          string   `yaml:"cve_id" json:"cve_id"`
	Description string   `yaml:"description" json:"description"`
	Status      string   `yaml:"status" json:"status"`
	Date        string   `yaml:"date" json:"date"`
	Comments    string   `yaml:"comments" json:"comments"`
	Link        string   `yaml:"link" json:"link"`
	References  []string `yaml:"references" json:"references"`
}

type Bug struct {
	Name        string   `yaml:"name" json:"name"`
	ID          string   `yaml:"id" json:"id"`
	Description string   `yaml:"description" json:"description"`
	Status      string   `yaml:"status" json:"status"`
	Level       string   `yaml:"level" json:"level"`
	Date        string   `yaml:"date" json:"date"`
	Link        string   `yaml:"link" json:"link"`
	Comments    string   `yaml:"comments" json:"comments"`
	References  []string `yaml:"references" json:"references"`
}
