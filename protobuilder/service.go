package protobuilder

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/jhump/protoreflect/v2/internal"
	"github.com/jhump/protoreflect/v2/protomessage"
)

// ServiceBuilder is a builder used to construct a protoreflect.ServiceDescriptor.
//
// To create a new ServiceBuilder, use NewService.
type ServiceBuilder struct {
	baseBuilder

	Options *descriptorpb.ServiceOptions

	methods []*MethodBuilder
	symbols map[protoreflect.Name]*MethodBuilder
}

var _ Builder = (*ServiceBuilder)(nil)

// NewService creates a new ServiceBuilder for a service with the given name.
func NewService(name protoreflect.Name) *ServiceBuilder {
	return &ServiceBuilder{
		baseBuilder: baseBuilderWithName(name),
		symbols:     map[protoreflect.Name]*MethodBuilder{},
	}
}

// FromService returns a ServiceBuilder that is effectively a copy of the given
// descriptor.
//
// Note that it is not just the given service that is copied but its entire
// file. So the caller can get the parent element of the returned builder and
// the result would be a builder that is effectively a copy of the service
// descriptor's parent file.
//
// This means that service builders created from descriptors do not need to be
// explicitly assigned to a file in order to preserve the original service's
// package name.
func FromService(sd protoreflect.ServiceDescriptor) (*ServiceBuilder, error) {
	if fb, err := FromFile(sd.ParentFile()); err != nil {
		return nil, err
	} else if sb, ok := fb.findFullyQualifiedElement(sd.FullName()).(*ServiceBuilder); ok {
		return sb, nil
	} else {
		return nil, fmt.Errorf("could not find service %s after converting file %q to builder", sd.FullName(), sd.ParentFile().Path())
	}
}

func fromService(sd protoreflect.ServiceDescriptor) (*ServiceBuilder, error) {
	sb := NewService(sd.Name())
	var err error
	sb.Options, err = protomessage.As[*descriptorpb.ServiceOptions](sd.Options())
	if err != nil {
		return nil, err
	}
	setComments(&sb.comments, sd.ParentFile().SourceLocations().ByDescriptor(sd))

	mtds := sd.Methods()
	for i, length := 0, mtds.Len(); i < length; i++ {
		mtd := mtds.Get(i)
		if mtb, err := fromMethod(mtd); err != nil {
			return nil, err
		} else if err := sb.TryAddMethod(mtb); err != nil {
			return nil, err
		}
	}

	return sb, nil
}

// SetName changes this service's name, returning the service builder for method
// chaining. If the given new name is not valid (e.g. TrySetName would have
// returned an error) then this method will panic.
func (sb *ServiceBuilder) SetName(newName protoreflect.Name) *ServiceBuilder {
	if err := sb.TrySetName(newName); err != nil {
		panic(err)
	}
	return sb
}

// TrySetName changes this service's name. It will return an error if the given
// new name is not a valid protobuf identifier or if the parent file builder
// already has an element with the given name.
func (sb *ServiceBuilder) TrySetName(newName protoreflect.Name) error {
	return sb.baseBuilder.setName(sb, newName)
}

// SetComments sets the comments associated with the service. This method
// returns the service builder, for method chaining.
func (sb *ServiceBuilder) SetComments(c Comments) *ServiceBuilder {
	sb.comments = c
	return sb
}

// Children returns any builders assigned to this service builder. These will
// be the service's methods.
func (sb *ServiceBuilder) Children() []Builder {
	var ch []Builder
	for _, mtb := range sb.methods {
		ch = append(ch, mtb)
	}
	return ch
}

func (sb *ServiceBuilder) findChild(name protoreflect.Name) Builder {
	child := sb.symbols[name]
	if child == nil {
		// NB: can't just return child because then it would
		//     be a typed nil, not a nil interface
		return nil
	}
	return child
}

func (sb *ServiceBuilder) removeChild(b Builder) {
	if p, ok := b.Parent().(*ServiceBuilder); !ok || p != sb {
		return
	}
	sb.methods = deleteBuilder(b.Name(), sb.methods).([]*MethodBuilder)
	delete(sb.symbols, b.Name())
	b.setParent(nil)
}

func (sb *ServiceBuilder) renamedChild(b Builder, oldName protoreflect.Name) error {
	if p, ok := b.Parent().(*ServiceBuilder); !ok || p != sb {
		return nil
	}

	if err := sb.addSymbol(b.(*MethodBuilder)); err != nil {
		return err
	}
	delete(sb.symbols, oldName)
	return nil
}

func (sb *ServiceBuilder) addSymbol(b *MethodBuilder) error {
	if _, ok := sb.symbols[b.Name()]; ok {
		return fmt.Errorf("service %s already contains method named %q", FullName(sb), b.Name())
	}
	sb.symbols[b.Name()] = b
	return nil
}

