// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileinit

import (
	"bytes"
	"fmt"
	"reflect"

	defval "google.golang.org/protobuf/internal/encoding/defval"
	wire "google.golang.org/protobuf/internal/encoding/wire"
	fieldnum "google.golang.org/protobuf/internal/fieldnum"
	pimpl "google.golang.org/protobuf/internal/impl"
	ptype "google.golang.org/protobuf/internal/prototype"
	pvalue "google.golang.org/protobuf/internal/value"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	preg "google.golang.org/protobuf/reflect/protoregistry"
)

func (file *fileDesc) lazyInit() *fileLazy {
	file.once.Do(func() {
		file.unmarshalFull(file.rawDesc)
		file.resolveImports()
		file.resolveEnums()
		file.resolveMessages()
		file.resolveExtensions()
		file.resolveServices()
		file.finishInit()
	})
	return file.lazy
}

func (file *fileDesc) resolveImports() {
	// TODO: Resolve file dependencies.
}

func (file *fileDesc) resolveEnums() {
	enumDecls := file.GoTypes[:len(file.allEnums)]
	for i := range file.allEnums {
		ed := &file.allEnums[i]

		// Associate the EnumType with a concrete Go type.
		enumCache := map[pref.EnumNumber]pref.Enum{}
		ed.lazy.typ = reflect.TypeOf(enumDecls[i])
		ed.lazy.new = func(n pref.EnumNumber) pref.Enum {
			if v, ok := enumCache[n]; ok {
				return v
			}
			v := reflect.New(ed.lazy.typ).Elem()
			v.SetInt(int64(n))
			return v.Interface().(pref.Enum)
		}
		for i := range ed.lazy.values.list {
			n := ed.lazy.values.list[i].number
			enumCache[n] = ed.lazy.new(n)
		}
	}
}

func (file *fileDesc) resolveMessages() {
	messageDecls := file.GoTypes[len(file.allEnums):]
	for i := range file.allMessages {
		md := &file.allMessages[i]

		// Associate the MessageType with a concrete Go type.
		if !md.isMapEntry {
			md.lazy.typ = reflect.TypeOf(messageDecls[i])
			md.lazy.new = func() pref.Message {
				t := md.lazy.typ.Elem()
				return reflect.New(t).Interface().(pref.ProtoMessage).ProtoReflect()
			}
		}

		// Resolve message field dependencies.
		for j := range md.lazy.fields.list {
			fd := &md.lazy.fields.list[j]
			if fd.isWeak {
				continue
			}

			switch fd.kind {
			case pref.EnumKind:
				fd.enumType = file.popEnumDependency()
			case pref.MessageKind, pref.GroupKind:
				fd.messageType = file.popMessageDependency()
			}
			fd.isMap = file.isMapEntry(fd.messageType)
			if !fd.hasPacked && file.lazy.syntax != pref.Proto2 && fd.cardinality == pref.Repeated {
				switch fd.kind {
				case pref.StringKind, pref.BytesKind, pref.MessageKind, pref.GroupKind:
					fd.isPacked = false
				default:
					fd.isPacked = true
				}
			}
			fd.defVal.lazyInit(fd.kind, file.enumValuesOf(fd.enumType))
		}
	}
}

