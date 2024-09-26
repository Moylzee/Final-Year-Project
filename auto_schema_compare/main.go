package main

import (
	config "auto_schema_compare/config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	printBanner()
	log.Println("Running Auto Schema Compare")

	// Create the directory
	err := createDirectory()
	if err != nil {
		log.Fatalf("Error creating directories: %v", err)
	}
	
	// Retrieve the swagger file
	err = retrieveSwaggerFile(config.SwaggerUrl, config.NewSwaggerFilePath)
	if err != nil {
		log.Fatalf("Error retrieving Swagger file: %v", err)
	}

	// Write the Current Swagger definitions to a file
	writeSwaggerDefinitions()
}

func createDirectory() error {
	log.Println("Creating Directory")
	// Create the directory if it doesn't exist
	err := os.MkdirAll(config.DirName, 0755)
	if err != nil {
		return err
	}

	log.Println("Directory created")
	return nil
}

func retrieveSwaggerFile(url, filePath string) error {
	log.Println("Retrieving Swagger File")
	// Retrieve the Swagger file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error retrieving Swagger file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error retrieving Swagger file: %s", resp.Status)
	}

	log.Println("Retrieved Swagger file")

	err = writeSwaggerToFile(resp.Body, filePath)
	if err != nil {
		return fmt.Errorf("error writing Swagger file to disk: %v", err)
	}

	return nil
}

// Function to write the Entire Swagger to a file
func writeSwaggerToFile(body io.Reader, filePath string) error {
	log.Println("Writing Swagger To File")
	
	// Read the body of the response
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("error reading Swagger body: %v", err)
	}

	// Write the body to a file
	err = ioutil.WriteFile(filePath, bodyBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing Swagger file to disk: %v", err)
	}

	log.Printf("Swagger JSON written to %s\n", filePath)
	return nil
}

// Function to extract definitions and write them to currentSwaggerDefinitions.json
func writeSwaggerDefinitions() {
	// Read the Swagger file
	swaggerFile, err := ioutil.ReadFile(config.CurrentSwaggerFilePath)
	if err != nil {
		log.Fatalf("Error reading Swagger file: %v", err)
	}

	// Parse the JSON
	var swaggerData map[string]interface{}
	err = json.Unmarshal(swaggerFile, &swaggerData)
	if err != nil {
		log.Fatalf("Error parsing Swagger JSON: %v", err)
	}

	// Extract the "definitions" field
	definitions, ok := swaggerData["definitions"]
	if !ok {
		log.Fatalf("No 'definitions' field found in Swagger JSON")
	}

	// Convert definitions to JSON format
	definitionsJSON, err := json.MarshalIndent(definitions, "", "  ")
	if err != nil {
		log.Fatalf("Error converting definitions to JSON: %v", err)
	}

	// Write the definitions to a file
	definitionsFilePath := fmt.Sprintf("%s/current_swagger_definitions.json", config.DirName)
	err = ioutil.WriteFile(definitionsFilePath, definitionsJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing definitions to file: %v", err)
	}

	log.Printf("Swagger definitions written to %s\n", definitionsFilePath)
}

func printBanner() {
    banner := `
	      _
 ___    ___  | |__     ___   ________     ____ 
/ __|  / __| |  _ \   / _ \ |  _   _ \   / _  |
\__ \ | (__  | | | | |  __/ | | | | | | | (_| |
|___/  \___| |_| |_|  \___| |_| |_| |_|  \____| -Brian Moyles
`
    fmt.Println(banner)
}