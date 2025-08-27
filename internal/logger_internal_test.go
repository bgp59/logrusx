package logrusx_internal

import (
	"testing"
)

func testLogAddModuleDirPathPrefix(t *testing.T, mdpc *ModuleDirPathCache, prefix string, expectedPrefixList []string) {
	mdpc.addPrefix(prefix)
	if len(mdpc.prefixList) != len(expectedPrefixList) {
		t.Errorf("len(prefixList): want %d, got %d", len(expectedPrefixList), len(mdpc.prefixList))
	}
	for i, expected := range expectedPrefixList {
		if mdpc.prefixList[i] != expected {
			t.Errorf("prefixList[%d]: want %#v, got %#v", i, expected, mdpc.prefixList[i])
		}
	}
}

func testLogStripModuleDirPathPrefix(t *testing.T, mdpc *ModuleDirPathCache, filePath string, expected string) {
	result := mdpc.stripPrefix(filePath)
	if result != expected {
		t.Errorf("%#v: stripPrefix(%#v): want %#v, got %#v", mdpc, filePath, expected, result)
	}
}

func TestLogAddModuleDirPathPrefix(t *testing.T) {
	mdpc := &ModuleDirPathCache{}

	for _, tc := range []struct {
		prefix             string
		expectedPrefixList []string
	}{
		{"a/b", []string{"a/b"}},
		{"a/b/c", []string{"a/b/c", "a/b"}},
		{"a", []string{"a/b/c", "a/b", "a"}},
		{"a", []string{"a/b/c", "a/b", "a"}},
		{"a/b/c/d", []string{"a/b/c/d", "a/b/c", "a/b", "a"}},
		{"a/b", []string{"a/b/c/d", "a/b/c", "a/b", "a"}},
		{"b/b", []string{"a/b/c/d", "a/b/c", "b/b", "a/b", "a"}},
	} {
		testLogAddModuleDirPathPrefix(t, mdpc, tc.prefix, tc.expectedPrefixList)
	}
}

func TestStripPrefixMatch(t *testing.T) {
	mdpc := &ModuleDirPathCache{
		prefixList: []string{"a/b/c/", "c/d/", "e/"},
	}

	for _, tc := range []struct {
		filePath string
		expected string
	}{
		{"a/b/c/d/e/f", "d/e/f"},
		{"c/d/e/f/g", "e/f/g"},
		{"e/f/g/h", "f/g/h"},
	} {
		testLogStripModuleDirPathPrefix(t, mdpc, tc.filePath, tc.expected)
	}
}

func TestStripPrefixNoMatch(t *testing.T) {
	for _, tc := range []struct {
		keepNDirs int
		filePath  string
		expected  string
	}{
		{2, "a/b/c", "a/b/c"},
		{3, "x/y/c/d", "x/y/c/d"},
		{1, "x/y/z/e", "z/e"},
	} {
		testLogStripModuleDirPathPrefix(t, &ModuleDirPathCache{keepNDirs: tc.keepNDirs}, tc.filePath, tc.expected)
	}
}