// GetMethod returns the method with the given name. If no such method exists in
// the service, nil is returned.
func (sb *ServiceBuilder) GetMethod(name protoreflect.Name) *MethodBuilder {
	return sb.symbols[name]
}

// RemoveMethod removes the method with the given name. If no such method exists
// in the service, this is a no-op. This returns the service builder, for method
// chaining.
func (sb *ServiceBuilder) RemoveMethod(name protoreflect.Name) *ServiceBuilder {
	sb.TryRemoveMethod(name)
	return sb
}

// TryRemoveMethod removes the method with the given name and returns false if
// the service has no such method.
func (sb *ServiceBuilder) TryRemoveMethod(name protoreflect.Name) bool {
	if mtb, ok := sb.symbols[name]; ok {
		sb.removeChild(mtb)
		return true
	}
	return false
}

// AddMethod adds the given method to this servuce. If an error prevents the
// method  from being added, this method panics. This returns the service
// builder, for method chaining.
func (sb *ServiceBuilder) AddMethod(mtb *MethodBuilder) *ServiceBuilder {
	if err := sb.TryAddMethod(mtb); err != nil {
		panic(err)
	}
	return sb
}

// TryAddMethod adds the given field to this service, returning any error that
// prevents the field from being added (such as a name collision with another
// method already added to the service).
func (sb *ServiceBuilder) TryAddMethod(mtb *MethodBuilder) error {
	if err := sb.addSymbol(mtb); err != nil {
		return err
	}
	Unlink(mtb)
	mtb.setParent(sb)
	sb.methods = append(sb.methods, mtb)
	return nil
}

// SetOptions sets the service options for this service and returns the service,
// for method chaining.
func (sb *ServiceBuilder) SetOptions(options *descriptorpb.ServiceOptions) *ServiceBuilder {
	sb.Options = options
	return sb
}

func (sb *ServiceBuilder) buildProto(path []int32, sourceInfo *descriptorpb.SourceCodeInfo) (*descriptorpb.ServiceDescriptorProto, error) {
	addCommentsTo(sourceInfo, path, &sb.comments)

	methods := make([]*descriptorpb.MethodDescriptorProto, 0, len(sb.methods))
	for _, mtb := range sb.methods {
		path := append(path, internal.ServiceMethodsTag, int32(len(methods)))
		mtd, err := mtb.buildProto(path, sourceInfo)
		if err != nil {
			return nil, err
		}
		methods = append(methods, mtd)
	}

	return &descriptorpb.ServiceDescriptorProto{
		Name:    proto.String(string(sb.name)),
		Options: sb.Options,
		Method:  methods,
	}, nil
}

// Build constructs a service descriptor based on the contents of this service
// builder. If there are any problems constructing the descriptor, including
// resolving symbols referenced by the builder or failing to meet certain
// validation rules, an error is returned.
func (sb *ServiceBuilder) Build() (protoreflect.ServiceDescriptor, error) {
	sd, err := sb.BuildDescriptor()
	if err != nil {
		return nil, err
	}
	return sd.(protoreflect.ServiceDescriptor), nil
}

// BuildDescriptor constructs a service descriptor based on the contents of this
// service builder. Most usages will prefer Build() instead, whose return type
// is a concrete descriptor type. This method is present to satisfy the Builder
// interface.
func (sb *ServiceBuilder) BuildDescriptor() (protoreflect.Descriptor, error) {
	return doBuild(sb, BuilderOptions{})
}

// MethodBuilder is a builder used to construct a protoreflect.MethodDescriptor. A
// method builder *must* be added to a service before calling its Build()
// method.
//
// To create a new MethodBuilder, use NewMethod.
type MethodBuilder struct {
	baseBuilder

	Options  *descriptorpb.MethodOptions
	ReqType  *RpcType
	RespType *RpcType
}

var _ Builder = (*MethodBuilder)(nil)

// NewMethod creates a new MethodBuilder for a method with the given name and
// request and response types.
func NewMethod(name protoreflect.Name, req, resp *RpcType) *MethodBuilder {
	return &MethodBuilder{
		baseBuilder: baseBuilderWithName(name),
		ReqType:     req,
		RespType:    resp,
	}
}

// FromMethod returns a MethodBuilder that is effectively a copy of the given
// descriptor.
//
// Note that it is not just the given method that is copied but its entire file.
// So the caller can get the parent element of the returned builder and the
// result would be a builder that is effectively a copy of the method
// descriptor's parent service.
//
// This means that method builders created from descriptors do not need to be
// explicitly assigned to a file in order to preserve the original method's
// package name.
func FromMethod(mtd protoreflect.MethodDescriptor) (*MethodBuilder, error) {
	if fb, err := FromFile(mtd.ParentFile()); err != nil {
		return nil, err
	} else if mtb, ok := fb.findFullyQualifiedElement(mtd.FullName()).(*MethodBuilder); ok {
		return mtb, nil
	} else {
		return nil, fmt.Errorf("could not find method %s after converting file %q to builder", mtd.FullName(), mtd.ParentFile().Path())
	}
}

