// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileinit

import (
	wire "google.golang.org/protobuf/internal/encoding/wire"
	fieldnum "google.golang.org/protobuf/internal/fieldnum"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

func newFileDesc(fb FileBuilder) *fileDesc {
	file := &fileDesc{fileInit: fileInit{
		GoTypes:           fb.GoTypes,
		DependencyIndexes: fb.DependencyIndexes,
	}, rawDesc: fb.RawDescriptor}
	file.initDecls(len(fb.EnumOutputTypes), len(fb.MessageOutputTypes), len(fb.ExtensionOutputTypes))
	file.unmarshalSeed(fb.RawDescriptor)

	// Extended message dependencies are eagerly handled since registration
	// needs this information at program init time.
	for i := range file.allExtensions {
		xd := &file.allExtensions[i]
		xd.extendedType = file.popMessageDependency()
	}

	file.checkDecls()
	return file
}

// initDecls pre-allocates slices for the exact number of enums, messages
// (including map entries), and extensions declared in the proto file.
// This is done to avoid regrowing the slice, which would change the address
// for any previously seen declaration.
//
// The alloc methods "allocates" slices by pulling from the capacity.
func (fd *fileDecls) initDecls(numEnums, numMessages, numExtensions int) {
	*fd = fileDecls{
		allEnums:      make([]enumDesc, 0, numEnums),
		allMessages:   make([]messageDesc, 0, numMessages),
		allExtensions: make([]extensionDesc, 0, numExtensions),
	}
}

func (fd *fileDecls) allocEnums(n int) []enumDesc {
	total := len(fd.allEnums)
	es := fd.allEnums[total : total+n]
	fd.allEnums = fd.allEnums[:total+n]
	return es
}
func (fd *fileDecls) allocMessages(n int) []messageDesc {
	total := len(fd.allMessages)
	ms := fd.allMessages[total : total+n]
	fd.allMessages = fd.allMessages[:total+n]
	return ms
}
func (fd *fileDecls) allocExtensions(n int) []extensionDesc {
	total := len(fd.allExtensions)
	xs := fd.allExtensions[total : total+n]
	fd.allExtensions = fd.allExtensions[:total+n]
	return xs
}

// checkDecls performs a sanity check that the expected number of expected
// declarations matches the number that were found in the descriptor proto.
func (fd *fileDecls) checkDecls() {
	if len(fd.allEnums) != cap(fd.allEnums) ||
		len(fd.allMessages) != cap(fd.allMessages) ||
		len(fd.allExtensions) != cap(fd.allExtensions) {
		panic("mismatching cardinality")
	}
}

func (fd *fileDesc) unmarshalSeed(b []byte) {
	nb := getNameBuilder()
	defer putNameBuilder(nb)

	var prevField pref.FieldNumber
	var numEnums, numMessages, numExtensions, numServices int
	var posEnums, posMessages, posExtensions, posServices int
	b0 := b
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.FileDescriptorProto_Name:
				fd.path = nb.MakeString(v)
			case fieldnum.FileDescriptorProto_Package:
				fd.protoPackage = pref.FullName(nb.MakeString(v))
			case fieldnum.FileDescriptorProto_EnumType:
				if prevField != fieldnum.FileDescriptorProto_EnumType {
					if numEnums > 0 {
						panic("non-contiguous repeated field")
					}
					posEnums = len(b0) - len(b) - n - m
				}
				numEnums++
			case fieldnum.FileDescriptorProto_MessageType:
				if prevField != fieldnum.FileDescriptorProto_MessageType {
					if numMessages > 0 {
						panic("non-contiguous repeated field")
					}
					posMessages = len(b0) - len(b) - n - m
				}
				numMessages++
			case fieldnum.FileDescriptorProto_Extension:
				if prevField != fieldnum.FileDescriptorProto_Extension {
					if numExtensions > 0 {
						panic("non-contiguous repeated field")
					}
					posExtensions = len(b0) - len(b) - n - m
				}
				numExtensions++
			case fieldnum.FileDescriptorProto_Service:
				if prevField != fieldnum.FileDescriptorProto_Service {
					if numServices > 0 {
						panic("non-contiguous repeated field")
					}
					posServices = len(b0) - len(b) - n - m
				}
				numServices++
			}
			prevField = num
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
			prevField = -1 // ignore known field numbers of unknown wire type
		}
	}

	// Must allocate all declarations before parsing each descriptor type
	// to ensure we handled all descriptors in "flattened ordering".
	if numEnums > 0 {
		fd.enums.list = fd.allocEnums(numEnums)
	}
	if numMessages > 0 {
		fd.messages.list = fd.allocMessages(numMessages)
	}
	if numExtensions > 0 {
		fd.extensions.list = fd.allocExtensions(numExtensions)
	}
	if numServices > 0 {
		fd.services.list = make([]serviceDesc, numServices)
	}

	if numEnums > 0 {
		b := b0[posEnums:]
		for i := range fd.enums.list {
			_, n := wire.ConsumeVarint(b)
			v, m := wire.ConsumeBytes(b[n:])
			fd.enums.list[i].unmarshalSeed(v, nb, fd, fd, i)
			b = b[n+m:]
		}
	}
	if numMessages > 0 {
		b := b0[posMessages:]
		for i := range fd.messages.list {
			_, n := wire.ConsumeVarint(b)
			v, m := wire.ConsumeBytes(b[n:])
			fd.messages.list[i].unmarshalSeed(v, nb, fd, fd, i)
			b = b[n+m:]
		}
	}
	if numExtensions > 0 {
		b := b0[posExtensions:]
		for i := range fd.extensions.list {
			_, n := wire.ConsumeVarint(b)
			v, m := wire.ConsumeBytes(b[n:])
			fd.extensions.list[i].unmarshalSeed(v, nb, fd, fd, i)
			b = b[n+m:]
		}
	}
	if numServices > 0 {
		b := b0[posServices:]
		for i := range fd.services.list {
			_, n := wire.ConsumeVarint(b)
			v, m := wire.ConsumeBytes(b[n:])
			fd.services.list[i].unmarshalSeed(v, nb, fd, fd, i)
			b = b[n+m:]
		}
	}
}

