/*
Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may not use this file except in compliance with the License. A copy of the License is located at

	http://aws.amazon.com/apache2.0/

	or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

*/

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"

	"github.com/seregayoga/twitchstream/pkg/config"
	"github.com/seregayoga/twitchstream/pkg/handler"
	twitchapi "github.com/seregayoga/twitchstream/pkg/twitch"
)

func main() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})

	cfg := &config.Config{}
	err := envconfig.Process("ts", cfg)
	if err != nil {
		log.Fatal(err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  cfg.RedirectURL,
	}

	cookieStore := sessions.NewCookieStore([]byte(cfg.CookieSecret))
	twitchAPI := twitchapi.NewAPI(cfg.ClientID)
	handlers := handler.NewHandlers(oauth2Config, cookieStore, twitchAPI)

	handler.HandleFunc("/", handlers.HandleIndex)
	handler.HandleFunc("/login", handlers.HandleLogin)
	handler.HandleFunc("/redirect", handlers.HandleOAuth2Callback)
	handler.HandleFunc("/stream", handlers.HandleStream)
	handler.HandleFunc("/events", handlers.HandleEvents)

	adress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Println("Listening " + adress)
	log.Println(http.ListenAndServe(adress, nil))
}
