// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package role

import (
	"context"

	oupkg "github.com/thunder-id/thunderid/internal/ou"
)

// ouRoleResolverAdapter implements oupkg.OURoleResolver using the role store.
// This adapter allows the OU package to query role data without directly
// accessing the ROLE table, breaking the cross-DB access boundary.
type ouRoleResolverAdapter struct {
	store roleStoreInterface
}

// newOURoleResolver creates a new OURoleResolver backed by the given role store.
func newOURoleResolver(store roleStoreInterface) oupkg.OURoleResolver {
	return &ouRoleResolverAdapter{store: store}
}

// GetRoleCountByOUID returns the count of roles belonging to the given organization unit.
func (a *ouRoleResolverAdapter) GetRoleCountByOUID(ctx context.Context, ouID string) (int, error) {
	return a.store.GetRoleListCountByOUID(ctx, ouID)
}

// GetRoleListByOUID returns a paginated list of roles belonging to the given organization unit.
func (a *ouRoleResolverAdapter) GetRoleListByOUID(
	ctx context.Context, ouID string, limit, offset int,
) ([]oupkg.Role, error) {
	roles, err := a.store.GetRoleListByOUID(ctx, ouID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]oupkg.Role, len(roles))
	for i, r := range roles {
		result[i] = oupkg.Role{ID: r.ID, Name: r.Name, Description: r.Description, IsReadOnly: r.IsReadOnly}
	}

	return result, nil
}
