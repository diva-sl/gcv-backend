package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Metric struct {
	Label string `json:"label" bson:"label"`
	Value string `json:"value" bson:"value"`
}

type Screenshot struct {
	Label string `json:"label" bson:"label"`
	Path  string `json:"path" bson:"path"`
	Type  string `json:"type" bson:"type"` // "storefront" or "admin"
}

type Project struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProjectID    string             `json:"projectId" bson:"projectId"` // Unique URL key (e.g. "kiddostyle")
	Title        string             `json:"title" bson:"title"`
	Client       string             `json:"client" bson:"client"`
	Category     string             `json:"category" bson:"category"` // "platform" or "design"
	Description  string             `json:"description" bson:"description"`
	Outcome      string             `json:"outcome" bson:"outcome"`
	Tags         []string           `json:"tags" bson:"tags"`
	Image        string             `json:"image" bson:"image"` // Main S3 thumbnail image link
	Challenge    string             `json:"challenge" bson:"challenge"`
	Solution     string             `json:"solution" bson:"solution"`
	Architecture []string           `json:"architecture" bson:"architecture"`
	Metrics      []Metric           `json:"metrics" bson:"metrics"`
	SiteUrl       string             `json:"siteUrl" bson:"siteUrl"`
	AdminUrl      string             `json:"adminUrl" bson:"adminUrl"`
	DesktopMockup string             `json:"desktopMockup" bson:"desktopMockup"` // Dynamic MacBook device mockup
	TabletMockup  string             `json:"tabletMockup" bson:"tabletMockup"`   // Dynamic iPad device mockup
	MobileMockup  string             `json:"mobileMockup" bson:"mobileMockup"`   // Dynamic iPhone device mockup
	Screenshots   []Screenshot       `json:"screenshots" bson:"screenshots"`
}