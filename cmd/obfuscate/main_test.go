package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestExtractHTMLNamesFindsMarkupAndInlineSelectors(t *testing.T) {
	content := `<div class="card dark" id="menu" data-track="yes"></div>
<style>.button:hover { color: red; } .md\:flex:hover { display: flex; } [data-state] { color: blue; }</style>
<script>document.querySelector('.action'); document.getElementById('dialog');</script>`

	names := extractHTMLNames(content)
	got := make(map[string]bool, len(names))
	for _, name := range names {
		got[name.Key()] = true
	}
	for _, want := range []string{
		"class:card", "class:dark", "id:menu", "data:track",
		"class:button", "class:md:flex", "data:state", "class:action", "id:dialog",
	} {
		if !got[want] {
			t.Errorf("extracted names do not contain %q; got %v", want, got)
		}
	}
	if got["class:button:hover"] {
		t.Error("pseudo-class suffix was included in the extracted class name")
	}
}

func TestCollectNamesProtectsCriticalNamesAndMappingIsDeterministic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.html")
	if err := os.WriteFile(path, []byte(`<div class="dark card" id="theme-toggle" data-contact-form=""></div>`), 0600); err != nil {
		t.Fatal(err)
	}

	names := collectNames([]string{path}, false)
	mapping := buildMapping(names)
	if _, ok := mapping["class:dark"]; ok {
		t.Fatal("protected dark class was included in the mapping")
	}
	if _, ok := mapping["id:theme-toggle"]; ok {
		t.Fatal("protected theme-toggle ID was included in the mapping")
	}
	if _, ok := mapping["data:contact-form"]; ok {
		t.Fatal("protected contact-form attribute was included in the mapping")
	}

	first := deterministicToken("class:card", "c")
	second := deterministicToken("class:card", "c")
	if first != second || !strings.HasPrefix(first, "c") {
		t.Fatalf("token generation is not deterministic or has the wrong prefix: %q, %q", first, second)
	}
}

func TestReplaceHTMLRewritesNamesAndReferences(t *testing.T) {
	mapping := map[string]string{
		"class:card": "c-card",
		"id:field":   "i-field",
		"id:help":    "i-help",
		"data:track": "d-track",
	}
	input := `<label for="field" class="card dark" data-track="yes"></label><input id="field" aria-describedby="help"><a href="#field">go</a><script>document.querySelector('.card');</script><style>.card:hover {}</style>`
	got := replaceHTML(input, mapping, &sync.Map{})

	for _, want := range []string{
		`for="i-field"`, `class="c-card dark"`, `data-d-track="yes"`,
		`id="i-field"`, `aria-describedby="i-help"`, `href="#i-field"`,
		`querySelector('.c-card')`, `.c-card:hover`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("rewritten HTML does not contain %q:\n%s", want, got)
		}
	}
}
