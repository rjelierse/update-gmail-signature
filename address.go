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
	"log"

	"google.golang.org/api/admin/directory/v1"
)

type addresses []*admin.UserAddress

func parseAddresses(o interface{}) (addrs addresses) {
	data, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("Failed to marshall addresses: %v", err)
	}

	if err = json.Unmarshal(data, &addrs); err != nil {
		log.Fatalf("Failed to unmarshall addresses: %v", err)
	}

	return addrs
}

func (addrs addresses) Type(t string) *admin.UserAddress {
	for _, addr := range addrs {
		if addr.Type == t {
			return addr
		}
	}

	return nil
}
