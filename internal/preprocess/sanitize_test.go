package preprocess_test

import (
	"hval/internal/preprocess"
	"testing"
)

func TestSanitize(t *testing.T) {
	input := "./../../test/values3.yaml"
	// output := "./tmp/values3.yaml"
	v, err := preprocess.NewValues(input, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := v.Open(); err != nil {
		t.Fatal(err)
	}

	if err := v.Sanitize(false); err != nil {
		t.Fatal(err)
	}

	t.FailNow()
}
