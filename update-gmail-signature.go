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
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/gmail/v1"
	"html/template"
	"io/ioutil"
	"log"
)

type signatureFields struct {
	Name   string
	Title  string
	Mobile string
}

func getFields(user *admin.User) (fields signatureFields) {
	fields.Name = user.Name.FullName
	if org := parseOrganizations(user.Organizations).Primary(); org != nil {
		fields.Title = org.Title
	}
	if phone := parsePhoneNumbers(user.Phones).Type("mobile"); phone != nil {
		fields.Mobile = phone.Value
	}
	return
}

func getUsers(domain string, config *jwt.Config) []*admin.User {
	client, err := admin.New(config.Client(context.Background()))
	if err != nil {
		log.Fatal("Unable to create Admin SDK client:", err)
	}
	users, err := client.Users.List().Domain(domain).Do()
	if err != nil {
		log.Fatal("Unable to retrieve list of users:", err)
	}
	return users.Users
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

func main() {
	var credentialsPath string
	var templatePath string
	var domain string
	var subject string

	flag.StringVar(&credentialsPath, "secret", "client_secret.json", "The path to Google API credentials JSON")
	flag.StringVar(&templatePath, "template", "template.html", "The path to the signature template")
	flag.StringVar(&domain, "domain", "", "The organisational domain to retrieve the users from")
	flag.StringVar(&subject, "subject", "", "The person to impersonate when retrieving the users")
	flag.Parse()

	if len(domain) == 0 {
		log.Fatal("Please specify a domain")
	}

	credentials, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(credentials, gmail.GmailSettingsBasicScope, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	if len(subject) > 0 {
		config.Subject = subject
	}

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Unable to parse signature template: %v", err)
	}

	for _, user := range getUsers(domain, config) {
		// Don't attempt to update signature for users without GMail
		if user.IsMailboxSetup {
			setSignature(user, tpl, config)
		}
	}
}
