// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attributecache

import (
	"errors"

	tidcommon "github.com/thunder-id/thunderid/pkg/thunderidengine/common"
)

// Store-level errors.
var (
	// errAttributeCacheNotFound is returned when an attribute cache entry is not found.
	errAttributeCacheNotFound = errors.New("attribute cache not found")
)

// Client-facing service errors.
var (
	// ErrorInvalidRequestFormat is returned when the request format is invalid.
	ErrorInvalidRequestFormat = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "ACS-1001",
		Error: tidcommon.I18nMessage{
			Key:          "error.attributecache.invalid_request_format",
			DefaultValue: "Invalid request format",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.attributecache.invalid_request_format_description",
			DefaultValue: "The request body is malformed or contains invalid data",
		},
	}

	// ErrorMissingCacheID is returned when cache ID is missing.
	ErrorMissingCacheID = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "ACS-1002",
		Error: tidcommon.I18nMessage{
			Key:          "error.attributecache.missing_cache_id",
			DefaultValue: "Missing cache ID",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.attributecache.missing_cache_id_description",
			DefaultValue: "Cache ID is required",
		},
	}

	// ErrorAttributeCacheNotFound is returned when an attribute cache entry is not found.
	ErrorAttributeCacheNotFound = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "ACS-1003",
		Error: tidcommon.I18nMessage{
			Key:          "error.attributecache.cache_not_found",
			DefaultValue: "Attribute cache not found",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.attributecache.cache_not_found_description",
			DefaultValue: "The attribute cache entry with the specified ID does not exist",
		},
	}

	// ErrorMissingAttributes is returned when attributes are missing.
	ErrorMissingAttributes = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "ACS-1004",
		Error: tidcommon.I18nMessage{
			Key:          "error.attributecache.missing_attributes",
			DefaultValue: "Missing attributes",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.attributecache.missing_attributes_description",
			DefaultValue: "Attributes are required",
		},
	}

	// ErrorInvalidExpiryTime is returned when expiry time is invalid.
	ErrorInvalidExpiryTime = tidcommon.ServiceError{
		Type: tidcommon.ClientErrorType,
		Code: "ACS-1005",
		Error: tidcommon.I18nMessage{
			Key:          "error.attributecache.invalid_expiry_time",
			DefaultValue: "Invalid expiry time",
		},
		ErrorDescription: tidcommon.I18nMessage{
			Key:          "error.attributecache.invalid_expiry_time_description",
			DefaultValue: "Expiry time must be in the future",
		},
	}
)
