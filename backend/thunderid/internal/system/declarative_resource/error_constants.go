// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package declarativeresource

import (
	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

var (
	// ErrorDeclarativeResourceCreateOperation is the error returned when
	// a declarative resource create operation is attempted.
	ErrorDeclarativeResourceCreateOperation = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "DCR-1001",
		Error: tidcommon.I18nMessage{
			Key:          "error.declarative_resource.create_operation_not_allowed",
			DefaultValue: "Declarative resource create operation is not allowed",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.declarative_resource.create_operation_not_allowed_description",
			DefaultValue: "Creating declarative resources is not permitted",
		},
	}

	// ErrorDeclarativeResourceUpdateOperation is the error returned when
	// a declarative resource update operation is attempted.
	ErrorDeclarativeResourceUpdateOperation = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "DCR-1002",
		Error: tidcommon.I18nMessage{
			Key:          "error.declarative_resource.update_operation_not_allowed",
			DefaultValue: "Declarative resource update operation is not allowed",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.declarative_resource.update_operation_not_allowed_description",
			DefaultValue: "Updating declarative resources is not permitted",
		},
	}

	// ErrorDeclarativeResourceDeleteOperation is the error returned when
	// a declarative resource delete operation is attempted.
	ErrorDeclarativeResourceDeleteOperation = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "DCR-1003",
		Error: tidcommon.I18nMessage{
			Key:          "error.declarative_resource.delete_operation_not_allowed",
			DefaultValue: "Declarative resource delete operation is not allowed",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.declarative_resource.delete_operation_not_allowed_description",
			DefaultValue: "Deleting declarative resources is not permitted",
		},
	}
)
