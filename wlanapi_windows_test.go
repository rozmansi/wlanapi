/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2020 Simon Rozman <simon@rozman.si>. All Rights Reserved.
 */

package wlanapi

import (
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestStruct(t *testing.T) {
	if unsafe.Sizeof(InterfaceInfo{}) != 532 {
		t.Errorf("InterfaceInfo wrong size: %v", unsafe.Sizeof(InterfaceInfo{}))
	}
	if unsafe.Sizeof(InterfaceInfoList{}) != 8 {
		t.Errorf("InterfaceInfoList wrong size: %v", unsafe.Sizeof(InterfaceInfoList{}))
	}
}

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

		t.Logf("Interface: %v, state: %v, GUID: %v",
			ii.InterfaceDescription(),
			ii.State,
			ii.InterfaceGUID)

		if ii.State == InterfaceStateNotReady {
			continue
		}

		err = session.SetProfileEAPXMLUserData(&ii.InterfaceGUID, "foobar", 0, "<foobar></foobar>")
		if err == nil {
			t.Errorf("SetProfileEAPXMLUserData error expected")
		}
		if err != windows.E_NOT_SET {
			t.Errorf("SetProfileEAPXMLUserData returned other error than expected: %v", err)
		}
	}
}