func (file *fileDesc) resolveExtensions() {
	for i := range file.allExtensions {
		xd := &file.allExtensions[i]

		// Associate the ExtensionType with a concrete Go type.
		var typ reflect.Type
		switch xd.lazy.kind {
		case pref.EnumKind, pref.MessageKind, pref.GroupKind:
			typ = reflect.TypeOf(file.GoTypes[file.DependencyIndexes[0]])
		default:
			typ = goTypeForPBKind[xd.lazy.kind]
		}
		switch xd.lazy.cardinality {
		case pref.Optional:
			switch xd.lazy.kind {
			case pref.EnumKind:
				et := pimpl.Export{}.EnumTypeOf(reflect.Zero(typ).Interface())
				xd.lazy.typ = typ
				xd.lazy.new = func() pref.Value {
					return xd.lazy.defVal.get()
				}
				xd.lazy.valueOf = func(v interface{}) pref.Value {
					ev := v.(pref.Enum)
					return pref.ValueOf(ev.Number())
				}
				xd.lazy.interfaceOf = func(pv pref.Value) interface{} {
					return et.New(pv.Enum())
				}
			case pref.MessageKind, pref.GroupKind:
				mt := pimpl.Export{}.MessageTypeOf(reflect.Zero(typ).Interface())
				xd.lazy.typ = typ
				xd.lazy.new = func() pref.Value {
					return pref.ValueOf(mt.New())
				}
				xd.lazy.valueOf = func(v interface{}) pref.Value {
					mv := v.(pref.ProtoMessage).ProtoReflect()
					return pref.ValueOf(mv)
				}
				xd.lazy.interfaceOf = func(pv pref.Value) interface{} {
					return pv.Message().Interface()
				}
			default:
				xd.lazy.typ = goTypeForPBKind[xd.lazy.kind]
				xd.lazy.new = func() pref.Value {
					return xd.lazy.defVal.get()
				}
				xd.lazy.valueOf = func(v interface{}) pref.Value {
					return pref.ValueOf(v)
				}
				xd.lazy.interfaceOf = func(pv pref.Value) interface{} {
					return pv.Interface()
				}
			}
		case pref.Repeated:
			c := pvalue.NewConverter(typ, xd.lazy.kind)
			xd.lazy.typ = reflect.PtrTo(reflect.SliceOf(typ))
			xd.lazy.new = func() pref.Value {
				v := reflect.New(xd.lazy.typ.Elem()).Interface()
				return pref.ValueOf(pvalue.ListOf(v, c))
			}
			xd.lazy.valueOf = func(v interface{}) pref.Value {
				return pref.ValueOf(pvalue.ListOf(v, c))
			}
			xd.lazy.interfaceOf = func(pv pref.Value) interface{} {
				return pv.List().(pvalue.Unwrapper).ProtoUnwrap()
			}
		default:
			panic(fmt.Sprintf("invalid cardinality: %v", xd.lazy.cardinality))
		}

		// Resolve extension field dependency.
		switch xd.lazy.kind {
		case pref.EnumKind:
			xd.lazy.enumType = file.popEnumDependency()
		case pref.MessageKind, pref.GroupKind:
			xd.lazy.messageType = file.popMessageDependency()
		}
		xd.lazy.defVal.lazyInit(xd.lazy.kind, file.enumValuesOf(xd.lazy.enumType))
	}
}

var goTypeForPBKind = map[pref.Kind]reflect.Type{
	pref.BoolKind:     reflect.TypeOf(bool(false)),
	pref.Int32Kind:    reflect.TypeOf(int32(0)),
	pref.Sint32Kind:   reflect.TypeOf(int32(0)),
	pref.Sfixed32Kind: reflect.TypeOf(int32(0)),
	pref.Int64Kind:    reflect.TypeOf(int64(0)),
	pref.Sint64Kind:   reflect.TypeOf(int64(0)),
	pref.Sfixed64Kind: reflect.TypeOf(int64(0)),
	pref.Uint32Kind:   reflect.TypeOf(uint32(0)),
	pref.Fixed32Kind:  reflect.TypeOf(uint32(0)),
	pref.Uint64Kind:   reflect.TypeOf(uint64(0)),
	pref.Fixed64Kind:  reflect.TypeOf(uint64(0)),
	pref.FloatKind:    reflect.TypeOf(float32(0)),
	pref.DoubleKind:   reflect.TypeOf(float64(0)),
	pref.StringKind:   reflect.TypeOf(string("")),
	pref.BytesKind:    reflect.TypeOf([]byte(nil)),
}

func (file *fileDesc) resolveServices() {
	for i := range file.services.list {
		sd := &file.services.list[i]

		// Resolve method dependencies.
		for j := range sd.lazy.methods.list {
			md := &sd.lazy.methods.list[j]
			md.inputType = file.popMessageDependency()
			md.outputType = file.popMessageDependency()
		}
	}
}

// isMapEntry reports whether the message is a map entry, being careful to
// avoid calling the IsMapEntry method if the message is declared
// within the same file (which would cause a recursive init deadlock).
func (fd *fileDesc) isMapEntry(md pref.MessageDescriptor) bool {
	if md == nil {
		return false
	}
	if md, ok := md.(*messageDescriptor); ok && md.parentFile == fd {
		return md.isMapEntry
	}
	return md.IsMapEntry()
}

