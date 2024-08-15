// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_reservation

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ReservationResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"block": schema.StringAttribute{
				Required: 					true,
				Description:         "Name of the target Block",
				MarkdownDescription: "Name of the target Block",
			},
			"cidr": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("x.x.x.x/x"), ""),
				},
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"created_on": schema.NumberAttribute{
				Computed: true,
			},
			"desc": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"reservation": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "ID of the target Reservation",
				MarkdownDescription: "ID of the target Reservation",
			},
			"reverse_search": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"settled_by": schema.StringAttribute{
				Computed: true,
			},
			"settled_on": schema.NumberAttribute{
				Computed: true,
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"smallest_cidr": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"space": schema.StringAttribute{
				Required: 					true,
				Description:         "Name of the target Space",
				MarkdownDescription: "Name of the target Space",
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"tag": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{},
				CustomType: TagType{
					ObjectType: types.ObjectType{
						AttrTypes: TagValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
		},
	}
}

type ReservationModel struct {
	Block         types.String `tfsdk:"block"`
	Cidr          types.String `tfsdk:"cidr"`
	CreatedBy     types.String `tfsdk:"created_by"`
	CreatedOn     types.Number `tfsdk:"created_on"`
	Desc          types.String `tfsdk:"desc"`
	Id            types.String `tfsdk:"id"`
	Reservation   types.String `tfsdk:"reservation"`
	ReverseSearch types.Bool   `tfsdk:"reverse_search"`
	SettledBy     types.String `tfsdk:"settled_by"`
	SettledOn     types.Number `tfsdk:"settled_on"`
	Size          types.Int64  `tfsdk:"size"`
	SmallestCidr  types.Bool   `tfsdk:"smallest_cidr"`
	Space         types.String `tfsdk:"space"`
	Status        types.String `tfsdk:"status"`
	Tag           TagValue     `tfsdk:"tag"`
}

var _ basetypes.ObjectTypable = TagType{}

type TagType struct {
	basetypes.ObjectType
}

func (t TagType) Equal(o attr.Type) bool {
	other, ok := o.(TagType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t TagType) String() string {
	return "TagType"
}

func (t TagType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if diags.HasError() {
		return nil, diags
	}

	return TagValue{
		state: attr.ValueStateKnown,
	}, diags
}

func NewTagValueNull() TagValue {
	return TagValue{
		state: attr.ValueStateNull,
	}
}

func NewTagValueUnknown() TagValue {
	return TagValue{
		state: attr.ValueStateUnknown,
	}
}

func NewTagValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (TagValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing TagValue Attribute Value",
				"While creating a TagValue value, a missing attribute value was detected. "+
					"A TagValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TagValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid TagValue Attribute Type",
				"While creating a TagValue value, an invalid attribute value was detected. "+
					"A TagValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TagValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("TagValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra TagValue Attribute Value",
				"While creating a TagValue value, an extra attribute value was detected. "+
					"A TagValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra TagValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewTagValueUnknown(), diags
	}

	if diags.HasError() {
		return NewTagValueUnknown(), diags
	}

	return TagValue{
		state: attr.ValueStateKnown,
	}, diags
}

func NewTagValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) TagValue {
	object, diags := NewTagValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewTagValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t TagType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewTagValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewTagValueUnknown(), nil
	}

	if in.IsNull() {
		return NewTagValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewTagValueMust(TagValue{}.AttributeTypes(ctx), attributes), nil
}

func (t TagType) ValueType(ctx context.Context) attr.Value {
	return TagValue{}
}

var _ basetypes.ObjectValuable = TagValue{}

type TagValue struct {
	state attr.ValueState
}

func (v TagValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 0)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 0)

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v TagValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v TagValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v TagValue) String() string {
	return "TagValue"
}

func (v TagValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{})

	return objVal, diags
}

func (v TagValue) Equal(o attr.Value) bool {
	other, ok := o.(TagValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	return true
}

func (v TagValue) Type(ctx context.Context) attr.Type {
	return TagType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v TagValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{}
}