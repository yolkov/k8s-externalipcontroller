// Copyright 2016 Mirantis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package extensions

import (
	"testing"
)

func checkFreeIP(t *testing.T, p *IpClaimPool, expected string) {
	actual, err := p.AvailableIP()
	if err != nil {
		t.Errorf("Error must not occur during AvailableIP() method; details --> %v\n", err)
	}

	if expected != actual {
		t.Errorf("Actual free IP %v is not as expected %v\n", actual, expected)
	}
}

func checkNoFreeIPError(t *testing.T, p *IpClaimPool) {
	ip, err := p.AvailableIP()

	if len(ip) != 0 {
		t.Error("AvailableIP must return empty string as IP value in case there is no free IP")
	}
	if err == nil {
		t.Error("AvailableIP must return error in case there is no free IP")
	}
	expectedErrMsg := "There is no free IP left in the pool"
	if err.Error() != expectedErrMsg {
		t.Errorf("Message of the returned by AvailableIP error is not as expected ('%v'). Actual: %v", expectedErrMsg, err)
	}
}

func TestIpClaimPoolAvailableIPFromCIDR(t *testing.T) {
	cidr := "192.168.16.248/29"
	allocated := map[string]string{
		"192.168.16.249": "test-claim-249",
		"192.168.16.250": "test-claim-250",
		"192.168.16.252": "test-claim-252",
		"192.168.16.253": "test-claim-253",
	}
	expectedIP := "192.168.16.251"

	ClaimPool := &IpClaimPool{
		Spec: IpClaimPoolSpec{
			CIDR:      cidr,
			Ranges:    nil,
			Allocated: allocated,
		},
	}

	checkFreeIP(t, ClaimPool, expectedIP)

	allocated["192.168.16.254"] = "test-claim-254"
	allocated["192.168.16.251"] = "test-claim-251"

	checkNoFreeIPError(t, ClaimPool)
}

func TestIpClaimPoolAvailableIPFromRange(t *testing.T) {
	cidr := "192.168.16.248/29"
	IPranges := [][]string{[]string{"192.168.16.250", "192.168.16.252"}}

	allocated := map[string]string{
		"192.168.16.250": "test-claim-250",
		"192.168.16.251": "test-claim-251",
	}

	expectedIP := "192.168.16.252"

	ClaimPool := &IpClaimPool{
		Spec: IpClaimPoolSpec{
			CIDR:      cidr,
			Ranges:    IPranges,
			Allocated: allocated,
		},
	}

	checkFreeIP(t, ClaimPool, expectedIP)

	allocated["192.168.16.252"] = "test-claim-252"

	checkNoFreeIPError(t, ClaimPool)
}

func TestIpClaimPoolRangesProperlyProcessed(t *testing.T) {
	cidr := "192.168.16.248/29"
	IPranges := [][]string{
		[]string{"192.168.16.249", "192.168.16.250"},
		[]string{"192.168.16.252", "192.168.16.253"},
	}

	allocated := map[string]string{
		"192.168.16.249": "test-claim-249",
		"192.168.16.250": "test-claim-250",
		"192.168.16.252": "test-claim-252",
	}

	expectedIP := "192.168.16.253"

	ClaimPool := &IpClaimPool{
		Spec: IpClaimPoolSpec{
			CIDR:      cidr,
			Ranges:    IPranges,
			Allocated: allocated,
		},
	}

	checkFreeIP(t, ClaimPool, expectedIP)

	allocated["192.168.16.253"] = "test-claim-253"

	checkNoFreeIPError(t, ClaimPool)
}
