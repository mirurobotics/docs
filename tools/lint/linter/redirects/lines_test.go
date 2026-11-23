package redirects

import (
	"reflect"
	"testing"
)

func TestLineForOffset(t *testing.T) {
	text := "abc\ndef\nghi"
	cases := []struct {
		off  int
		want int
	}{
		{0, 1},
		{2, 1},
		{3, 1},
		{4, 2},
		{7, 2},
		{8, 3},
		{11, 3},
	}
	for _, tc := range cases {
		got := lineForOffset(text, tc.off)
		if got != tc.want {
			t.Errorf("lineForOffset(%d): got %d, want %d", tc.off, got, tc.want)
		}
	}
}

func TestLineLookup(t *testing.T) {
	t.Run("count_zero", func(t *testing.T) {
		got := lineLookup("{}", 0)
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("decoder_error", func(t *testing.T) {
		got := lineLookup("not json", 1)
		want := []int{1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("redirects_token_then_eof", func(t *testing.T) {
		got := lineLookup(`{"redirects":`, 1)
		want := []int{1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("array_open_then_eof_mid_scan", func(t *testing.T) {
		got := lineLookup(`{"redirects": [`, 1)
		want := []int{1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("more_entries_than_count", func(t *testing.T) {
		text := "{\n" +
			"  \"redirects\": [\n" +
			"    {\"a\":1},\n" +
			"    {\"b\":2},\n" +
			"    {\"c\":3}\n" +
			"  ]\n" +
			"}\n"
		got := lineLookup(text, 2)
		if len(got) != 2 {
			t.Fatalf("got len %d, want 2", len(got))
		}
		if got[0] != 3 || got[1] != 4 {
			t.Errorf("got %v, want [3 4]", got)
		}
	})

	t.Run("nested_objects", func(t *testing.T) {
		text := "{\n" +
			"  \"redirects\": [\n" +
			"    {\"x\": {\"nested\": 1}}\n" +
			"  ]\n" +
			"}\n"
		got := lineLookup(text, 1)
		if len(got) != 1 || got[0] != 3 {
			t.Errorf("got %v, want [3]", got)
		}
	})
}

func TestTokenStart(t *testing.T) {
	t.Run("whitespace_only", func(t *testing.T) {
		got := tokenStart("   abc", 0)
		if got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})

	t.Run("comma_then_brace", func(t *testing.T) {
		got := tokenStart(",   {", 0)
		if got != 4 {
			t.Errorf("got %d, want 4", got)
		}
	})

	t.Run("offset_at_eof", func(t *testing.T) {
		got := tokenStart("abc", 3)
		if got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})

	t.Run("tab_and_cr", func(t *testing.T) {
		got := tokenStart("\t\r\n{", 0)
		if got != 3 {
			t.Errorf("got %d, want 3", got)
		}
	})
}
