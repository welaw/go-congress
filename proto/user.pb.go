// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/user.proto

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type User struct {
	Uid             string   `protobuf:"bytes,2,opt,name=uid" json:"uid,omitempty"`
	Username        string   `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	Email           string   `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	FullName        string   `protobuf:"bytes,5,opt,name=full_name,json=fullName" json:"full_name,omitempty"`
	Biography       string   `protobuf:"bytes,6,opt,name=biography" json:"biography,omitempty"`
	Upstream        string   `protobuf:"bytes,7,opt,name=upstream" json:"upstream,omitempty"`
	UpstreamBody    string   `protobuf:"bytes,8,opt,name=upstream_body,json=upstreamBody" json:"upstream_body,omitempty"`
	PictureUrl      string   `protobuf:"bytes,9,opt,name=picture_url,json=pictureUrl" json:"picture_url,omitempty"`
	ProviderId      string   `protobuf:"bytes,10,opt,name=provider_id,json=providerId" json:"provider_id,omitempty"`
	Name            string   `protobuf:"bytes,11,opt,name=name" json:"name,omitempty"`
	EmailPrivate    bool     `protobuf:"varint,13,opt,name=email_private,json=emailPrivate" json:"email_private,omitempty"`
	Roles           []string `protobuf:"bytes,14,rep,name=roles" json:"roles,omitempty"`
	Provider        string   `protobuf:"bytes,17,opt,name=provider" json:"provider,omitempty"`
	FullNamePrivate bool     `protobuf:"varint,21,opt,name=full_name_private,json=fullNamePrivate" json:"full_name_private,omitempty"`
	BioguideId      string   `protobuf:"bytes,22,opt,name=bioguide_id,json=bioguideId" json:"bioguide_id,omitempty"`
	LisId           string   `protobuf:"bytes,23,opt,name=lis_id,json=lisId" json:"lis_id,omitempty"`
	FoundId         string   `protobuf:"bytes,24,opt,name=found_id,json=foundId" json:"found_id,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto1.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *User) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *User) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *User) GetBiography() string {
	if m != nil {
		return m.Biography
	}
	return ""
}

func (m *User) GetUpstream() string {
	if m != nil {
		return m.Upstream
	}
	return ""
}

func (m *User) GetUpstreamBody() string {
	if m != nil {
		return m.UpstreamBody
	}
	return ""
}

func (m *User) GetPictureUrl() string {
	if m != nil {
		return m.PictureUrl
	}
	return ""
}

func (m *User) GetProviderId() string {
	if m != nil {
		return m.ProviderId
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetEmailPrivate() bool {
	if m != nil {
		return m.EmailPrivate
	}
	return false
}

func (m *User) GetRoles() []string {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *User) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *User) GetFullNamePrivate() bool {
	if m != nil {
		return m.FullNamePrivate
	}
	return false
}

func (m *User) GetBioguideId() string {
	if m != nil {
		return m.BioguideId
	}
	return ""
}

func (m *User) GetLisId() string {
	if m != nil {
		return m.LisId
	}
	return ""
}

func (m *User) GetFoundId() string {
	if m != nil {
		return m.FoundId
	}
	return ""
}

func init() {
	proto1.RegisterType((*User)(nil), "proto.User")
}

func init() { proto1.RegisterFile("proto/user.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x51, 0xcb, 0x4e, 0xeb, 0x30,
	0x10, 0x55, 0x6f, 0xfa, 0x48, 0xa6, 0xed, 0xa5, 0xb5, 0x28, 0x98, 0x87, 0xd4, 0x0a, 0x36, 0x15,
	0x0b, 0x58, 0xf0, 0x07, 0xec, 0xb2, 0x41, 0xa8, 0x52, 0xd7, 0x51, 0x8a, 0x5d, 0xb0, 0xe4, 0xd6,
	0xd1, 0x24, 0xae, 0xd4, 0xdf, 0xe6, 0x0b, 0xd0, 0x8c, 0xe3, 0xb0, 0xf2, 0x9c, 0x73, 0xe6, 0x79,
	0x0c, 0xb3, 0x0a, 0x5d, 0xe3, 0x5e, 0x7c, 0xad, 0xf1, 0x99, 0x43, 0x31, 0xe0, 0xe7, 0xe1, 0x27,
	0x81, 0xfe, 0xb6, 0xd6, 0x28, 0x66, 0x90, 0x78, 0xa3, 0xe4, 0xbf, 0x55, 0x6f, 0x9d, 0x6d, 0x28,
	0x14, 0xb7, 0x90, 0x52, 0xfe, 0xb1, 0x3c, 0x68, 0x99, 0x30, 0xdd, 0x61, 0x71, 0x09, 0x03, 0x7d,
	0x28, 0x8d, 0x95, 0x7d, 0x16, 0x02, 0x10, 0x77, 0x90, 0xed, 0xbd, 0xb5, 0x05, 0x97, 0x0c, 0x42,
	0x09, 0x11, 0xef, 0x54, 0x72, 0x0f, 0xd9, 0xce, 0xb8, 0x2f, 0x2c, 0xab, 0xef, 0xb3, 0x1c, 0xb2,
	0xf8, 0x47, 0xf0, 0xb0, 0xaa, 0x6e, 0x50, 0x97, 0x07, 0x39, 0x6a, 0x87, 0xb5, 0x58, 0x3c, 0xc2,
	0x34, 0xc6, 0xc5, 0xce, 0xa9, 0xb3, 0x4c, 0x39, 0x61, 0x12, 0xc9, 0x37, 0xa7, 0xce, 0x62, 0x09,
	0xe3, 0xca, 0x7c, 0x36, 0x1e, 0x75, 0xe1, 0xd1, 0xca, 0x8c, 0x53, 0xa0, 0xa5, 0xb6, 0x68, 0x39,
	0x01, 0xdd, 0xc9, 0x28, 0x8d, 0x85, 0x51, 0x12, 0xda, 0x84, 0x96, 0xca, 0x95, 0x10, 0xd0, 0xe7,
	0xc5, 0xc7, 0xac, 0x70, 0x4c, 0xa3, 0xf9, 0xb4, 0xa2, 0x42, 0x73, 0x2a, 0x1b, 0x2d, 0xa7, 0xab,
	0xde, 0x3a, 0xdd, 0x4c, 0x98, 0xfc, 0x08, 0x1c, 0x99, 0x81, 0xce, 0xea, 0x5a, 0xfe, 0x5f, 0x25,
	0x64, 0x06, 0x03, 0xba, 0x28, 0x36, 0x97, 0xf3, 0x70, 0x51, 0xc4, 0xe2, 0x09, 0xe6, 0x9d, 0x51,
	0x5d, 0xeb, 0x05, 0xb7, 0xbe, 0x88, 0x86, 0xc5, 0xee, 0x4b, 0x18, 0x93, 0x4d, 0xde, 0x28, 0x4d,
	0x7b, 0x5f, 0x85, 0xbd, 0x23, 0x95, 0x2b, 0xb1, 0x80, 0xa1, 0x35, 0x35, 0x69, 0xd7, 0xe1, 0x33,
	0xac, 0xa9, 0x73, 0x25, 0x6e, 0x20, 0xdd, 0x3b, 0x7f, 0x54, 0x24, 0x48, 0x16, 0x46, 0x8c, 0x73,
	0xb5, 0x1b, 0xf2, 0xdf, 0xbf, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0x35, 0x44, 0x13, 0x34, 0x16,
	0x02, 0x00, 0x00,
}