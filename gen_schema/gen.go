// Generate the scheme.go file in the root of this repo

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	binaryName  = "oapi-codegen"
	swaggerFile = "swagger.json"
	swaggerURL  = "https://rootly.com/swagger/v1/swagger.json"
)

func main() {
	getSwaggerConfig()
	genCode()
	cleanUp()
}

// Get the swagger config from rootly.com/api
func getSwaggerConfig() {
	// Making request
	log.Println("Getting latest swagger config from rootly.com/api")
	resp, err := http.Get(swaggerURL)
	if err != nil {
		log.Fatal("Failed to get latest swagger config from ", swaggerURL, "\n\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(
			"Failed to get latest swagger config from ",
			swaggerURL,
			"with exit code of ",
			resp.Status,
		)
	}
	log.Println("Got latest swagger data from rootly.com/api")

	// Decoding the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to parse data from ", swaggerURL, "\n\n", err)
	}

	// Writing the body to ./swagger.json
	err = ioutil.WriteFile(swaggerFile, body, 0666)
	if err != nil {
		log.Fatal("Failed to write to ", swaggerFile, "\n\n", err)
	}
	log.Println("Wrote swagger config to swagger.json")
}

// Gen the schema.gen.go file
func genCode() {
	log.Println("Generating schema.gen.go")

	// Getting the codegen tool path
	codegenPath, err := exec.LookPath(binaryName)
	if err != nil {
		log.Fatal("Is ", binaryName, " installed? ", "\n\n", err)
	}

	// Running the command and getting the stdout
	cmd := exec.Command(codegenPath, "-package", "rootly", swaggerFile)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed to run code gen", "\n\n", err)
	}

	// Writing the stdout to the schema.gen.go file
	err = ioutil.WriteFile("../schema/schema.gen.go", output, 0666)
	if err != nil {
		log.Fatal("Failed to write data to scheme file", "\n\n", err)
	}
	log.Println("Generated schema.gen.go")
}

// Remove the ./swagger.json file
func cleanUp() {
	err := os.Remove(swaggerFile)
	if err != nil {
		log.Fatal("Failed to delete ", swaggerFile, "\n\n", err)
	}
	log.Println("Removed swagger.json")
	log.Println(
		"You might need to remove duplites in the generated file (see https://github.com/deepmap/oapi-codegen/pull/200 for more info)",
	)
}
