package bf

import (
	"bytes"
	"testing"
)

func TestNormalize(t *testing.T) {
	if !bytes.Equal(Normalize([]byte(`[`)), []byte(``)) {
		t.Fail()
	}
	if !bytes.Equal(Normalize([]byte(`]`)), []byte(``)) {
		t.Fail()
	}
	if !bytes.Equal(Normalize([]byte(`][`)), []byte(``)) {
		t.Fail()
	}
	if !bytes.Equal(Normalize([]byte(`[]`)), []byte(`[]`)) {
		t.Fail()
	}
}
