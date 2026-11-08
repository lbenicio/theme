package main

import (
	"strings"
	"testing"
)

func TestCategorizeConventionalCommits(t *testing.T) {
	commits := []commit{
		{hash: "0000001", message: "feat(ui): Add navigation"},
		{hash: "0000002", message: "feat: Add search"},
		{hash: "0000008", message: "feat!: Replace legacy API"},
		{hash: "0000003", message: "fix: Handle empty title"},
		{hash: "0000004", message: "security: Patch token handling"},
		{hash: "0000005", message: "chore: Refresh dependencies"},
		{hash: "0000006", message: "misc: Unclassified change"},
		{hash: "0000007", message: "Update contributor guide"},
	}

	categories := categorize(commits)
	wantNames := []string{"Added", "Fixed", "Security", "Chore", "Other"}
	if len(categories) != len(wantNames) {
		t.Fatalf("got %d categories, want %d: %#v", len(categories), len(wantNames), categories)
	}
	for i, want := range wantNames {
		if categories[i].name != want {
			t.Errorf("category %d = %q, want %q", i, categories[i].name, want)
		}
	}

	if got := categories[0].commits[0].message; got != "Add navigation" {
		t.Errorf("scoped feature message = %q, want stripped description", got)
	}
	if got := categories[0].commits[1].message; got != "Add search" {
		t.Errorf("feature commits are not sorted alphabetically: got second message %q", got)
	}
	if got := categories[0].commits[2].message; got != "Replace legacy API" {
		t.Errorf("breaking feature message = %q, want parsed description", got)
	}
	if got := categories[len(categories)-1].commits[0].message; got != "Unclassified change" {
		t.Errorf("unknown conventional message = %q, want description", got)
	}
	if got := categories[len(categories)-1].commits[1].message; got != "Update contributor guide" {
		t.Errorf("plain commit message = %q, want original message", got)
	}
}

func TestGenerateEntryIncludesVersionCategoryAndShortHash(t *testing.T) {
	entry := generateEntry("1.2.3", []category{{
		name:    "Fixed",
		commits: []commit{{hash: "abcdef0123456789", message: "Handle empty title"}},
	}})

	for _, want := range []string{"## [1.2.3] - ", "### Fixed", "- Handle empty title (abcdef0)"} {
		if !strings.Contains(entry, want) {
			t.Errorf("generated entry does not contain %q:\n%s", want, entry)
		}
	}
}
