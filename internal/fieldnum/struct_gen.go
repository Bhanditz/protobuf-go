// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style.
// license that can be found in the LICENSE file.

// Code generated by generate-protos. DO NOT EDIT.

package fieldnum

// Field numbers for google.protobuf.Struct.
const (
	Struct_Fields = 1 // repeated google.protobuf.Struct.FieldsEntry
)

// Field numbers for google.protobuf.Struct.FieldsEntry.
const (
	Struct_FieldsEntry_Key   = 1 // optional string
	Struct_FieldsEntry_Value = 2 // optional google.protobuf.Value
)

// Field numbers for google.protobuf.Value.
const (
	Value_NullValue   = 1 // optional google.protobuf.NullValue
	Value_NumberValue = 2 // optional double
	Value_StringValue = 3 // optional string
	Value_BoolValue   = 4 // optional bool
	Value_StructValue = 5 // optional google.protobuf.Struct
	Value_ListValue   = 6 // optional google.protobuf.ListValue
)

// Field numbers for google.protobuf.ListValue.
const (
	ListValue_Values = 1 // repeated google.protobuf.Value
)
