// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cert

import (
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/database/provider"
)

// Initialize initializes and returns the certificate service.
func Initialize(cacheManager cache.CacheManagerInterface, dbProvider provider.DBProviderInterface) (
	CertificateServiceInterface, error) {
	txn, err := dbProvider.GetConfigDBTransactioner()
	if err != nil {
		return nil, err
	}
	certByIDCache := cache.GetCache[*Certificate](cacheManager, "CertificateByIDCache")
	certByReferenceCache := cache.GetCache[*Certificate](cacheManager, "CertificateByReferenceCache")
	certStore := newCachedBackedCertificateStore(certByIDCache, certByReferenceCache)
	certService := newCertificateService(certStore, txn)
	return certService, nil
}
