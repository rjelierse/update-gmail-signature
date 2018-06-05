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
	"bytes"
	"fmt"
	"html/template"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/admin/directory/v1"
	gmail "google.golang.org/api/gmail/v1"
)

type signatureFields struct {
	Name    string
	Title   string
	Mobile  string
	Phone   string
	Address string
}

func getFields(user *admin.User) (fields signatureFields) {
	fields.Name = user.Name.FullName
	if org := parseOrganizations(user.Organizations).Primary(); org != nil {
		fields.Title = org.Title
	}
	phones := parsePhoneNumbers(user.Phones)
	if phone := phones.Type("mobile"); phone != nil {
		fields.Mobile = phone.Value
	}
	if phone := phones.Type("work"); phone != nil {
		fields.Phone = phone.Value
	}
	if address := parseAddresses(user.Addresses).Type("work"); address != nil {
		fields.Address = address.Formatted
	}
	return
}

func setSignature(user *admin.User, tpl *template.Template, config *jwt.Config) {
	fmt.Printf("Updating signature for user %v\n", user.PrimaryEmail)
	var buf bytes.Buffer

	v := getFields(user)

	if err := tpl.Execute(&buf, v); err != nil {
		log.Fatal("Cannot execute template:", err)
	}

	result := buf.String()

	config.Subject = user.PrimaryEmail
	client, err := gmail.New(config.Client(context.Background()))
	if err != nil {
		log.Fatal("Unable to create GMail client:", err)
	}

	_, err = client.Users.Settings.SendAs.Patch(user.PrimaryEmail, user.PrimaryEmail, &gmail.SendAs{Signature: result}).Do()
	if err != nil {
		log.Fatal("Unable to set signature:", err)
	}
}
