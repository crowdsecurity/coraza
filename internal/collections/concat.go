// Copyright 2023 Juan Pablo Tosso and the OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package collections

import (
	"regexp"
	"strings"

	"github.com/corazawaf/coraza/v3/collection"
	"github.com/corazawaf/coraza/v3/internal/corazarules"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/corazawaf/coraza/v3/types/variables"
)

// ConcatCollection is a collection view over multiple collections.
type ConcatCollection struct {
	data     []collection.Collection
	variable variables.RuleVariable
	cache    concatCache
}

var _ collection.Collection = &ConcatCollection{}

func NewConcatCollection(variable variables.RuleVariable, data ...collection.Collection) *ConcatCollection {
	return &ConcatCollection{
		data:     data,
		variable: variable,
	}
}

// FindAll returns all matches for all collections
func (c *ConcatCollection) FindAll() []types.MatchData {
	parts := make([][]types.MatchData, len(c.data))
	for i, d := range c.data {
		parts[i] = d.FindAll()
	}
	return c.cache.concat(c.variable, parts)
}

// Name returns the name for the current CollectionconcatCollection
func (c *ConcatCollection) Name() string {
	return c.variable.Name()
}

// ConcatKeyed is a collection view over multiple keyed collections.
type ConcatKeyed struct {
	data     []collection.Keyed
	variable variables.RuleVariable
	cache    concatCache
}

var _ collection.Keyed = &ConcatKeyed{}

func NewConcatKeyed(variable variables.RuleVariable, data ...collection.Keyed) *ConcatKeyed {
	return &ConcatKeyed{
		data:     data,
		variable: variable,
	}
}

func (c *ConcatKeyed) Get(key string) []string {
	keyL := strings.ToLower(key)
	var res []string
	for _, c := range c.data {
		res = append(res, c.Get(keyL)...)
	}
	return res
}

// FindRegex returns a slice of MatchData for the regex
func (c *ConcatKeyed) FindRegex(key *regexp.Regexp) []types.MatchData {
	var res []types.MatchData
	for _, d := range c.data {
		res = append(res, replaceVariable(c.variable, d.FindRegex(key))...)
	}
	return res
}

// FindString returns a slice of MatchData for the string
func (c *ConcatKeyed) FindString(key string) []types.MatchData {
	if key == "" {
		// Children answer an empty key with their shared FindAll result,
		// which must not be relabelled in place below.
		return c.FindAll()
	}
	var res []types.MatchData
	for _, d := range c.data {
		res = append(res, replaceVariable(c.variable, d.FindString(key))...)
	}
	return res
}

// FindAll returns all matches for all collections
func (c *ConcatKeyed) FindAll() []types.MatchData {
	parts := make([][]types.MatchData, len(c.data))
	for i, d := range c.data {
		parts[i] = d.FindAll()
	}
	return c.cache.concat(c.variable, parts)
}

// Name returns the name for the current CollectionconcatCollection
func (c *ConcatKeyed) Name() string {
	return c.variable.Name()
}

// concatCache reuses the concatenation while every part is the very same
// slice as last time, which is the case when the parts come from unchanged
// caching collections. Holding the old parts keeps their memory alive, so a
// new slice can never reuse the same address.
type concatCache struct {
	parts [][]types.MatchData
	res   []types.MatchData
}

func sameSlice(a, b []types.MatchData) bool {
	return len(a) == len(b) && (len(a) == 0 || &a[0] == &b[0])
}

func (c *concatCache) concat(v variables.RuleVariable, parts [][]types.MatchData) []types.MatchData {
	if c.parts != nil && len(c.parts) == len(parts) {
		hit := true
		for i := range parts {
			if !sameSlice(c.parts[i], parts[i]) {
				hit = false
				break
			}
		}
		if hit {
			return c.res
		}
	}

	n := 0
	for _, p := range parts {
		n += len(p)
	}
	var res []types.MatchData
	if n > 0 {
		// Parts may be shared by their collection: copy rather than relabel them in place.
		buf := make([]corazarules.MatchData, n)
		res = make([]types.MatchData, n)
		i := 0
		for _, p := range parts {
			for _, m := range p {
				buf[i] = corazarules.MatchData{
					Variable_: v,
					Key_:      m.Key(),
					Value_:    m.Value(),
				}
				res[i] = &buf[i]
				i++
			}
		}
	}
	c.parts, c.res = parts, res
	return res
}

// replaceVariable ensures a returned match references the variable of a concatenated variable,
// not original one. It relabels md in place, so md must be a fresh slice owned by the caller,
// never a shared FindAll result.
func replaceVariable(v variables.RuleVariable, md []types.MatchData) []types.MatchData {
	for _, m := range md {
		m.(*corazarules.MatchData).Variable_ = v
	}
	return md
}
