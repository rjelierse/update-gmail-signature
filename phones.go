package main

import (
	"google.golang.org/api/admin/directory/v1"
	"encoding/json"
	"log"
)

type PhoneNumbers []*admin.UserPhone

func parsePhoneNumbers(o interface{}) (numbers PhoneNumbers) {
	data, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("Failed to encode phone numbers: %v", err)
	}

	if err = json.Unmarshal(data, &numbers); err != nil {
		log.Fatalf("Failed to decode phone numbers: %v", err)
	}

	return numbers
}

func (nos PhoneNumbers) Type(t string) *admin.UserPhone {
	for _, no := range nos {
		if no.Type == t {
			return no
		}
	}

	return nil
}
