// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetHeaders(t *testing.T) {
	testCases := []struct {
		realmRoles     []string
		resourceRoles  []string
		expectedResult string
		expectedError  error
	}{
		// Global role
		{
			resourceRoles:  []string{"admin"},
			expectedResult: "project1|project2|project3|edgenode-system",
			expectedError:  nil,
		},
		// Global role but also restricted viewer role present - still an admin
		{
			resourceRoles:  []string{"admin"},
			realmRoles:     []string{"project1_tc-r"},
			expectedResult: "project1|project2|project3|edgenode-system",
			expectedError:  nil,
		},
		// Restricted viewer role
		{
			realmRoles:     []string{"project1_tc-r", "project2_tc-r"},
			expectedResult: "project1|project2",
			expectedError:  nil,
		},
		// Unknown project role
		{
			realmRoles:     []string{"projectunknown_tc-r"},
			expectedResult: "",
			expectedError:  errNoProjectAccess,
		},
		// Unknown role
		{
			resourceRoles:  []string{"administrator"},
			expectedResult: "",
			expectedError:  errNoProjectAccess,
		},
	}

	p := &proxy{
		projects:              []string{"project1", "project2", "project3"},
		includeEdgenodeSystem: true,
	}

	for _, tc := range testCases {
		result, err := p.setHeaders(tc.realmRoles, tc.resourceRoles)
		require.Equal(t, tc.expectedResult, result)
		require.Equal(t, tc.expectedError, err)
	}
}
