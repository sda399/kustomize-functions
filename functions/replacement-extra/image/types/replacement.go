// Copyright 2021 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

type TargetSelectorRegex struct {
	targetSelector *types.TargetSelector
	selectRegex    *types.SelectorRegex
	rejectRegex    []*types.SelectorRegex
}

func NewTargetSelectorRegex(ts *types.TargetSelector) (*TargetSelectorRegex, error) {
	tsr := new(TargetSelectorRegex)
	tsr.targetSelector = ts
	var err error

	tsr.selectRegex, err = types.NewSelectorRegex(ts.Select)
	if err != nil {
		return nil, err
	}

	rej := []*types.SelectorRegex{}
	for _, r := range ts.Reject {
		if r.IsEmpty() {
			continue
		}
		rr, err := types.NewSelectorRegex(r)
		if err != nil {
			return nil, err
		}
		rej = append(rej, rr)
	}
	tsr.rejectRegex = rej

	return tsr, nil
}

func (tsr *TargetSelectorRegex) Selects(id resid.ResId) bool {
	return tsr.selectRegex.MatchGvk(id.Gvk) && tsr.selectRegex.MatchName(id.Name) && tsr.selectRegex.MatchNamespace(id.Namespace)
}

func (tsr *TargetSelectorRegex) RejectsAny(ids []resid.ResId) bool {
	for _, r := range tsr.rejectRegex {
		for _, id := range ids {
			if r.MatchGvk(id.Gvk) && r.MatchName(id.Name) && r.MatchNamespace(id.Namespace) {
				return true
			}
		}
	}
	return false
}
