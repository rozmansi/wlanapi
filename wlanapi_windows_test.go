/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2020 Simon Rozman <simon@rozman.si>. All Rights Reserved.
 */

package wlanapi

import (
	"testing"

	"golang.org/x/sys/windows"
)

func Test(t *testing.T) {
	session, version, err := CreateClientSession(2)
	if err != nil {
		t.Errorf("Error creating client session: %v", err)
	}
	defer session.Close()
	if version < 2 {
		t.Errorf("Invalid version: %v", version)
	}

	ifaces, err := session.Interfaces()
	if err != nil {
		t.Errorf("Error enumerating interfaces: %v", err)
	}
	defer ifaces.Close()

	for i := uint32(0); i < ifaces.NumberOfItems; i++ {
		ii := ifaces.Item(i)
		if ii.State == InterfaceStateNotReady {
			continue
		}

		desc := ii.InterfaceDescription()
		t.Logf("Interface: %v", desc)

		err = session.SetProfileEAPXMLUserData(&ii.InterfaceGUID, "foobar", 0, "<foobar></foobar>")
		if err == nil {
			t.Errorf("SetProfileEAPXMLUserData error expected")
		}
		if err != windows.E_NOT_SET {
			t.Errorf("SetProfileEAPXMLUserData returned other error than expected: %v", err)
		}
	}
}