// enumValuesOf retrieves the list of enum values for the given enum,
// being careful to avoid calling the Values method if the enum is declared
// within the same file (which would cause a recursive init deadlock).
func (fd *fileDesc) enumValuesOf(ed pref.EnumDescriptor) pref.EnumValueDescriptors {
	if ed == nil {
		return nil
	}
	if ed, ok := ed.(*enumDesc); ok && ed.parentFile == fd {
		return &ed.lazy.values
	}
	return ed.Values()
}

func (fd *fileDesc) popEnumDependency() pref.EnumDescriptor {
	depIdx := fd.popDependencyIndex()
	if depIdx < len(fd.allEnums)+len(fd.allMessages) {
		return &fd.allEnums[depIdx]
	} else {
		return pimpl.Export{}.EnumDescriptorOf(fd.GoTypes[depIdx])
	}
}

func (fd *fileDesc) popMessageDependency() pref.MessageDescriptor {
	depIdx := fd.popDependencyIndex()
	if depIdx < len(fd.allEnums)+len(fd.allMessages) {
		return fd.allMessages[depIdx-len(fd.allEnums)].asDesc()
	} else {
		return pimpl.Export{}.MessageDescriptorOf(fd.GoTypes[depIdx])
	}
}

func (fi *fileInit) popDependencyIndex() int {
	depIdx := fi.DependencyIndexes[0]
	fi.DependencyIndexes = fi.DependencyIndexes[1:]
	return int(depIdx)
}

func (fi *fileInit) finishInit() {
	if len(fi.DependencyIndexes) > 0 {
		panic("unused dependencies")
	}
	*fi = fileInit{} // clear fileInit for GC to reclaim resources
}

type defaultValue struct {
	has   bool
	val   pref.Value
	enum  pref.EnumValueDescriptor
	check func() // only set for non-empty bytes
}

func (dv *defaultValue) get() pref.Value {
	if dv.check != nil {
		dv.check()
	}
	return dv.val
}

func (dv *defaultValue) lazyInit(k pref.Kind, eds pref.EnumValueDescriptors) {
	if dv.has {
		switch k {
		case pref.EnumKind:
			// File descriptors always store default enums by name.
			dv.enum = eds.ByName(pref.Name(dv.val.String()))
			dv.val = pref.ValueOf(dv.enum.Number())
		case pref.BytesKind:
			// Store a copy of the default bytes, so that we can detect
			// accidental mutations of the original value.
			b := append([]byte(nil), dv.val.Bytes()...)
			dv.check = func() {
				if !bytes.Equal(b, dv.val.Bytes()) {
					// TODO: Avoid panic if we're running with the race detector
					// and instead spawn a goroutine that periodically resets
					// this value back to the original to induce a race.
					panic("detected mutation on the default bytes")
				}
			}
		}
	} else {
		switch k {
		case pref.BoolKind:
			dv.val = pref.ValueOf(false)
		case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
			dv.val = pref.ValueOf(int32(0))
		case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
			dv.val = pref.ValueOf(int64(0))
		case pref.Uint32Kind, pref.Fixed32Kind:
			dv.val = pref.ValueOf(uint32(0))
		case pref.Uint64Kind, pref.Fixed64Kind:
			dv.val = pref.ValueOf(uint64(0))
		case pref.FloatKind:
			dv.val = pref.ValueOf(float32(0))
		case pref.DoubleKind:
			dv.val = pref.ValueOf(float64(0))
		case pref.StringKind:
			dv.val = pref.ValueOf(string(""))
		case pref.BytesKind:
			dv.val = pref.ValueOf([]byte(nil))
		case pref.EnumKind:
			dv.enum = eds.Get(0)
			dv.val = pref.ValueOf(dv.enum.Number())
		}
	}
}

