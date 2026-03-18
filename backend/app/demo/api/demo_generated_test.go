package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratedDemoApiCarriesPathID(t *testing.T) {
	checks := map[string][]string{
		filepath.Join("internal", "types", "types.go"): {
			"type GetItemReq struct",
			"Id int64",
			"type UpdateItemReq struct",
			"Id int64",
		},
		filepath.Join("internal", "handler", "getitemhandler.go"): {
			"req.Id",
		},
		filepath.Join("internal", "handler", "updateitemhandler.go"): {
			"req.Id",
		},
		filepath.Join("internal", "logic", "getitemlogic.go"): {
			"req.Id",
		},
		filepath.Join("internal", "logic", "updateitemlogic.go"): {
			"req.Id",
		},
	}

	for relPath, wantSnippets := range checks {
		content, err := os.ReadFile(relPath)
		if err != nil {
			t.Fatalf("read %s: %v", relPath, err)
		}

		text := string(content)
		for _, snippet := range wantSnippets {
			if !strings.Contains(text, snippet) {
				t.Fatalf("%s missing snippet %q", relPath, snippet)
			}
		}
	}
}
