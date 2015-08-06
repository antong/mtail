// Copyright 2011 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package vm

import (
	"regexp"
	"testing"
	"time"

	"github.com/google/mtail/metrics"
	"github.com/kylelemons/godebug/pretty"
)

var instructions = []struct {
	name          string
	i             instr
	re            []*regexp.Regexp
	str           []string
	reversedStack []interface{} // stack is inverted to be pushed onto vm stack

	expectedStack  []interface{}
	expectedThread thread
}{
	// Composite literals require too many explicit conversions.
	{"inc",
		instr{inc, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{0},
		[]interface{}{},
		thread{pc: 0, matches: map[int][]string{}},
	},
	{"inc by int",
		instr{inc, 2},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{0, 1}, // first is metric 0 "foo", second is the inc val.
		[]interface{}{},
		thread{pc: 0, matches: map[int][]string{}},
	},
	{"inc by string",
		instr{inc, 2},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{0, "1"}, // first is metric 0 "foo", second is the inc val.
		[]interface{}{},
		thread{pc: 0, matches: map[int][]string{}},
	},
	{"set int",
		instr{set, 2},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1, 2}, // set metric 1 "bar"
		[]interface{}{},
		thread{pc: 0, matches: map[int][]string{}},
	},
	{"set str",
		instr{set, 2},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1, "2"},
		[]interface{}{},
		thread{pc: 0, matches: map[int][]string{}},
	},
	{"match",
		instr{match, 0},
		[]*regexp.Regexp{regexp.MustCompile("a*b")},
		[]string{},
		[]interface{}{},
		[]interface{}{},
		thread{match: true, pc: 0, matches: map[int][]string{0: {"aaaab"}}},
	},
	{"cmp lt",
		instr{cmp, -1},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1, "2"},
		[]interface{}{},
		thread{pc: 0, match: true, matches: map[int][]string{}}},
	{"cmp eq",
		instr{cmp, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"2", "2"},
		[]interface{}{},
		thread{pc: 0, match: true, matches: map[int][]string{}}},
	{"cmp gt",
		instr{cmp, 1},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{},
		thread{pc: 0, match: true, matches: map[int][]string{}}},
	{"cmp le",
		instr{cmp, 1},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, "2"},
		[]interface{}{},
		thread{pc: 0, match: false, matches: map[int][]string{}}},
	{"cmp ne",
		instr{cmp, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"1", "2"},
		[]interface{}{},
		thread{pc: 0, match: false, matches: map[int][]string{}}},
	{"cmp ge",
		instr{cmp, -1},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 2},
		[]interface{}{},
		thread{pc: 0, match: false, matches: map[int][]string{}}},
	{"jnm",
		instr{jnm, 37},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{},
		thread{pc: 37, matches: map[int][]string{}}},
	{"jm",
		instr{jm, 37},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{},
		thread{pc: 0, matches: map[int][]string{}}},
	{"strptime",
		instr{strptime, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"2012/01/18 06:25:00", "2006/01/02 15:04:05"},
		[]interface{}{},
		thread{pc: 0, time: time.Date(2012, 1, 18, 6, 25, 0, 0, time.UTC),
			matches: map[int][]string{}}},
	{"add",
		instr{add, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(3)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"sub",
		instr{sub, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"mul",
		instr{sub, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"div",
		instr{sub, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{4, 2},
		[]interface{}{int64(2)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"tolower",
		instr{tolower, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"mIxeDCasE"},
		[]interface{}{"mixedcase"},
		thread{pc: 0, matches: map[int][]string{}}},
	{"length",
		instr{length, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"1234"},
		[]interface{}{4},
		thread{pc: 0, matches: map[int][]string{}}},
	{"length 0",
		instr{length, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{""},
		[]interface{}{0},
		thread{pc: 0, matches: map[int][]string{}}},
	{"shl",
		instr{shl, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(4)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"shr",
		instr{shr, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"and",
		instr{and, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(0)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"or",
		instr{or, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(3)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"xor",
		instr{xor, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(3)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"xor 2",
		instr{xor, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 3},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"xor 3",
		instr{xor, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{-1, 3},
		[]interface{}{int64(^3)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"not",
		instr{not, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{0},
		[]interface{}{int64(-1)},
		thread{pc: 0, matches: map[int][]string{}}},
	{"pow",
		instr{pow, 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 2},
		[]interface{}{int64(4)},
		thread{pc: 0, matches: map[int][]string{}}},
}

// TestInstrs tests that each instruction behaves as expected through one
// instruction cycle.
func TestInstrs(t *testing.T) {
	for _, tc := range instructions {
		var m []*metrics.Metric
		m = append(m,
			metrics.NewMetric("foo", "test", metrics.Counter),
			metrics.NewMetric("bar", "test", metrics.Counter))

		v := New(tc.name, tc.re, tc.str, m, []instr{tc.i}, true)
		v.t = new(thread)
		v.t.stack = make([]interface{}, 0)
		for _, item := range tc.reversedStack {
			v.t.Push(item)
		}
		v.t.matches = make(map[int][]string, 0)
		v.input = "aaaab"
		v.execute(v.t, tc.i)

		diff := pretty.Compare(tc.expectedStack, v.t.stack)
		if len(diff) > 0 {
			t.Errorf("%s: unexpected virtual machine stack state.\n%s", tc.name, diff)
		}
		// patch in the thread stack because otherwise the test table is huge
		tc.expectedThread.stack = tc.expectedStack

		if diff = pretty.Compare(v.t, &tc.expectedThread); len(diff) > 0 {
			t.Errorf("%s: unexpected virtual machine thread state.\n%s", tc.name, diff)
		}
	}
}
