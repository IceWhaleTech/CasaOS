package service

import (
	"testing"
)

func TestSearch(t *testing.T) {
	if d, e := NewOtherService().Search("test"); e != nil || d == nil {

		t.Error("then test search error", e)
	}
}
