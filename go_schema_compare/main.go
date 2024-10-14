package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var directory = "../auto_schema_compare/swagger_files/currentSwaggerDefinitions.json"

var objects = []string{
	"DataTable",
	"EmergencyGroup",
	"Flow",
	"Grammar",
	"GrammarLanguage",
	"IVR",
	"ScheduleGroup",
	"Schedule",
	"Prompt",
	"AuthzDivision",
	"DomainOrganizationRoleCreate",
	"InstagramIntegrationRequest",
	"MessagingSettingRequest",
	"SupportedContent",
	"ExternalMetricDefinitionCreateRequest",
	"ExternalContact",
	"FlowOutcome",
	"FlowMilestone",
	"GroupCreate",
	"RoleDivisionGrants",
	"CreateIntegrationRequest",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"CallableTimeSet",
	"RoutingQueue",
}

func main() {
	// Read the JSON file
	data, err := os.ReadFile(directory)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}

	// Unmarshal the JSON data
	var resourceData map[string]interface{}
	if err := json.Unmarshal(data, &resourceData); err != nil {
		log.Fatalf("failed to unmarshal JSON: %s", err)
	}

	// Create the directory if it doesn't exist
	if _, err := os.Stat("object"); os.IsNotExist(err) {
		if err := os.Mkdir("object", os.ModePerm); err != nil {
			log.Fatalf("failed to create directory: %s", err)
		}
	}
	
    for _, resource := range objects {
        // Check if the resource exists
        resourceObject, exists := resourceData[resource].(map[string]interface{})
        if !exists {
            log.Printf("resource %s not found in JSON", resource)
            continue
        }

        seenRefs := make(map[string]bool)
        fullyResolvedSchema := resolveRefs(resourceObject, resourceData, seenRefs)

        // Marshal the fully resolved schema back into JSON for output
        resourceJSON, err := json.MarshalIndent(fullyResolvedSchema, "", "  ")
        if err != nil {
            log.Fatalf("Failed to marshal resource JSON: %s", err)
        }

		// Write the resource JSON to a file
		fileName := fmt.Sprintf("object/%s.json", resource)
		if err := os.WriteFile(fileName, resourceJSON, os.ModePerm); err != nil {
			log.Fatalf("Failed to write resource JSON to file: %s", err)
		}
		fmt.Printf("Resource %s written to %s\n", resource, fileName)
    }
}

func resolveRefs(schema map[string]interface{}, definitions map[string]interface{}, seenRefs map[string]bool) map[string]interface{} {
	for key, value := range schema {
		switch typedValue := value.(type) {
		case map[string]interface{}:
			// Check for $ref in the object
			if ref, ok := typedValue["$ref"]; ok {
				refString, isString := ref.(string)
				if isString {
					// Resolve the reference
					refPath := strings.TrimPrefix(refString, "#/definitions/")
					// Check for cyclic reference
					if seenRefs[refPath] {
						log.Printf("Cyclic reference detected: %s", refPath)
						continue
					}
					seenRefs[refPath] = true

					// Find the referenced object in definitions
					if refDef, found := definitions[refPath]; found {
						// Resolve the referenced definition recursively
						resolvedDef := resolveRefs(refDef.(map[string]interface{}), definitions, seenRefs)
						// Replace the $ref with the resolved definition
						schema[key] = resolvedDef
					} else {
						log.Printf("Reference not found in definitions: %s", refPath)
					}
				}
			} else {
				// Continue resolving nested objects
				schema[key] = resolveRefs(typedValue, definitions, seenRefs)
			}
		case []interface{}:
			// If the value is an array, resolve each item
			for i, item := range typedValue {
				if itemMap, ok := item.(map[string]interface{}); ok {
					typedValue[i] = resolveRefs(itemMap, definitions, seenRefs)
				}
			}
		}
	}
	return schema
}