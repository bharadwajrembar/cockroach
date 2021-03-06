// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package tree

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/timeofday"
	"github.com/stretchr/testify/assert"
)

func TestAllTypesCastableToString(t *testing.T) {
	defer leaktest.AfterTest(t)()
	for _, typ := range types.Scalar {
		if ok, _ := isCastDeepValid(typ, types.String); !ok {
			t.Errorf("%s is not castable to STRING, all types should be", typ)
		}
	}
}

func TestAllTypesCastableFromString(t *testing.T) {
	defer leaktest.AfterTest(t)()
	for _, typ := range types.Scalar {
		if ok, _ := isCastDeepValid(types.String, typ); !ok {
			t.Errorf("%s is not castable from STRING, all types should be", typ)
		}
	}
}

func TestCompareTimestamps(t *testing.T) {
	defer leaktest.AfterTest(t)()

	pacificTimeZone := int32(7 * 60 * 60)
	sydneyTimeZone := int32(-10 * 60 * 60)

	testCases := []struct {
		desc     string
		left     Datum
		right    Datum
		expected int
	}{
		{
			desc:     "same DTime are equal",
			left:     MakeDTime(timeofday.New(12, 0, 0, 0)),
			right:    MakeDTime(timeofday.New(12, 0, 0, 0)),
			expected: 0,
		},
		{
			desc:     "same DTimeTZ are equal",
			left:     NewDTimeTZFromOffset(timeofday.New(22, 0, 0, 0), sydneyTimeZone),
			right:    NewDTimeTZFromOffset(timeofday.New(22, 0, 0, 0), sydneyTimeZone),
			expected: 0,
		},
		{
			desc:     "equal wall clock time for DTime and DTimeTZ, with TimeTZ ahead",
			left:     MakeDTime(timeofday.New(12, 0, 0, 0)),
			right:    NewDTimeTZFromOffset(timeofday.New(22, 0, 0, 0), sydneyTimeZone),
			expected: 1,
		},
		{
			desc:     "equal wall clock time for DTime and DTimeTZ, with TimeTZ behind",
			left:     MakeDTime(timeofday.New(12, 0, 0, 0)),
			right:    NewDTimeTZFromOffset(timeofday.New(5, 0, 0, 0), pacificTimeZone),
			expected: -1,
		},
		{
			desc:     "equal wall clock time for DTime and DTimeTZ, with TimeTZ ahead",
			left:     NewDTimeTZFromOffset(timeofday.New(22, 0, 0, 0), sydneyTimeZone),
			right:    NewDTimeTZFromOffset(timeofday.New(5, 0, 0, 0), pacificTimeZone),
			expected: -1,
		},
		{
			desc:     "wall clock time different for DTimeTZ and DTimeTZ",
			left:     NewDTimeTZFromOffset(timeofday.New(23, 0, 0, 0), sydneyTimeZone),
			right:    NewDTimeTZFromOffset(timeofday.New(5, 0, 0, 0), pacificTimeZone),
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("%s cmp %s", tc.left.String(), tc.right.String()),
			func(t *testing.T) {
				assert.Equal(nil, tc.expected, compareTimestamps(nil, tc.left, tc.right))
				assert.Equal(nil, -tc.expected, compareTimestamps(nil, tc.right, tc.left))
			},
		)
	}
}
