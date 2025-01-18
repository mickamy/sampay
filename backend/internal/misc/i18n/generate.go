//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.Open("resources/ja.yaml")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var messages map[string]any
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&messages); err != nil {
		log.Fatalf("failed to decode YAML: %v", err)
	}

	output := &strings.Builder{}
	fmt.Fprintln(output, "package i18n")
	fmt.Fprintln(output, "// Code generated by generate.go; DO NOT EDIT.")
	fmt.Fprintln(output, "")
	fmt.Fprintln(output, "type MessageID string")
	fmt.Fprintln(output, "")
	fmt.Fprintln(output, "func (id MessageID) String() string { return string(id) }")
	fmt.Fprintln(output, "")
	fmt.Fprintln(output, "const (")

	makeConstants(messages, "", output)

	fmt.Fprintln(output, ")")

	if err := os.WriteFile("message_id.go", []byte(output.String()), 0644); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	fmt.Println("i18n: message_id.go generated")
}

func makeConstants(data map[string]any, prefix string, output *strings.Builder) {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		constName := makeConstantName(fullKey)

		switch v := data[key].(type) {
		case string:
			fmt.Fprintf(output, "\t%s MessageID = \"%s\"\n", constName, fullKey)
		case map[string]any:
			makeConstants(v, fullKey, output)
		default:
			log.Fatalf("unexpected type: %T", v)
		}
	}
}

func makeConstantName(key string) string {
	parts := strings.Split(key, ".")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}
