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
	"flag"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/gmail/v1"
	"html/template"
	"io/ioutil"
	"log"
)

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
		log.Fatal("Unable to read client secret file:", err)
	}

	config, err := google.JWTConfigFromJSON(credentials, gmail.GmailSettingsBasicScope, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		log.Fatal("Unable to parse client secret file to config:", err)
	}

	if len(subject) > 0 {
		config.Subject = subject
	}

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("Unable to parse signature template:", err)
	}

	for _, user := range getUsers(domain, config) {
		// Don't attempt to update signature for users without GMail
		if user.IsMailboxSetup {
			setSignature(user, tpl, config)
		}
	}
}
