// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style.
// license that can be found in the LICENSE file.

// Code generated by generate-types. DO NOT EDIT.

package proto

import (
	"google.golang.org/protobuf/internal/encoding/wire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func sizeSingular(num wire.Number, kind protoreflect.Kind, v protoreflect.Value) int {
	switch kind {
	case protoreflect.BoolKind:
		return wire.SizeVarint(wire.EncodeBool(v.Bool()))
	case protoreflect.EnumKind:
		return wire.SizeVarint(uint64(v.Enum()))
	case protoreflect.Int32Kind:
		return wire.SizeVarint(uint64(int32(v.Int())))
	case protoreflect.Sint32Kind:
		return wire.SizeVarint(wire.EncodeZigZag(int64(int32(v.Int()))))
	case protoreflect.Uint32Kind:
		return wire.SizeVarint(uint64(uint32(v.Uint())))
	case protoreflect.Int64Kind:
		return wire.SizeVarint(uint64(v.Int()))
	case protoreflect.Sint64Kind:
		return wire.SizeVarint(wire.EncodeZigZag(v.Int()))
	case protoreflect.Uint64Kind:
		return wire.SizeVarint(v.Uint())
	case protoreflect.Sfixed32Kind:
		return wire.SizeFixed32()
	case protoreflect.Fixed32Kind:
		return wire.SizeFixed32()
	case protoreflect.FloatKind:
		return wire.SizeFixed32()
	case protoreflect.Sfixed64Kind:
		return wire.SizeFixed64()
	case protoreflect.Fixed64Kind:
		return wire.SizeFixed64()
	case protoreflect.DoubleKind:
		return wire.SizeFixed64()
	case protoreflect.StringKind:
		return wire.SizeBytes(len([]byte(v.String())))
	case protoreflect.BytesKind:
		return wire.SizeBytes(len(v.Bytes()))
	case protoreflect.MessageKind:
		return wire.SizeBytes(sizeMessage(v.Message()))
	case protoreflect.GroupKind:
		return wire.SizeGroup(num, sizeMessage(v.Message()))
	default:
		return 0
	}
}
