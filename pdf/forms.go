/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

//
// Higher level manipulation of forms (AcroForm).
//

package pdf

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/unidoc/unidoc/common"
)

type PdfAcroForm struct {
	Fields          *[]*PdfField
	NeedAppearances PdfObject
	SigFlags        PdfObject
	CO              PdfObject
	DR              PdfObject
	DA              PdfObject
	Q               PdfObject
	XFA             PdfObject
}

func (r *PdfReader) newPdfAcroFormFromDict(d PdfObjectDictionary) (*PdfAcroForm, error) {
	acroForm := PdfAcroForm{}

	if obj, has := d["Fields"]; has {
		obj, err := r.traceToObject(obj)
		if err != nil {
			return nil, err
		}
		fieldArray, ok := TraceToDirectObject(obj).(*PdfObjectArray)
		if !ok {
			return nil, fmt.Errorf("Fields not an array (%T)", obj)
		}

		fields := []*PdfField{}
		for _, obj := range fieldArray {
			obj, err := r.traceToObject(obj)
			if err != nil {
				return nil, err
			}
			fDict, ok := TraceToDirectObject(obj).(*PdfObjectDictionary)
			if !ok {
				return nil, fmt.Errorf("Invalid Fields entry: %T", obj)
			}
			field := newPdfFieldFromDict(fDict)
			fields = append(fields, field)
		}
		acroForm.Fields = &fields
	}

	if obj, has := d["NeedAppearances"]; has {
		acroForm.NeedAppearances = obj
	}
	if obj, has := d["SigFlags"]; has {
		acroForm.SigFlags = obj
	}
	if obj, has := d["CO"]; has {
		acroForm.CO = obj
	}
	if obj, has := d["DR"]; has {
		acroForm.DR = obj
	}
	if obj, has := d["DA"]; has {
		acroForm.DA = obj
	}
	if obj, has := d["Q"]; has {
		acroForm.Q = obj
	}
	if obj, has := d["XFA"]; has {
		acroForm.XFA = obj
	}

	return &acroForm, nil
}

func (this *PdfAcroForm) ToPdfObject() PdfObject {

}

func (r *PdfReader) newPdfFieldDict(d PdfObjectDictionary) (*PdfField, error) {
	field := PdfField{}

	// Field type (required in terminal fields).
	// Can be /Btn /Tx /Ch /Sig
	// Required for a terminal field (inheritable).
	if obj, has := d["FT"]; has {
		obj, err = r.traceToObject(obj)
		if err != nil {
			return nil, err
		}
		name, ok := obj.(*PdfObjectName)
		if !ok {
			return nil, fmt.Errorf("Invalid type of FT field (%T)", obj)
		}

		acroForm.FT = name
	}

	// In a non-terminal field, the Kids array shall refer to field dictionaries that are immediate descendants of this field.
	// In a terminal field, the Kids array ordinarily shall refer to one or more separate widget annotations that are associated
	// with this field. However, if there is only one associated widget annotation, and its contents have been merged into the field
	// dictionary, Kids shall be omitted.

	// Terminal field if:
	// 1. Kids pointing to widget annotations (Kids[0].(*PdfField) casting fails)
	// 2. Kids empty or missing and a the dict has a Subtype equivalent to "Widget"

	if obj, has := d["Parent"]; has {
		field.Parent = obj
	}
	if obj, has := d["Kids"]; has {
		field.Kids = obj
	}

	// Partial field name (Optional)
	if obj, has := d["T"]; has {
		field.T = obj
	}
	// Alternate description (Optional)
	if obj, has := d["TU"]; has {
		field.TU = obj
	}
	// Mapping name (Optional)
	if obj, has := d["TM"]; has {
		field.TM = obj
	}
	// Field flag. (Optional; inheritable)
	if obj, has := d["Ff"]; has {
		field.Ff = obj
	}
	// Value (Optional; inheritable) - Various types depending on the field type.
	if obj, has := d["V"]; has {
		field.V = obj
	}
	// Default value for reset (Optional; inheritable)
	if obj, has := d["DV"]; has {
		field.DV = obj
	}
	// Additional actions dictionary (Optional)
	if obj, has := d["AA"]; has {
		field.AA = obj
	}

	return &field, nil
}

type PdfAnnotationWidget struct {
	Subtype PdfObject // Widget (required)
	H       PdfObject
	MK      PdfObject
	A       PdfObject
	AA      PdfObject
	BS      PdfObject
	Parent  *PdfIndirectObject // Max 1 parent; Gets tricky for both form and annotation refs?  Seems to usually refer to the page one.
}

type PdfField struct {
	FT     *PdfObjectName // field type
	Parent *PdfField
	// In a non-terminal field, the Kids array shall refer to field dictionaries that are immediate descendants of this field.
	// In a terminal field, the Kids array ordinarily shall refer to one or more separate widget annotations that are associated
	// with this field. However, if there is only one associated widget annotation, and its contents have been merged into the field
	// dictionary, Kids shall be omitted.
	Kids PdfObject
	T    PdfObject
	TU   PdfObject
	TM   PdfObject
	Ff   PdfObject // field flag
	V    PdfObject //value
	DV   PdfObject
	AA   PdfObject
	// Widget annotation can be merged in.
	// PdfAnnotationWidget
}

func (this *PdfField) ToPdfObject() {
	// If Kids refer only to a single pdf widget annotation widget, then can merge it in.
}

type PdfFieldVariableText struct {
	PdfField
	DA PdfObject
	Q  PdfObject
	DS PdfObject
	RV PdfObject
}

type PdfFieldText struct {
	PdfField
	// Text field
	MaxLen PdfObject // inheritable
}

type PdfFieldChoice struct {
	PdfField
	// Choice fields.
	Opt PdfObject
	TI  PdfObject
	I   PdfObject
}

type PdfFieldSignature struct {
	PdfField
	// Signature fields (Table 232).
	Lock PdfObject
	SV   PdfObject
}
