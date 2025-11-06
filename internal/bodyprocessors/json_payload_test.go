package bodyprocessors_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/corazawaf/coraza/v3/experimental/plugins/plugintypes"
	"github.com/corazawaf/coraza/v3/internal/bodyprocessors"
	"github.com/corazawaf/coraza/v3/internal/corazawaf"
)

func getJSONBodyProcessor(t *testing.T) plugintypes.BodyProcessor {
	t.Helper()
	j, err := bodyprocessors.GetBodyProcessor("json")
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func TestJSONPayload(t *testing.T) {
	payload := strings.TrimSpace(`{"foo": "bar", "baz": 1, "qux": [1, 2, 3], "quux": {"corge": "grault"}}`)

	j := getJSONBodyProcessor(t)

	v := corazawaf.NewTransactionVariables()
	if err := j.ProcessRequest(strings.NewReader(payload), v, plugintypes.BodyProcessorOptions{
		Mime: "application/json",
	}); err != nil {
		t.Fatal(err)
	}

	if v.RawRequestBody().Get() != payload {
		t.Errorf("Expected %s, got %s", payload, v.RawRequestBody().Get())
	}
	if rbl, _ := strconv.Atoi(v.RawRequestBodyLength().Get()); rbl != len(payload) {
		t.Errorf("Expected %d, got %s", len(payload), v.RawRequestBodyLength().Get())
	}
	if v.ArgsPost().Get("json.foo")[0] != "bar" {
		t.Errorf("Expected %s, got %s", "bar", v.ArgsPost().Get("json.foo"))
	}
	if v.ArgsPost().Get("json.baz")[0] != "1" {
		t.Errorf("Expected %s, got %s", "1", v.ArgsPost().Get("json.baz"))
	}
	if v.ArgsPost().Get("json.qux.0")[0] != "1" {
		t.Errorf("Expected %s, got %s", "1", v.ArgsPost().Get("json.qux.0"))
	}
	if v.ArgsPost().Get("json.qux.1")[0] != "2" {
		t.Errorf("Expected %s, got %s", "2", v.ArgsPost().Get("json.qux.1"))
	}
	if v.ArgsPost().Get("json.qux.2")[0] != "3" {
		t.Errorf("Expected %s, got %s", "3", v.ArgsPost().Get("json.qux.2"))
	}
	if v.ArgsPost().Get("json.quux.corge")[0] != "grault" {
		t.Errorf("Expected %s, got %s", "grault", v.ArgsPost().Get("json.quux.corge"))
	}
	if len(v.ArgsPost().Get("json.quux")) != 0 {
		t.Errorf("Expected %s, got %s", "", v.ArgsPost().Get("json.quux"))
	}
}
