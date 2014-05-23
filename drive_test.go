// Copyright 2014 ALTOROS
// Licensed under the AGPLv3, see LICENSE file for details.

package gosigma

import (
	"testing"

	"github.com/Altoros/gosigma/data"
	"github.com/Altoros/gosigma/mock"
)

func newDataDrive(uuid string) *data.Drive {
	result := &data.Drive{
		Resource: data.Resource{URI: "uri", UUID: uuid},
		LibraryDrive: data.LibraryDrive{
			Arch:      "arch",
			ImageType: "image-type",
			OS:        "os",
			Paid:      true,
		},
		Affinities:      []string{"aff1", "aff2"},
		AllowMultimount: true,
		Jobs: []data.Resource{
			data.Resource{URI: "uri-0", UUID: "uuid-0"},
			data.Resource{URI: "uri-1", UUID: "uuid-1"},
		},
		Media:       "media",
		Meta:        map[string]string{"key1": "value1", "key2": "value2"},
		Name:        "name",
		Owner:       &data.Resource{URI: "owner-uri", UUID: "owner-uuid"},
		Size:        1000,
		Status:      "status",
		StorageType: "storage-type",
	}
	return result
}

func testDrive(t *testing.T, d Drive, uuid string, short bool) {
	checkv := func(v, wants string) {
		if v != wants {
			t.Errorf("value %q, wants %q", v, wants)
		}
	}
	checkb := func(v, wants bool) {
		if v != wants {
			t.Errorf("value %s, wants %s", v, wants)
		}
	}
	checki := func(v, wants uint64) {
		if v != wants {
			t.Errorf("value %s, wants %s", v, wants)
		}
	}
	checkg := func(d Drive, k, wants string) {
		if v, ok := d.Get(k); !ok || v != wants {
			t.Errorf("value of Get(%q) = %q, %v; wants %s", k, v, ok, wants)
		}
	}

	checkv(d.URI(), "uri")
	checkv(d.UUID(), uuid)
	checkv(d.Owner().URI(), "owner-uri")
	checkv(d.Owner().UUID(), "owner-uuid")
	checkv(d.Status(), "status")

	if short {
		checkv(d.Name(), "")
		checkv(d.Media(), "")
		checkv(d.StorageType(), "")
		return
	}

	checkv(d.Arch(), "arch")
	checkv(d.ImageType(), "image-type")
	checkv(d.OS(), "os")
	checkb(d.Paid(), true)

	if v := d.Affinities(); len(v) != 2 {
		t.Error("Affinities failed:", len(v), v)
	} else {
		checkv(v[0], "aff1")
		checkv(v[1], "aff2")
	}

	checkb(d.AllowMultimount(), true)

	if v := d.Jobs(); len(v) != 2 {
		t.Error("Jobs failed:", v)
	} else {
		checkv(v[0].URI(), "uri-0")
		checkv(v[0].UUID(), "uuid-0")
		checkv(v[1].URI(), "uri-1")
		checkv(v[1].UUID(), "uuid-1")
	}

	checkv(d.Media(), "media")

	checkg(d, "key1", "value1")
	checkg(d, "key2", "value2")
	checkv(d.Name(), "name")
	checki(d.Size(), 1000)
	checkv(d.StorageType(), "storage-type")
}

func TestClientDriveEmpty(t *testing.T) {
	d := &drive{obj: &data.Drive{}}
	if d.Owner() != nil {
		t.Error("invalid owner")
	}
}

func TestClientDrivesEmpty(t *testing.T) {
	mock.Drives.ResetDrives()

	cli, err := createTestClient(t)
	if err != nil || cli == nil {
		t.Error("NewClient() failed:", err, cli)
		return
	}

	check := func(rqspec RequestSpec, libspec LibrarySpec) {
		drives, err := cli.Drives(rqspec, libspec)
		if err != nil {
			t.Error(err)
		}
		if len(drives) > 0 {
			t.Errorf("%v", drives)
		}
	}

	check(RequestShort, LibraryAccount)
	check(RequestDetail, LibraryAccount)
}

func TestClientDrives(t *testing.T) {
	mock.Drives.ResetDrives()

	mock.Drives.AddDrive(newDataDrive("uuid-0"))
	mock.Drives.AddDrive(newDataDrive("uuid-1"))

	cli, err := createTestClient(t)
	if err != nil {
		t.Error(err)
		return
	}

	drives, err := cli.Drives(RequestShort, LibraryAccount)
	if err != nil {
		t.Error(err)
		return
	}

	if len(drives) != 2 {
		t.Errorf("invalid len: %v", drives)
		return
	}

	for i, uuid := range []string{"uuid-0", "uuid-1"} {
		d := drives[i]

		if d.String() == "" {
			t.Error("Empty string representation")
			return
		}

		testDrive(t, d, uuid, true)

		// refresh
		if err := d.Refresh(); err != nil {
			t.Error(err)
			return
		}

		testDrive(t, d, uuid, false)

		if err := d.Remove(); err != nil {
			t.Error("Drive remove fail:", err)
			return
		}

		// failed refresh
		if err := d.Refresh(); err == nil {
			t.Error("Drive refresh must fail")
			return
		}
	}

	mock.Drives.ResetDrives()
}

func TestClientDrive(t *testing.T) {
	mock.Drives.ResetDrives()

	mock.Drives.AddDrive(newDataDrive("uuid"))

	cli, err := createTestClient(t)
	if err != nil {
		t.Error(err)
		return
	}

	d, err := cli.Drive("uuid", LibraryAccount)
	if err != nil {
		t.Error(err)
		return
	}

	if d.String() == "" {
		t.Error("Empty string representation")
		return
	}

	testDrive(t, d, "uuid", false)

	if err := d.Remove(); err != nil {
		t.Error("Drive remove fail:", err)
		return
	}

	// failed refresh
	if err := d.Refresh(); err == nil {
		t.Error("Drive refresh must fail")
		return
	}

	mock.Drives.ResetDrives()
}
