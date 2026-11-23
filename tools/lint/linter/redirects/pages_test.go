package redirects

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestCleanPath(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"/x", "x"},
		{"/x/", "x"},
		{"/x?q=1", "x"},
		{"/x#frag", "x"},
		{"/x?q=1#frag", "x"},
		{"/", ""},
		{"x", "x"},
		{"/x/?q=1", "x"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := cleanPath(tc.in)
			if got != tc.want {
				t.Errorf("cleanPath(%q): got %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSplitWildcard(t *testing.T) {
	cases := []struct {
		name        string
		in          string
		wantPrefix  []string
		wantHasWild bool
	}{
		{"no_wildcard", "x/y", []string{"x", "y"}, false},
		{"trailing_wildcard", "x/:slug*", []string{"x"}, true},
		{"trailing_wildcard_no_star", "x/:slug", []string{"x"}, true},
		{"wildcard_first_segment", ":slug*", nil, true},
		{"empty", "", nil, false},
		{"only_wildcard", ":slug*", nil, true},
		{"empty_segments_skipped", "//x", []string{"x"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotPrefix, gotHas := splitWildcard(tc.in)
			if !reflect.DeepEqual(gotPrefix, tc.wantPrefix) {
				t.Errorf("prefix: got %v, want %v", gotPrefix, tc.wantPrefix)
			}
			if gotHas != tc.wantHasWild {
				t.Errorf("hasWildcard: got %v, want %v", gotHas, tc.wantHasWild)
			}
		})
	}
}

func TestDirHasPages(t *testing.T) {
	t.Run("missing_dir", func(t *testing.T) {
		root := t.TempDir()
		if dirHasPages(filepath.Join(root, "nope")) {
			t.Error("expected false for missing dir")
		}
	})

	t.Run("md_extension", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"a.md": "x"})
		if !dirHasPages(root) {
			t.Error("expected true for dir with .md")
		}
	})

	t.Run("non_page_files_only", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"a.txt": "x"})
		if dirHasPages(root) {
			t.Error("expected false for dir with .txt only")
		}
	})

	t.Run("subdir_no_pages", func(t *testing.T) {
		root := setupContentRoot(t, map[string]string{"0empty/": "", "z.mdx": "x"})
		if !dirHasPages(root) {
			t.Error("expected true for dir with .mdx sibling and empty subdir")
		}
	})
}
