// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/alertmodelpayload.proto

package alertmodel

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AlertModelPayloadMessage struct {
	Id                    int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	SeverityType          int32    `protobuf:"varint,5,opt,name=severityType,proto3" json:"severityType,omitempty"`
	SeverityName          string   `protobuf:"bytes,6,opt,name=severityName,proto3" json:"severityName,omitempty"`
	AlertEntityId         int32    `protobuf:"varint,7,opt,name=alertEntityId,proto3" json:"alertEntityId,omitempty"`
	AlertEntityName       string   `protobuf:"bytes,8,opt,name=alertEntityName,proto3" json:"alertEntityName,omitempty"`
	AlertDeviceGroup      int32    `protobuf:"varint,9,opt,name=alertDeviceGroup,proto3" json:"alertDeviceGroup,omitempty"`
	AlertDeviceIpAddress  string   `protobuf:"bytes,10,opt,name=alertDeviceIpAddress,proto3" json:"alertDeviceIpAddress,omitempty"`
	AlertDeviceMacAddress string   `protobuf:"bytes,11,opt,name=alertDeviceMacAddress,proto3" json:"alertDeviceMacAddress,omitempty"`
	AlertDeviceIdentifier string   `protobuf:"bytes,12,opt,name=alertDeviceIdentifier,proto3" json:"alertDeviceIdentifier,omitempty"`
	AlertDeviceAssetTag   string   `protobuf:"bytes,13,opt,name=alertDeviceAssetTag,proto3" json:"alertDeviceAssetTag,omitempty"`
	AlertDeviceType       int64    `protobuf:"varint,14,opt,name=alertDeviceType,proto3" json:"alertDeviceType,omitempty"`
	DefinitionId          int64    `protobuf:"varint,15,opt,name=definitionId,proto3" json:"definitionId,omitempty"`
	CatalogName           string   `protobuf:"bytes,16,opt,name=catalogName,proto3" json:"catalogName,omitempty"`
	CategoryId            int32    `protobuf:"varint,17,opt,name=categoryId,proto3" json:"categoryId,omitempty"`
	CategoryName          string   `protobuf:"bytes,18,opt,name=categoryName,proto3" json:"categoryName,omitempty"`
	SubCategoryId         int32    `protobuf:"varint,19,opt,name=subCategoryId,proto3" json:"subCategoryId,omitempty"`
	SubCategoryName       string   `protobuf:"bytes,20,opt,name=subCategoryName,proto3" json:"subCategoryName,omitempty"`
	StatusType            int32    `protobuf:"varint,21,opt,name=statusType,proto3" json:"statusType,omitempty"`
	StatusName            string   `protobuf:"bytes,22,opt,name=statusName,proto3" json:"statusName,omitempty"`
	TimeStamp             string   `protobuf:"bytes,23,opt,name=timeStamp,proto3" json:"timeStamp,omitempty"`
	Message               string   `protobuf:"bytes,24,opt,name=message,proto3" json:"message,omitempty"`
	EemiMessage           string   `protobuf:"bytes,25,opt,name=eemiMessage,proto3" json:"eemiMessage,omitempty"`
	RecommendedAction     string   `protobuf:"bytes,26,opt,name=recommendedAction,proto3" json:"recommendedAction,omitempty"`
	AlertMessageId        string   `protobuf:"bytes,27,opt,name=alertMessageId,proto3" json:"alertMessageId,omitempty"`
	AlertVarBindDetails   string   `protobuf:"bytes,28,opt,name=alertVarBindDetails,proto3" json:"alertVarBindDetails,omitempty"`
	AlertMessageType      string   `protobuf:"bytes,29,opt,name=alertMessageType,proto3" json:"alertMessageType,omitempty"`
	MessageArgs           string   `protobuf:"bytes,30,opt,name=messageArgs,proto3" json:"messageArgs,omitempty"`
	OriginOfCondition     string   `protobuf:"bytes,31,opt,name=originOfCondition,proto3" json:"originOfCondition,omitempty"`
	OriginalDeviceType    int64    `protobuf:"varint,32,opt,name=originalDeviceType,proto3" json:"originalDeviceType,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *AlertModelPayloadMessage) Reset()         { *m = AlertModelPayloadMessage{} }
func (m *AlertModelPayloadMessage) String() string { return proto.CompactTextString(m) }
func (*AlertModelPayloadMessage) ProtoMessage()    {}
func (*AlertModelPayloadMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_68e636160d8cfe65, []int{0}
}

func (m *AlertModelPayloadMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlertModelPayloadMessage.Unmarshal(m, b)
}
func (m *AlertModelPayloadMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlertModelPayloadMessage.Marshal(b, m, deterministic)
}
func (m *AlertModelPayloadMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlertModelPayloadMessage.Merge(m, src)
}
func (m *AlertModelPayloadMessage) XXX_Size() int {
	return xxx_messageInfo_AlertModelPayloadMessage.Size(m)
}
func (m *AlertModelPayloadMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_AlertModelPayloadMessage.DiscardUnknown(m)
}

var xxx_messageInfo_AlertModelPayloadMessage proto.InternalMessageInfo

func (m *AlertModelPayloadMessage) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetSeverityType() int32 {
	if m != nil {
		return m.SeverityType
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetSeverityName() string {
	if m != nil {
		return m.SeverityName
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertEntityId() int32 {
	if m != nil {
		return m.AlertEntityId
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetAlertEntityName() string {
	if m != nil {
		return m.AlertEntityName
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertDeviceGroup() int32 {
	if m != nil {
		return m.AlertDeviceGroup
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetAlertDeviceIpAddress() string {
	if m != nil {
		return m.AlertDeviceIpAddress
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertDeviceMacAddress() string {
	if m != nil {
		return m.AlertDeviceMacAddress
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertDeviceIdentifier() string {
	if m != nil {
		return m.AlertDeviceIdentifier
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertDeviceAssetTag() string {
	if m != nil {
		return m.AlertDeviceAssetTag
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertDeviceType() int64 {
	if m != nil {
		return m.AlertDeviceType
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetDefinitionId() int64 {
	if m != nil {
		return m.DefinitionId
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetCatalogName() string {
	if m != nil {
		return m.CatalogName
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetCategoryId() int32 {
	if m != nil {
		return m.CategoryId
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetCategoryName() string {
	if m != nil {
		return m.CategoryName
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetSubCategoryId() int32 {
	if m != nil {
		return m.SubCategoryId
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetSubCategoryName() string {
	if m != nil {
		return m.SubCategoryName
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetStatusType() int32 {
	if m != nil {
		return m.StatusType
	}
	return 0
}

func (m *AlertModelPayloadMessage) GetStatusName() string {
	if m != nil {
		return m.StatusName
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetTimeStamp() string {
	if m != nil {
		return m.TimeStamp
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetEemiMessage() string {
	if m != nil {
		return m.EemiMessage
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetRecommendedAction() string {
	if m != nil {
		return m.RecommendedAction
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertMessageId() string {
	if m != nil {
		return m.AlertMessageId
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertVarBindDetails() string {
	if m != nil {
		return m.AlertVarBindDetails
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetAlertMessageType() string {
	if m != nil {
		return m.AlertMessageType
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetMessageArgs() string {
	if m != nil {
		return m.MessageArgs
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetOriginOfCondition() string {
	if m != nil {
		return m.OriginOfCondition
	}
	return ""
}

func (m *AlertModelPayloadMessage) GetOriginalDeviceType() int64 {
	if m != nil {
		return m.OriginalDeviceType
	}
	return 0
}

func init() {
	proto.RegisterType((*AlertModelPayloadMessage)(nil), "alertmodel.AlertModelPayloadMessage")
}

func init() {
	proto.RegisterFile("api/alertmodelpayload.proto", fileDescriptor_68e636160d8cfe65)
}

var fileDescriptor_68e636160d8cfe65 = []byte{
	// 560 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x94, 0x51, 0x6f, 0xd3, 0x30,
	0x10, 0xc7, 0xd5, 0xa1, 0x6d, 0xd4, 0xdb, 0xba, 0xd5, 0xdb, 0xc0, 0xb0, 0x31, 0xaa, 0x09, 0xa1,
	0x0a, 0xa1, 0x80, 0x80, 0x17, 0x1e, 0xdb, 0x0d, 0xa1, 0x3e, 0x14, 0xaa, 0x32, 0xf1, 0xee, 0xc6,
	0xd7, 0xc8, 0x52, 0x62, 0x47, 0xb6, 0x5b, 0xa9, 0x5f, 0x83, 0x4f, 0x8c, 0x7c, 0x6e, 0x1b, 0xa7,
	0xe9, 0x4b, 0x24, 0xff, 0xfe, 0x77, 0xff, 0xf8, 0xce, 0x67, 0x93, 0x1b, 0x5e, 0xca, 0x4f, 0x3c,
	0x07, 0xe3, 0x0a, 0x2d, 0x20, 0x2f, 0xf9, 0x2a, 0xd7, 0x5c, 0x24, 0xa5, 0xd1, 0x4e, 0x53, 0x52,
	0x09, 0xf7, 0xff, 0xda, 0x84, 0x0d, 0xfc, 0x72, 0xec, 0x97, 0x93, 0x10, 0x37, 0x06, 0x6b, 0x79,
	0x06, 0xb4, 0x43, 0x0e, 0xa4, 0x60, 0xad, 0x5e, 0xab, 0x7f, 0x38, 0x3d, 0x90, 0x82, 0xde, 0x93,
	0x53, 0x0b, 0x4b, 0x30, 0xd2, 0xad, 0x9e, 0x56, 0x25, 0xb0, 0x43, 0x54, 0x6a, 0x2c, 0x8e, 0xf9,
	0xc5, 0x0b, 0x60, 0x47, 0xbd, 0x56, 0xbf, 0x3d, 0xad, 0x31, 0xfa, 0x8e, 0x9c, 0xe1, 0x16, 0x7e,
	0x28, 0x27, 0xdd, 0x6a, 0x24, 0xd8, 0x31, 0x1a, 0xd5, 0x21, 0xed, 0x93, 0xf3, 0x08, 0xa0, 0xd9,
	0x73, 0x34, 0xdb, 0xc5, 0xf4, 0x03, 0xb9, 0x40, 0xf4, 0x08, 0x4b, 0x99, 0xc2, 0x4f, 0xa3, 0x17,
	0x25, 0x6b, 0xa3, 0x65, 0x83, 0xd3, 0x2f, 0xe4, 0x2a, 0x62, 0xa3, 0x72, 0x20, 0x84, 0x01, 0x6b,
	0x19, 0x41, 0xeb, 0xbd, 0x1a, 0xfd, 0x46, 0xae, 0x23, 0x3e, 0xe6, 0xe9, 0x26, 0xe9, 0x04, 0x93,
	0xf6, 0x8b, 0x3b, 0x59, 0x23, 0x01, 0xca, 0xc9, 0xb9, 0x04, 0xc3, 0x4e, 0x1b, 0x59, 0x95, 0x48,
	0x3f, 0x93, 0xcb, 0x48, 0x18, 0x58, 0x0b, 0xee, 0x89, 0x67, 0xec, 0x0c, 0x73, 0xf6, 0x49, 0xdb,
	0x3e, 0x05, 0x8c, 0x07, 0xd3, 0xe9, 0xb5, 0xfa, 0xcf, 0xa6, 0xbb, 0xd8, 0x9f, 0x8d, 0x80, 0xb9,
	0x54, 0xd2, 0x49, 0xad, 0x46, 0x82, 0x9d, 0x63, 0x58, 0x8d, 0xd1, 0x1e, 0x39, 0x49, 0xb9, 0xe3,
	0xb9, 0xce, 0xb0, 0xe3, 0x17, 0xf8, 0xdf, 0x18, 0xd1, 0x3b, 0x42, 0x52, 0xee, 0x20, 0xd3, 0xc6,
	0x1f, 0x5d, 0x17, 0xfb, 0x1c, 0x11, 0xff, 0x97, 0xcd, 0x0a, 0x2d, 0x68, 0x98, 0x80, 0x98, 0xf9,
	0x09, 0xb0, 0x8b, 0xd9, 0x43, 0x65, 0x73, 0x19, 0x26, 0xa0, 0x06, 0x7d, 0x65, 0x11, 0x40, 0xb3,
	0xab, 0x30, 0x01, 0x3b, 0xd8, 0xef, 0xc9, 0x3a, 0xee, 0x16, 0x16, 0xcb, 0xbf, 0x0e, 0x7b, 0xaa,
	0x48, 0xa5, 0xa3, 0xc9, 0x0b, 0x34, 0x89, 0x08, 0xbd, 0x25, 0x6d, 0x27, 0x0b, 0xf8, 0xe3, 0x78,
	0x51, 0xb2, 0x97, 0x28, 0x57, 0x80, 0x32, 0x72, 0x5c, 0x84, 0x2b, 0xc1, 0x18, 0x6a, 0x9b, 0xa5,
	0xef, 0x16, 0x40, 0x21, 0xd7, 0x17, 0x86, 0xbd, 0x0a, 0xdd, 0x8a, 0x10, 0xfd, 0x48, 0xba, 0x06,
	0x52, 0x5d, 0x14, 0xa0, 0x04, 0x88, 0x41, 0xea, 0xdb, 0xcc, 0x5e, 0x63, 0x5c, 0x53, 0xa0, 0xef,
	0x49, 0x07, 0x0f, 0x6d, 0x9d, 0x3d, 0x12, 0xec, 0x06, 0x43, 0x77, 0xe8, 0x76, 0x4a, 0xfe, 0x72,
	0x33, 0x94, 0x4a, 0x3c, 0x82, 0xe3, 0x32, 0xb7, 0xec, 0x36, 0x9a, 0x92, 0xba, 0xb4, 0xbd, 0x23,
	0x6b, 0x0f, 0xec, 0xd3, 0x1b, 0x0c, 0x6f, 0x70, 0x5f, 0xd5, 0xba, 0xc0, 0x81, 0xc9, 0x2c, 0xbb,
	0x0b, 0x55, 0x45, 0xc8, 0x57, 0xa5, 0x8d, 0xcc, 0xa4, 0xfa, 0x3d, 0x7f, 0xd0, 0x4a, 0xe0, 0xf0,
	0xb0, 0xb7, 0xa1, 0xaa, 0x86, 0x40, 0x13, 0x42, 0x03, 0xe4, 0x79, 0x34, 0xa4, 0x3d, 0x9c, 0xbe,
	0x3d, 0xca, 0x70, 0x42, 0xbe, 0xa7, 0xba, 0x48, 0x04, 0xe4, 0x79, 0x02, 0xca, 0x81, 0x29, 0x8d,
	0xb4, 0x90, 0xf8, 0x66, 0x69, 0x95, 0x48, 0xe5, 0x20, 0x33, 0xdc, 0x3b, 0x27, 0xb9, 0x9c, 0x25,
	0xb0, 0x04, 0xe5, 0xc2, 0x37, 0x97, 0xb3, 0x61, 0xb7, 0xf1, 0x9c, 0x4d, 0x5a, 0xb3, 0x23, 0x7c,
	0xf9, 0xbe, 0xfe, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x66, 0xd7, 0x57, 0xc0, 0x18, 0x05, 0x00, 0x00,
}
