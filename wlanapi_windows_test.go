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
	if unsafe.Offsetof(InterfaceInfo{}.InterfaceGUID) != 0 {
		t.Errorf("InterfaceInfo.InterfaceGUID wrong offset: %v", unsafe.Offsetof(InterfaceInfo{}.InterfaceGUID))
	}
	if unsafe.Offsetof(InterfaceInfo{}.interfaceDescription) != 16 {
		t.Errorf("InterfaceInfo.interfaceDescription wrong offset: %v", unsafe.Offsetof(InterfaceInfo{}.interfaceDescription))
	}
	if unsafe.Offsetof(InterfaceInfo{}.State) != 528 {
		t.Errorf("InterfaceInfo.State wrong offset: %v", unsafe.Offsetof(InterfaceInfo{}.State))
	}
	if unsafe.Sizeof(InterfaceInfo{}) != 532 {
		t.Errorf("InterfaceInfo wrong size: %v", unsafe.Sizeof(InterfaceInfo{}))
	}

	if unsafe.Offsetof(InterfaceInfoList{}.NumberOfItems) != 0 {
		t.Errorf("InterfaceInfoList.NumberOfItems wrong offset: %v", unsafe.Offsetof(InterfaceInfoList{}.NumberOfItems))
	}
	if unsafe.Offsetof(InterfaceInfoList{}.Index) != 4 {
		t.Errorf("InterfaceInfoList.Index wrong offset: %v", unsafe.Offsetof(InterfaceInfoList{}.Index))
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

	for _, ii := range ifaces.InterfaceInfo() {
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
