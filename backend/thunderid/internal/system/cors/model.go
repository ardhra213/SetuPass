// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cors

// entry is the discriminated-union type produced by YAML decoding. It carries
// either a literal allowed-origin string or a regex pattern. compile turns an
// entry into the corresponding compiled originRule. The interface and its
// implementations are unexported; callers outside the cors package construct
// OriginEntries values exclusively via YAML decoding through UnmarshalYAML.
type entry interface {
	isOriginEntry()
}

// literalEntry is the bare-string YAML form, e.g. "https://example.com" or
// the special-case "null".
type literalEntry struct {
	Value string
}

// isOriginEntry marks literalEntry as a member of the entry sealed-interface
// union. The method has no behavior and is never called directly; its sole
// purpose is to restrict entry implementations to types declared in this
// package.
func (literalEntry) isOriginEntry() {}

// regexEntry is the object YAML form, e.g. { regex: "\\Ahttps://..." }.
type regexEntry struct {
	Pattern string
}

// isOriginEntry marks regexEntry as a member of the entry sealed-interface
// union. The method has no behavior and is never called directly; its sole
// purpose is to restrict entry implementations to types declared in this
// package.
func (regexEntry) isOriginEntry() {}

// OriginEntries is the slice wrapper that carries the heterogeneous YAML
// schema for cors.allowed_origins. The element type is unexported, so callers
// outside the cors package construct values exclusively through YAML
// decoding. Custom YAML unmarshaling on this type dispatches between the two
// entry forms.
type OriginEntries []entry
