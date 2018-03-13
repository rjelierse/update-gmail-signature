// Copyright 2017 Raymond Jelierse
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"google.golang.org/api/admin/directory/v1"
	"log"
)

type phoneNumbers []*admin.UserPhone

func parsePhoneNumbers(o interface{}) (numbers phoneNumbers) {
	data, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("Failed to encode phone numbers: %v", err)
	}

	if err = json.Unmarshal(data, &numbers); err != nil {
		log.Fatalf("Failed to decode phone numbers: %v", err)
	}

	return numbers
}

func (nos phoneNumbers) Type(t string) *admin.UserPhone {
	for _, no := range nos {
		if no.Type == t {
			return no
		}
	}

	return nil
}
