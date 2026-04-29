package ui

import (
	"testing"

	"github.com/Gu1llaum-3/sshm/internal/config"
)

func TestApplySourceFileFilter(t *testing.T) {
	hosts := []config.SSHHost{
		{Name: "a", SourceFile: "/home/u/.ssh/config"},
		{Name: "b", SourceFile: "/home/u/.ssh/work.conf"},
		{Name: "c", SourceFile: "/home/u/.ssh/work.conf"},
		{Name: "d", SourceFile: "/home/u/.ssh/perso.conf"},
	}

	t.Run("empty selected returns all", func(t *testing.T) {
		got := applySourceFileFilter(hosts, "")
		if len(got) != len(hosts) {
			t.Fatalf("expected %d hosts, got %d", len(hosts), len(got))
		}
	})

	t.Run("matching path filters to that file", func(t *testing.T) {
		got := applySourceFileFilter(hosts, "/home/u/.ssh/work.conf")
		if len(got) != 2 {
			t.Fatalf("expected 2 hosts, got %d", len(got))
		}
		for _, h := range got {
			if h.SourceFile != "/home/u/.ssh/work.conf" {
				t.Errorf("unexpected host %q with SourceFile %q", h.Name, h.SourceFile)
			}
		}
	})

	t.Run("non-matching path returns empty", func(t *testing.T) {
		got := applySourceFileFilter(hosts, "/nowhere.conf")
		if len(got) != 0 {
			t.Fatalf("expected 0 hosts, got %d", len(got))
		}
	})

	t.Run("empty input returns empty", func(t *testing.T) {
		got := applySourceFileFilter(nil, "/home/u/.ssh/work.conf")
		if len(got) != 0 {
			t.Fatalf("expected 0 hosts, got %d", len(got))
		}
	})
}

// TestFilterByTag verifies that filterByTag matches both own and inherited tags
func TestFilterByTag(t *testing.T) {
	hosts := []config.SSHHost{
		{Name: "h1", Tags: []string{"dev"}, InheritedTags: []string{}},
		{Name: "h2", Tags: []string{}, InheritedTags: []string{"prod"}},
		{Name: "h3", Tags: []string{"backup"}, InheritedTags: []string{"prod"}},
		{Name: "h4", Tags: []string{}, InheritedTags: []string{}},
	}

	tests := []struct {
		name    string
		tag     string
		want    []string // expected host names
	}{
		{
			name: "match own tag",
			tag:  "dev",
			want: []string{"h1"},
		},
		{
			name: "match inherited tag",
			tag:  "prod",
			want: []string{"h2", "h3"},
		},
		{
			name: "match with deduplication",
			tag:  "backup",
			want: []string{"h3"},
		},
		{
			name: "no match",
			tag:  "nonexistent",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterByTag(hosts, tt.tag)
			if len(got) != len(tt.want) {
				t.Errorf("filterByTag(%q) = %d hosts, want %d", tt.tag, len(got), len(tt.want))
				return
			}
			gotNames := make(map[string]bool)
			for _, h := range got {
				gotNames[h.Name] = true
			}
			for _, name := range tt.want {
				if !gotNames[name] {
					t.Errorf("filterByTag(%q): missing host %s", tt.tag, name)
				}
			}
		})
	}
}

// TestFilterByTag_InheritedOnly verifies filtering works with inherited tags alone
func TestFilterByTag_InheritedOnly(t *testing.T) {
	hosts := []config.SSHHost{
		{Name: "inherited-host", Tags: []string{}, InheritedTags: []string{"prod"}},
	}
	got := filterByTag(hosts, "prod")
	if len(got) != 1 || got[0].Name != "inherited-host" {
		t.Errorf("filterByTag(\"prod\") = %v, want [inherited-host]", got)
	}
}

// TestFilterByTag_Deduplication verifies that matching a deduplicated tag works
func TestFilterByTag_Deduplication(t *testing.T) {
	hosts := []config.SSHHost{
		{Name: "dup-host", Tags: []string{"prod"}, InheritedTags: []string{"prod"}},
	}
	got := filterByTag(hosts, "prod")
	if len(got) != 1 || got[0].Name != "dup-host" {
		t.Errorf("filterByTag(\"prod\") with duplicate = %v, want [dup-host]", got)
	}
}