func (fd *fileDesc) unmarshalFull(b []byte) {
	nb := getNameBuilder()
	defer putNameBuilder(nb)

	var hasSyntax bool
	var enumIdx, messageIdx, extensionIdx, serviceIdx int
	fd.lazy = &fileLazy{}
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.FileDescriptorProto_PublicDependency:
				fd.lazy.imports[v].IsPublic = true
			case fieldnum.FileDescriptorProto_WeakDependency:
				fd.lazy.imports[v].IsWeak = true
			}
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.FileDescriptorProto_Syntax:
				hasSyntax = true
				switch string(v) {
				case "proto2":
					fd.lazy.syntax = pref.Proto2
				case "proto3":
					fd.lazy.syntax = pref.Proto3
				default:
					panic("invalid syntax")
				}
			case fieldnum.FileDescriptorProto_Dependency:
				fd.lazy.imports = append(fd.lazy.imports, pref.FileImport{
					FileDescriptor: ptype.PlaceholderFile(nb.MakeString(v), ""),
				})
			case fieldnum.FileDescriptorProto_EnumType:
				fd.enums.list[enumIdx].unmarshalFull(v, nb)
				enumIdx++
			case fieldnum.FileDescriptorProto_MessageType:
				fd.messages.list[messageIdx].unmarshalFull(v, nb)
				messageIdx++
			case fieldnum.FileDescriptorProto_Extension:
				fd.extensions.list[extensionIdx].unmarshalFull(v, nb)
				extensionIdx++
			case fieldnum.FileDescriptorProto_Service:
				fd.services.list[serviceIdx].unmarshalFull(v, nb)
				serviceIdx++
			case fieldnum.FileDescriptorProto_Options:
				fd.lazy.options = append(fd.lazy.options, v...)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}

	// If syntax is missing, it is assumed to be proto2.
	if !hasSyntax {
		fd.lazy.syntax = pref.Proto2
	}
}

func (ed *enumDesc) unmarshalFull(b []byte, nb *nameBuilder) {
	var rawValues [][]byte
	ed.lazy = new(enumLazy)
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.EnumDescriptorProto_Value:
				rawValues = append(rawValues, v)
			case fieldnum.EnumDescriptorProto_ReservedName:
				ed.lazy.resvNames.list = append(ed.lazy.resvNames.list, pref.Name(nb.MakeString(v)))
			case fieldnum.EnumDescriptorProto_ReservedRange:
				ed.lazy.resvRanges.list = append(ed.lazy.resvRanges.list, unmarshalEnumReservedRange(v))
			case fieldnum.EnumDescriptorProto_Options:
				ed.lazy.options = append(ed.lazy.options, v...)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}

	if len(rawValues) > 0 {
		ed.lazy.values.list = make([]enumValueDesc, len(rawValues))
		for i, b := range rawValues {
			ed.lazy.values.list[i].unmarshalFull(b, nb, ed.parentFile, ed, i)
		}
	}
}

func unmarshalEnumReservedRange(b []byte) (r [2]pref.EnumNumber) {
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.EnumDescriptorProto_EnumReservedRange_Start:
				r[0] = pref.EnumNumber(v)
			case fieldnum.EnumDescriptorProto_EnumReservedRange_End:
				r[1] = pref.EnumNumber(v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
	return r
}

func (vd *enumValueDesc) unmarshalFull(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	vd.parentFile = pf
	vd.parent = pd
	vd.index = i

	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.EnumValueDescriptorProto_Number:
				vd.number = pref.EnumNumber(v)
			}
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.EnumValueDescriptorProto_Name:
				vd.fullName = nb.AppendFullName(pd.FullName(), v)
			case fieldnum.EnumValueDescriptorProto_Options:
				vd.options = append(vd.options, v...)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}

