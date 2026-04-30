package v8_1

import (
	"encoding/json"
	"fmt"
	"strings"

	upgradetypes "cosmossdk.io/x/upgrade/types"
)

type SilentSlashRecord struct {
	Operator string `json:"operator"` // bech32 valoper
	Height   int64  `json:"height"`
	Fraction string `json:"fraction"` // math.LegacyDec string
	Reason   string `json:"reason"`   // "missing_signature" or "double_sign"
}

var SilentSlashes = []SilentSlashRecord{
	// TODO(v8.1): populate before release
}

func resolveSilentSlashes(plan upgradetypes.Plan) ([]SilentSlashRecord, error) {
	combined := append([]SilentSlashRecord{}, SilentSlashes...)

	if plan.Info != "" {
		var info struct {
			SilentSlashes []SilentSlashRecord `json:"silent_slashes"`
		}
		if err := json.Unmarshal([]byte(plan.Info), &info); err != nil {
			if strings.Contains(plan.Info, "silent_slashes") {
				return nil, fmt.Errorf("v8.1.0: malformed silent_slashes in plan.Info: %w", err)
			}
		} else {
			combined = append(combined, info.SilentSlashes...)
		}
	}

	if len(combined) == 0 {
		return nil, nil
	}

	// Dedupe by operator, last entry wins, preserve insertion order.
	seen := make(map[string]int, len(combined))
	deduped := make([]SilentSlashRecord, 0, len(combined))
	for _, r := range combined {
		if i, ok := seen[r.Operator]; ok {
			deduped[i] = r
			continue
		}
		seen[r.Operator] = len(deduped)
		deduped = append(deduped, r)
	}
	return deduped, nil
}

func affectedValidators(records []SilentSlashRecord) map[string]struct{} {
	out := make(map[string]struct{}, len(records))
	for _, s := range records {
		out[s.Operator] = struct{}{}
	}
	return out
}
