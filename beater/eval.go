package beater

import (
	"context"
	"encoding/json"

	"github.com/fupslot/fsstatbeat/config"
	"github.com/open-policy-agent/opa/rego"
)

func toMap(in interface{}) map[string]interface{} {
	var out map[string]interface{}
	b, _ := json.Marshal(in)
	json.Unmarshal(b, &out)
	return out
}

func StatEval(res config.Resource, in *FileState, ctx context.Context) (string, error) {
	out := map[string]interface{}{"file": toMap(in)}

	r := rego.New(
		rego.Query(res.Condition),
		rego.Input(out))

	rs, err := r.Eval(ctx)
	if err != nil {
		return "unknown", err
	}

	if rs[0].Expressions[0].Value.(bool) {
		return "pass", nil
	}

	return "fail", nil
}