func (md *messageDesc) unmarshalFull(b []byte, nb *nameBuilder) {
	var rawFields, rawOneofs [][]byte
	var enumIdx, messageIdx, extensionIdx int
	var isMapEntry bool
	md.lazy = new(messageLazy)
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.DescriptorProto_Field:
				rawFields = append(rawFields, v)
			case fieldnum.DescriptorProto_OneofDecl:
				rawOneofs = append(rawOneofs, v)
			case fieldnum.DescriptorProto_ReservedName:
				md.lazy.resvNames.list = append(md.lazy.resvNames.list, pref.Name(nb.MakeString(v)))
			case fieldnum.DescriptorProto_ReservedRange:
				md.lazy.resvRanges.list = append(md.lazy.resvRanges.list, unmarshalMessageReservedRange(v))
			case fieldnum.DescriptorProto_ExtensionRange:
				r, opts := unmarshalMessageExtensionRange(v)
				md.lazy.extRanges.list = append(md.lazy.extRanges.list, r)
				md.lazy.extRangeOptions = append(md.lazy.extRangeOptions, opts)
			case fieldnum.DescriptorProto_EnumType:
				md.enums.list[enumIdx].unmarshalFull(v, nb)
				enumIdx++
			case fieldnum.DescriptorProto_NestedType:
				md.messages.list[messageIdx].unmarshalFull(v, nb)
				messageIdx++
			case fieldnum.DescriptorProto_Extension:
				md.extensions.list[extensionIdx].unmarshalFull(v, nb)
				extensionIdx++
			case fieldnum.DescriptorProto_Options:
				md.unmarshalOptions(v, &isMapEntry)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}

	if len(rawFields) > 0 || len(rawOneofs) > 0 {
		md.lazy.fields.list = make([]fieldDesc, len(rawFields))
		md.lazy.oneofs.list = make([]oneofDesc, len(rawOneofs))
		for i, b := range rawFields {
			fd := &md.lazy.fields.list[i]
			fd.unmarshalFull(b, nb, md.parentFile, md.asDesc(), i)
			if fd.cardinality == pref.Required {
				md.lazy.reqNumbers.list = append(md.lazy.reqNumbers.list, fd.number)
			}
		}
		for i, b := range rawOneofs {
			od := &md.lazy.oneofs.list[i]
			od.unmarshalFull(b, nb, md.parentFile, md.asDesc(), i)
		}
	}

	if isMapEntry != md.isMapEntry {
		panic("mismatching map entry property")
	}
}

func (md *messageDesc) unmarshalOptions(b []byte, isMapEntry *bool) {
	md.lazy.options = append(md.lazy.options, b...)
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.MessageOptions_MapEntry:
				*isMapEntry = wire.DecodeBool(v)
			case fieldnum.MessageOptions_MessageSetWireFormat:
				md.lazy.isMessageSet = wire.DecodeBool(v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}

func unmarshalMessageReservedRange(b []byte) (r [2]pref.FieldNumber) {
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.DescriptorProto_ReservedRange_Start:
				r[0] = pref.FieldNumber(v)
			case fieldnum.DescriptorProto_ReservedRange_End:
				r[1] = pref.FieldNumber(v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
	return r
}

func unmarshalMessageExtensionRange(b []byte) (r [2]pref.FieldNumber, opts []byte) {
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.DescriptorProto_ExtensionRange_Start:
				r[0] = pref.FieldNumber(v)
			case fieldnum.DescriptorProto_ExtensionRange_End:
				r[1] = pref.FieldNumber(v)
			}
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.DescriptorProto_ExtensionRange_Options:
				opts = append(opts, v...)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
	return r, opts
}

func (fd *fieldDesc) unmarshalFull(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	fd.parentFile = pf
	fd.parent = pd
	fd.index = i

	var rawDefVal []byte
	var rawTypeName []byte
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldDescriptorProto_Number:
				fd.number = pref.FieldNumber(v)
			case fieldnum.FieldDescriptorProto_Label:
				fd.cardinality = pref.Cardinality(v)
			case fieldnum.FieldDescriptorProto_Type:
				fd.kind = pref.Kind(v)
			case fieldnum.FieldDescriptorProto_OneofIndex:
				// In messageDesc.UnmarshalFull, we allocate slices for both
				// the field and oneof descriptors before unmarshaling either
				// of them. This ensures pointers to slice elements are stable.
				od := &pd.(messageType).lazy.oneofs.list[v]
				od.fields.list = append(od.fields.list, fd)
				if fd.oneofType != nil {
					panic("oneof type already set")
				}
				fd.oneofType = od
			}
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldDescriptorProto_Name:
				fd.fullName = nb.AppendFullName(pd.FullName(), v)
			case fieldnum.FieldDescriptorProto_JsonName:
				fd.hasJSONName = true
				fd.jsonName = nb.MakeString(v)
			case fieldnum.FieldDescriptorProto_DefaultValue:
				fd.defVal.has = true
				rawDefVal = v
			case fieldnum.FieldDescriptorProto_TypeName:
				rawTypeName = v
			case fieldnum.FieldDescriptorProto_Options:
				fd.unmarshalOptions(v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}

	if !fd.hasJSONName {
		fd.jsonName = nb.MakeJSONName(fd.Name())
	}
	if rawDefVal != nil {
		var err error
		fd.defVal.val, err = defval.Unmarshal(string(rawDefVal), fd.kind, defval.Descriptor)
		if err != nil {
			panic(err)
		}
	}
	if fd.isWeak {
		if len(rawTypeName) == 0 || rawTypeName[0] != '.' {
			panic("weak target name must be fully qualified")
		}
		// Consult the global registry for weak messages.
		name := pref.FullName(rawTypeName[1:])
		fd.messageType, _ = preg.GlobalFiles.FindMessageByName(name)
		if fd.messageType == nil {
			fd.messageType = ptype.PlaceholderMessage(pref.FullName(rawTypeName[1:]))
		}
	}
}

