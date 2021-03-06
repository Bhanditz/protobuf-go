// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"container/list"
	"reflect"

	"google.golang.org/protobuf/internal/encoding/wire"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

var bytesType = reflect.TypeOf([]byte(nil))

func makeLegacyUnknownFieldsFunc(t reflect.Type) func(p *messageDataType) pref.UnknownFields {
	fu, ok := t.FieldByName("XXX_unrecognized")
	if !ok || fu.Type != bytesType {
		return nil
	}
	fieldOffset := offsetOf(fu)
	return func(p *messageDataType) pref.UnknownFields {
		if p.p.IsNil() {
			return emptyUnknownFields{}
		}
		rv := p.p.Apply(fieldOffset).AsValueOf(bytesType)
		return (*legacyUnknownBytes)(rv.Interface().(*[]byte))
	}
}

// legacyUnknownBytes is a wrapper around XXX_unrecognized that implements
// the protoreflect.UnknownFields interface. This is challenging since we are
// limited to a []byte, so we do not have much flexibility in the choice
// of data structure that would have been ideal.
type legacyUnknownBytes []byte

func (fs *legacyUnknownBytes) Len() int {
	// Runtime complexity: O(n)
	b := *fs
	m := map[pref.FieldNumber]bool{}
	for len(b) > 0 {
		num, _, n := wire.ConsumeField(b)
		m[num] = true
		b = b[n:]
	}
	return len(m)
}

func (fs *legacyUnknownBytes) Get(num pref.FieldNumber) (raw pref.RawFields) {
	// Runtime complexity: O(n)
	b := *fs
	for len(b) > 0 {
		num2, _, n := wire.ConsumeField(b)
		if num == num2 {
			raw = append(raw, b[:n]...)
		}
		b = b[n:]
	}
	return raw
}

func (fs *legacyUnknownBytes) Set(num pref.FieldNumber, raw pref.RawFields) {
	num2, _, _ := wire.ConsumeTag(raw)
	if len(raw) > 0 && (!raw.IsValid() || num != num2) {
		panic("invalid raw fields")
	}

	// Remove all current fields of num.
	// Runtime complexity: O(n)
	b := *fs
	out := (*fs)[:0]
	for len(b) > 0 {
		num2, _, n := wire.ConsumeField(b)
		if num != num2 {
			out = append(out, b[:n]...)
		}
		b = b[n:]
	}
	*fs = out

	// Append new fields of num.
	*fs = append(*fs, raw...)
}

func (fs *legacyUnknownBytes) Range(f func(pref.FieldNumber, pref.RawFields) bool) {
	type entry struct {
		num pref.FieldNumber
		raw pref.RawFields
	}

	// Collect up a list of all the raw fields.
	// We preserve the order such that the latest encountered fields
	// are presented at the end.
	//
	// Runtime complexity: O(n)
	b := *fs
	l := list.New()
	m := map[pref.FieldNumber]*list.Element{}
	for len(b) > 0 {
		num, _, n := wire.ConsumeField(b)
		if e, ok := m[num]; ok {
			x := e.Value.(*entry)
			x.raw = append(x.raw, b[:n]...)
			l.MoveToBack(e)
		} else {
			x := &entry{num: num}
			x.raw = append(x.raw, b[:n]...)
			m[num] = l.PushBack(x)
		}
		b = b[n:]
	}

	// Iterate over all the raw fields.
	// This ranges over a snapshot of the current state such that mutations
	// while ranging are not observable.
	//
	// Runtime complexity: O(n)
	for e := l.Front(); e != nil; e = e.Next() {
		x := e.Value.(*entry)
		if !f(x.num, x.raw) {
			return
		}
	}
}

func (fs *legacyUnknownBytes) IsSupported() bool {
	return true
}
