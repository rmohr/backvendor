// Copyright (C) 2018 Tim Waugh
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package backvendor

import (
	"testing"
)

func TestNewFileHashes(t *testing.T) {
	hashes, err := NewFileHashes("git", "testdata/gosource", nil)
	if err != nil {
		t.Fatal(err)
	}
	if hashes == nil {
		t.Fatal("NewFileHashes returned nil map")
	}
	if hashes.root != "testdata/gosource" {
		t.Fatalf("Incorrect hashes.root (%s)", hashes.root)
	}
	if hashes.vcsCmd != "git" {
		t.Fatalf("Incorrect hashes.vcsCmd (%s)", hashes.vcsCmd)
	}
	emptyhash := FileHash("e69de29bb2d1d6434b8b29ae775ad8c2e48c5391")
	expected := map[string]FileHash{
		"ignored.go":                                 emptyhash,
		"vendor/github.com/foo/bar/bar.go":           emptyhash,
		"vendor/github.com/eggs/ham/ham.go":          emptyhash,
		"vendor/github.com/eggs/ham/spam/ignored.go": emptyhash,
	}
	if len(hashes.hashes) != len(expected) {
		t.Fatalf("len(hashes[%v]) != %d", hashes, len(expected))
	}
	for key, value := range expected {
		got, ok := hashes.hashes[key]
		if !ok {
			t.Errorf("%s missing", key)
			continue
		}
		if got != value {
			t.Errorf("%s: wrong hash (%s != %s)", key, got, value)
		}
	}
}

func TestNewFileHashesExclude(t *testing.T) {
	excludes := make(map[string]struct{})
	excludes["testdata/gosource/ignored.go"] = struct{}{}
	hashes, err := NewFileHashes("git", "testdata/gosource", excludes)
	if err != nil {
		t.Fatal(err)
	}
	emptyhash := FileHash("e69de29bb2d1d6434b8b29ae775ad8c2e48c5391")
	expected := map[string]FileHash{
		"vendor/github.com/foo/bar/bar.go":           emptyhash,
		"vendor/github.com/eggs/ham/ham.go":          emptyhash,
		"vendor/github.com/eggs/ham/spam/ignored.go": emptyhash,
	}
	if len(hashes.hashes) != len(expected) {
		t.Fatalf("len(hashes[%v]) != %d", hashes, len(expected))
	}
	for key, value := range expected {
		got, ok := hashes.hashes[key]
		if !ok {
			t.Errorf("%s missing", key)
			continue
		}
		if got != value {
			t.Errorf("%s: wrong hash (%s != %s)", key, got, value)
		}
	}
}

func TestIsSubsetOf(t *testing.T) {
	hashes, err := NewFileHashes("git", "testdata/gosource", nil)
	if err != nil {
		t.Fatal(err)
	}

	if !hashes.IsSubsetOf(hashes) {
		t.Fatalf("not subset of self")
	}

	other := &FileHashes{
		vcsCmd: hashes.vcsCmd,
		root:   hashes.root,
		hashes: make(map[string]FileHash),
	}
	for k, v := range hashes.hashes {
		other.hashes[k] = v
	}
	hashes.hashes["foo"] = FileHash("")
	if hashes.IsSubsetOf(other) {
		t.Fail()
	}
}
