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
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/admin/directory/v1"
)

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

func getUser(subject string, config *jwt.Config) *admin.User {
	client, err := admin.New(config.Client(context.Background()))
	if err != nil {
		log.Fatal("Unable to create Admin SDK client:", err)
	}
	user, err := client.Users.Get(subject).Do()
	if err != nil {
		log.Fatal("Unable to get user ", subject, ":", err)
	}
	return user
}