func fromMethod(mtd protoreflect.MethodDescriptor) (*MethodBuilder, error) {
	req := RpcTypeImportedMessage(mtd.Input(), mtd.IsStreamingClient())
	resp := RpcTypeImportedMessage(mtd.Output(), mtd.IsStreamingServer())
	mtb := NewMethod(mtd.Name(), req, resp)
	var err error
	mtb.Options, err = protomessage.As[*descriptorpb.MethodOptions](mtd.Options())
	if err != nil {
		return nil, err
	}
	setComments(&mtb.comments, mtd.ParentFile().SourceLocations().ByDescriptor(mtd))

	return mtb, nil
}

// SetName changes this method's name, returning the method builder for method
// chaining. If the given new name is not valid (e.g. TrySetName would have
// returned an error) then this method will panic.
func (mtb *MethodBuilder) SetName(newName protoreflect.Name) *MethodBuilder {
	if err := mtb.TrySetName(newName); err != nil {
		panic(err)
	}
	return mtb
}

// TrySetName changes this method's name. It will return an error if the given
// new name is not a valid protobuf identifier or if the parent service builder
// already has a method with the given name.
func (mtb *MethodBuilder) TrySetName(newName protoreflect.Name) error {
	return mtb.baseBuilder.setName(mtb, newName)
}

// SetComments sets the comments associated with the method. This method
// returns the method builder, for method chaining.
func (mtb *MethodBuilder) SetComments(c Comments) *MethodBuilder {
	mtb.comments = c
	return mtb
}

// Children returns nil, since methods cannot have child elements. It is
// present to satisfy the Builder interface.
func (mtb *MethodBuilder) Children() []Builder {
	// methods do not have children
	return nil
}

func (mtb *MethodBuilder) findChild(_ protoreflect.Name) Builder {
	// methods do not have children
	return nil
}

func (mtb *MethodBuilder) removeChild(_ Builder) {
	// methods do not have children
}

func (mtb *MethodBuilder) renamedChild(_ Builder, _ protoreflect.Name) error {
	// methods do not have children
	return nil
}

// SetOptions sets the method options for this method and returns the method,
// for method chaining.
func (mtb *MethodBuilder) SetOptions(options *descriptorpb.MethodOptions) *MethodBuilder {
	mtb.Options = options
	return mtb
}

// SetRequestType changes the request type for the method and then returns the
// method builder, for method chaining.
func (mtb *MethodBuilder) SetRequestType(t *RpcType) *MethodBuilder {
	mtb.ReqType = t
	return mtb
}

// SetResponseType changes the response type for the method and then returns the
// method builder, for method chaining.
func (mtb *MethodBuilder) SetResponseType(t *RpcType) *MethodBuilder {
	mtb.RespType = t
	return mtb
}

func (mtb *MethodBuilder) buildProto(path []int32, sourceInfo *descriptorpb.SourceCodeInfo) (*descriptorpb.MethodDescriptorProto, error) {
	addCommentsTo(sourceInfo, path, &mtb.comments)

	mtd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String(string(mtb.name)),
		Options:    mtb.Options,
		InputType:  proto.String("." + string(mtb.ReqType.TypeName())),
		OutputType: proto.String("." + string(mtb.RespType.TypeName())),
	}
	if mtb.ReqType.IsStream {
		mtd.ClientStreaming = proto.Bool(true)
	}
	if mtb.RespType.IsStream {
		mtd.ServerStreaming = proto.Bool(true)
	}

	return mtd, nil
}

// Build constructs a method descriptor based on the contents of this method
// builder. If there are any problems constructing the descriptor, including
// resolving symbols referenced by the builder or failing to meet certain
// validation rules, an error is returned.
func (mtb *MethodBuilder) Build() (protoreflect.MethodDescriptor, error) {
	mtd, err := mtb.BuildDescriptor()
	if err != nil {
		return nil, err
	}
	return mtd.(protoreflect.MethodDescriptor), nil
}

// BuildDescriptor constructs a method descriptor based on the contents of this
// method builder. Most usages will prefer Build() instead, whose return type is
// a concrete descriptor type. This method is present to satisfy the Builder
// interface.
func (mtb *MethodBuilder) BuildDescriptor() (protoreflect.Descriptor, error) {
	return doBuild(mtb, BuilderOptions{})
}
