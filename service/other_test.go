package service

import (
	"testing"

	"go.uber.org/goleak"
)

func TestSearch(t *testing.T) {
	goleak.VerifyNone(t)

	if d, e := NewOtherService().Search("test"); e != nil || d == nil {
		t.Error("then test search error", e)
	}
}