func (fd *fieldDesc) unmarshalOptions(b []byte) {
	fd.options = append(fd.options, b...)
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldOptions_Packed:
				fd.hasPacked = true
				fd.isPacked = wire.DecodeBool(v)
			case fieldnum.FieldOptions_Weak:
				fd.isWeak = wire.DecodeBool(v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}

func (od *oneofDesc) unmarshalFull(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	od.parentFile = pf
	od.parent = pd
	od.index = i

	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.OneofDescriptorProto_Name:
				od.fullName = nb.AppendFullName(pd.FullName(), v)
			case fieldnum.OneofDescriptorProto_Options:
				od.options = append(od.options, v...)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}

func (xd *extensionDesc) unmarshalFull(b []byte, nb *nameBuilder) {
	var rawDefVal []byte
	xd.lazy = new(extensionLazy)
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldDescriptorProto_Label:
				xd.lazy.cardinality = pref.Cardinality(v)
			case fieldnum.FieldDescriptorProto_Type:
				xd.lazy.kind = pref.Kind(v)
			}
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldDescriptorProto_JsonName:
				xd.lazy.hasJSONName = true
				xd.lazy.jsonName = nb.MakeString(v)
			case fieldnum.FieldDescriptorProto_DefaultValue:
				xd.lazy.defVal.has = true
				rawDefVal = v
			case fieldnum.FieldDescriptorProto_Options:
				xd.unmarshalOptions(v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}

	if rawDefVal != nil {
		var err error
		xd.lazy.defVal.val, err = defval.Unmarshal(string(rawDefVal), xd.lazy.kind, defval.Descriptor)
		if err != nil {
			panic(err)
		}
	}
}

func (xd *extensionDesc) unmarshalOptions(b []byte) {
	xd.lazy.options = append(xd.lazy.options, b...)
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldOptions_Packed:
				xd.lazy.isPacked = wire.DecodeBool(v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}

func (sd *serviceDesc) unmarshalFull(b []byte, nb *nameBuilder) {
	var rawMethods [][]byte
	sd.lazy = new(serviceLazy)
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.ServiceDescriptorProto_Method:
				rawMethods = append(rawMethods, v)
			case fieldnum.ServiceDescriptorProto_Options:
				sd.lazy.options = append(sd.lazy.options, v...)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}

	if len(rawMethods) > 0 {
		sd.lazy.methods.list = make([]methodDesc, len(rawMethods))
		for i, b := range rawMethods {
			sd.lazy.methods.list[i].unmarshalFull(b, nb, sd.parentFile, sd, i)
		}
	}
}

func (md *methodDesc) unmarshalFull(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	md.parentFile = pf
	md.parent = pd
	md.index = i

	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.MethodDescriptorProto_ClientStreaming:
				md.isStreamingClient = wire.DecodeBool(v)
			case fieldnum.MethodDescriptorProto_ServerStreaming:
				md.isStreamingServer = wire.DecodeBool(v)
			}
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.MethodDescriptorProto_Name:
				md.fullName = nb.AppendFullName(pd.FullName(), v)
			case fieldnum.MethodDescriptorProto_Options:
				md.options = append(md.options, v...)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}
