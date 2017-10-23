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
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"fmt"
	"html/template"
	"bytes"
	"flag"
)

var gAdmin *admin.Service
var gMail *gmail.Service

type signatureFields struct {
	Name string
	Title string
	Mobile string
}

func getUsers() []*admin.User {
	users, err := gAdmin.Users.List().Customer("my_customer").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of users: %v", err)
	}

	return users.Users
}

func setSignature(user *admin.User, tpl *template.Template) {
	fmt.Printf("Updating signature for user %v\n", user.PrimaryEmail)
	var buf bytes.Buffer

	v := signatureFields{
		Name:  user.Name.FullName,
		Title: parseOrganizations(user.Organizations).Primary().Title,
		Mobile: parsePhoneNumbers(user.Phones).Type("mobile").Value,
	}

	if err := tpl.Execute(&buf, v); err != nil {
		log.Fatalf("Cannot execute template: %v", err)
	}

	result := buf.String()

	_, err := gMail.Users.Settings.SendAs.Patch(user.PrimaryEmail, user.PrimaryEmail, &gmail.SendAs{Signature: result}).Do()
	if err != nil {
		log.Fatalf("Unable to set signature for %v: %v", user.PrimaryEmail, err)
	}
}

func main() {
	var credentialsPath string
	var templatePath string

	flag.StringVar(&credentialsPath, "secret", "client_secret.json", "The path to Google API credentials JSON")
	flag.StringVar(&templatePath, "template", "template.html", "The path to the signature template")
	flag.Parse()

	credentials, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, gmail.GmailSettingsBasicScope, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(context.Background(), config)

	gAdmin, err = admin.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve G Suite Admin client: %v", err)
	}

	gMail, err = gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve GMail client: %v", err)
	}

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Unable to parse signature template: %v", err)
	}

	for _, user := range getUsers() {
		// Don't attempt to update signature for users without GMail
		if user.IsMailboxSetup {
			setSignature(user, tpl)
		}
	}
}
