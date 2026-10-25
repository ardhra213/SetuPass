// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package utils

import "fmt"

// QueryParamInclude is the query parameter name for the include parameter.
const QueryParamInclude = "include"

// IncludeValueDisplay is the value for the include query parameter to request display attributes.
const IncludeValueDisplay = "display"

// IncludeDisplayQuery is the query string fragment appended to pagination links
// when the include=display parameter is active.
const IncludeDisplayQuery = "&" + QueryParamInclude + "=" + IncludeValueDisplay

// DisplayQueryParam returns IncludeDisplayQuery if includeDisplay is true, empty string otherwise.
func DisplayQueryParam(includeDisplay bool) string {
	if includeDisplay {
		return IncludeDisplayQuery
	}
	return ""
}

// Link represents a pagination link in API responses.
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// BuildPaginationLinks builds pagination links for paginated API responses.
// extraQuery is an optional query string fragment (e.g. "&include=display") appended to each link.
func BuildPaginationLinks(base string, limit, offset, totalCount int, extraQuery string) []Link {
	links := make([]Link, 0)

	if limit <= 0 {
		return links
	}

	if offset > 0 {
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=0&limit=%d%s", base, limit, extraQuery),
			Rel:  "first",
		})

		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=%d&limit=%d%s", base, prevOffset, limit, extraQuery),
			Rel:  "prev",
		})
	}

	if offset+limit < totalCount {
		nextOffset := offset + limit
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=%d&limit=%d%s", base, nextOffset, limit, extraQuery),
			Rel:  "next",
		})
	}

	lastPageOffset := ((totalCount - 1) / limit) * limit
	if offset < lastPageOffset {
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=%d&limit=%d%s", base, lastPageOffset, limit, extraQuery),
			Rel:  "last",
		})
	}

	return links
}
