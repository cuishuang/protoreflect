// Code generated by protoc-gen-gosrcinfo. DO NOT EDIT.
// source: desc_test_proto3.proto

package testprotos

import (
	sourceinfo "github.com/jhump/protoreflect/v2/sourceinfo"
)

func init() {
	srcInfo := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x3c, 0xd3, 0xdf, 0x8e, 0xd3, 0x3c,
		0x10, 0x05, 0x70, 0xdb, 0xc7, 0xf6, 0xe7, 0x75, 0xb3, 0x6d, 0x73, 0xd2, 0x3f, 0x5f, 0x93, 0x6e,
		0x9b, 0xa4, 0x15, 0x2b, 0x01, 0xd2, 0x4a, 0x5c, 0x80, 0x40, 0x20, 0x40, 0xac, 0x10, 0xda, 0xe5,
		0x02, 0x71, 0x8f, 0x78, 0xff, 0xb7, 0x40, 0x93, 0x7a, 0x7a, 0x97, 0xdf, 0x9e, 0x99, 0xc9, 0xd8,
		0x9b, 0xe6, 0x48, 0x6f, 0xcc, 0xc9, 0xe6, 0x94, 0x6d, 0x45, 0x18, 0x43, 0x79, 0x4a, 0x84, 0x33,
		0xcf, 0xf9, 0x26, 0xbb, 0x34, 0xbb, 0x3c, 0xa6, 0x6c, 0x1d, 0xe1, 0x4d, 0x23, 0x7f, 0x84, 0x21,
		0xa2, 0x69, 0xa7, 0x47, 0x4b, 0xfc, 0x67, 0x86, 0x9c, 0xb3, 0x0b, 0x86, 0xfe, 0xc6, 0xcc, 0x6d,
		0xce, 0x19, 0xc1, 0x58, 0xe2, 0x26, 0x2c, 0xf2, 0x2c, 0xfb, 0x60, 0x9c, 0x21, 0x72, 0x5a, 0xe5,
		0x2a, 0x07, 0x81, 0x15, 0x2d, 0x54, 0x8e, 0xc8, 0x6c, 0x4a, 0xa1, 0x25, 0x66, 0xe9, 0xff, 0x12,
		0xd9, 0x49, 0x54, 0x39, 0x62, 0xb6, 0xde, 0x96, 0x42, 0x47, 0x54, 0xa9, 0x29, 0x91, 0xb4, 0x55,
		0x69, 0xae, 0x92, 0xac, 0x66, 0x29, 0x04, 0x71, 0x7b, 0x2d, 0x94, 0x6d, 0x6f, 0xaf, 0x85, 0x70,
		0xc4, 0x6d, 0x4d, 0xd9, 0xdd, 0x1b, 0xfa, 0xa5, 0xd9, 0x4e, 0xbb, 0x7b, 0x59, 0x70, 0x99, 0xa6,
		0x95, 0xfc, 0xb4, 0x7b, 0x9d, 0xce, 0xd2, 0x22, 0xf0, 0xa2, 0xa5, 0x2a, 0x12, 0x75, 0xdd, 0xa9,
		0x2c, 0x51, 0xef, 0x8f, 0x2a, 0x10, 0xf5, 0x78, 0x2a, 0x43, 0x2c, 0xc1, 0xb4, 0x2d, 0x91, 0x0d,
		0xa2, 0xb9, 0x4a, 0xb2, 0x05, 0x55, 0x20, 0xb8, 0xde, 0x94, 0x36, 0x47, 0x34, 0x69, 0x5f, 0x22,
		0x17, 0x45, 0x8d, 0xca, 0x12, 0xcd, 0x4a, 0x47, 0xca, 0x31, 0x9b, 0xb6, 0x2b, 0x6d, 0x20, 0x56,
		0xe9, 0x6b, 0x89, 0x10, 0x45, 0x6f, 0x55, 0x96, 0x58, 0xbd, 0xfb, 0xa8, 0x92, 0xca, 0xcf, 0x5f,
		0x4a, 0x9b, 0x27, 0xd6, 0xd7, 0x93, 0xfa, 0x28, 0xda, 0xa9, 0x2c, 0xb1, 0x6e, 0xf5, 0x6c, 0x1e,
		0xc4, 0xfa, 0x7a, 0xb6, 0x40, 0x6c, 0xd2, 0xeb, 0x12, 0x85, 0x28, 0xea, 0x55, 0x96, 0xd8, 0x0c,
		0xf7, 0x2a, 0x10, 0x9b, 0x97, 0xaf, 0xa6, 0xfb, 0xb6, 0xf4, 0x3b, 0xb3, 0xbf, 0xdc, 0xb7, 0xdc,
		0xc0, 0x2e, 0xad, 0xa6, 0x71, 0x56, 0xee, 0xbb, 0x4d, 0xa7, 0xa9, 0xc5, 0x4e, 0x37, 0xdc, 0xa6,
		0x56, 0x65, 0x89, 0xb6, 0x3b, 0xa8, 0x40, 0xb4, 0xc3, 0x58, 0xda, 0x2c, 0xd1, 0x25, 0x8d, 0xac,
		0x17, 0x2d, 0x55, 0x81, 0xe8, 0xea, 0x8d, 0x4a, 0x2a, 0xb7, 0x3b, 0x15, 0x88, 0x6e, 0x7f, 0x27,
		0x3b, 0x45, 0x43, 0x7f, 0x90, 0x9f, 0x42, 0xce, 0x88, 0xf2, 0xaa, 0xc3, 0xe5, 0x1b, 0x88, 0xd3,
		0x37, 0x70, 0x4c, 0x3f, 0xa5, 0x25, 0x5e, 0xfe, 0xcf, 0xc7, 0x6a, 0xab, 0x72, 0xc4, 0x71, 0x77,
		0x56, 0x81, 0x38, 0x3e, 0x3c, 0x95, 0x36, 0x4b, 0xf4, 0xe9, 0x47, 0x89, 0xe4, 0xbd, 0x7d, 0xd5,
		0xa9, 0x02, 0xd1, 0xdf, 0x9d, 0x54, 0x8e, 0xe8, 0xcf, 0x0f, 0x2a, 0x10, 0xfd, 0x87, 0xef, 0x65,
		0x88, 0x23, 0x86, 0xf4, 0xb7, 0x44, 0x32, 0x72, 0xa8, 0xf6, 0x2a, 0xc9, 0x0e, 0xef, 0x55, 0x91,
		0x18, 0x1e, 0x9f, 0x54, 0x20, 0x86, 0xe7, 0x3f, 0x65, 0x08, 0x88, 0x31, 0xfd, 0x2e, 0x91, 0x7c,
		0x03, 0x63, 0x75, 0x50, 0x05, 0x62, 0xec, 0x5f, 0xa8, 0x1c, 0x31, 0xde, 0xbf, 0x51, 0x45, 0x62,
		0xfc, 0xf4, 0x4d, 0x25, 0x53, 0x1e, 0x7f, 0xfd, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x4c, 0xf0, 0x36,
		0x6d, 0x33, 0x04, 0x00, 0x00,
	}
	sourceinfo.Register("desc_test_proto3.proto", srcInfo)
}
