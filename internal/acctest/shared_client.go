package acctest

import (
	"log"
	"os"

	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/apiclient"
	"github.com/rootlyhq/terraform-provider-rootly/v5/meta"
	rootly "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

var (
	TestApiHost  string
	TestApiToken string

	SharedClient *rootly.ClientWithResponses
)

func init() {
	if v := os.Getenv("ROOTLY_API_URL"); v != "" {
		TestApiHost = v
	} else {
		TestApiHost = "https://api.rootly.com"
	}

	if v := os.Getenv("ROOTLY_API_TOKEN"); v != "" {
		TestApiToken = v
	}

	_, client, err := apiclient.New(TestApiHost, TestApiToken, meta.GetVersion())
	if err != nil {
		log.Fatalln(err.Error())
	}

	SharedClient = client
}