func (ed *enumDesc) unmarshalSeed(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	ed.parentFile = pf
	ed.parent = pd
	ed.index = i

	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.EnumDescriptorProto_Name:
				ed.fullName = nb.AppendFullName(pd.FullName(), v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}

func (md *messageDesc) unmarshalSeed(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	md.parentFile = pf
	md.parent = pd
	md.index = i

	var prevField pref.FieldNumber
	var numEnums, numMessages, numExtensions int
	var posEnums, posMessages, posExtensions int
	b0 := b
	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.DescriptorProto_Name:
				md.fullName = nb.AppendFullName(pd.FullName(), v)
			case fieldnum.DescriptorProto_EnumType:
				if prevField != fieldnum.DescriptorProto_EnumType {
					if numEnums > 0 {
						panic("non-contiguous repeated field")
					}
					posEnums = len(b0) - len(b) - n - m
				}
				numEnums++
			case fieldnum.DescriptorProto_NestedType:
				if prevField != fieldnum.DescriptorProto_NestedType {
					if numMessages > 0 {
						panic("non-contiguous repeated field")
					}
					posMessages = len(b0) - len(b) - n - m
				}
				numMessages++
			case fieldnum.DescriptorProto_Extension:
				if prevField != fieldnum.DescriptorProto_Extension {
					if numExtensions > 0 {
						panic("non-contiguous repeated field")
					}
					posExtensions = len(b0) - len(b) - n - m
				}
				numExtensions++
			}
			prevField = num
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
			prevField = -1 // ignore known field numbers of unknown wire type
		}
	}

	// Must allocate all declarations before parsing each descriptor type
	// to ensure we handled all descriptors in "flattened ordering".
	if numEnums > 0 {
		md.enums.list = md.parentFile.allocEnums(numEnums)
	}
	if numMessages > 0 {
		md.messages.list = md.parentFile.allocMessages(numMessages)
	}
	if numExtensions > 0 {
		md.extensions.list = md.parentFile.allocExtensions(numExtensions)
	}

	if numEnums > 0 {
		b := b0[posEnums:]
		for i := range md.enums.list {
			_, n := wire.ConsumeVarint(b)
			v, m := wire.ConsumeBytes(b[n:])
			md.enums.list[i].unmarshalSeed(v, nb, pf, md.asDesc(), i)
			b = b[n+m:]
		}
	}
	if numMessages > 0 {
		b := b0[posMessages:]
		for i := range md.messages.list {
			_, n := wire.ConsumeVarint(b)
			v, m := wire.ConsumeBytes(b[n:])
			md.messages.list[i].unmarshalSeed(v, nb, pf, md.asDesc(), i)
			b = b[n+m:]
		}
	}
	if numExtensions > 0 {
		b := b0[posExtensions:]
		for i := range md.extensions.list {
			_, n := wire.ConsumeVarint(b)
			v, m := wire.ConsumeBytes(b[n:])
			md.extensions.list[i].unmarshalSeed(v, nb, pf, md.asDesc(), i)
			b = b[n+m:]
		}
	}
}

func (xd *extensionDesc) unmarshalSeed(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	xd.parentFile = pf
	xd.parent = pd
	xd.index = i

	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.VarintType:
			v, m := wire.ConsumeVarint(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldDescriptorProto_Number:
				xd.number = pref.FieldNumber(v)
			}
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.FieldDescriptorProto_Name:
				xd.fullName = nb.AppendFullName(pd.FullName(), v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}

func (sd *serviceDesc) unmarshalSeed(b []byte, nb *nameBuilder, pf *fileDesc, pd pref.Descriptor, i int) {
	sd.parentFile = pf
	sd.parent = pd
	sd.index = i

	for len(b) > 0 {
		num, typ, n := wire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case wire.BytesType:
			v, m := wire.ConsumeBytes(b)
			b = b[m:]
			switch num {
			case fieldnum.ServiceDescriptorProto_Name:
				sd.fullName = nb.AppendFullName(pd.FullName(), v)
			}
		default:
			m := wire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}
