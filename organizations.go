package main

import (
	"google.golang.org/api/admin/directory/v1"
	"encoding/json"
	"log"
)

type Organizations []*admin.UserOrganization

func parseOrganizations(o interface{}) Organizations {
	bytes, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("Could not encode organizations: %v", err)
	}

	var orgs Organizations

	err = json.Unmarshal(bytes, &orgs)
	if err != nil {
		log.Fatalf("Cound not decode organizations: %v", err)
	}

	return orgs
}

func (orgs Organizations) Primary() *admin.UserOrganization {
	for _, org := range orgs {
		if org.Primary {
			return org
		}
	}

	return nil
}
