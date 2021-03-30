package main

import (
	"encoding/xml"

	"github.com/sirupsen/logrus"

	"github.com/FlowingSPDG/vMix-profile-export/models"
)

// parseProfile parse vMix profile data
func parseProfile(prof []byte) (*models.Profile, error) {
	p := &models.Profile{}
	if err := xml.Unmarshal(prof, p); err != nil {
		logrus.Fatalln("Failed to unmarshal profile XML :", err)
		return nil, err
	}
	return p, nil
}
