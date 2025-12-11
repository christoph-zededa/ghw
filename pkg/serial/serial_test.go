//go:build linux
// +build linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package serial

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaypipes/ghw/pkg/option"
)

func TestSerial(t *testing.T) {
	root, err := os.MkdirTemp("", "ghw-serial-test-")
	if err != nil {
		t.Fatalf("could not create temp directory: %v", err)
	}
	defer os.RemoveAll(root)

	// Create /sys/class/tty
	ttyDir := filepath.Join(root, "sys", "class", "tty")
	if err := os.MkdirAll(ttyDir, 0755); err != nil {
		t.Fatalf("could not create tty directory: %v", err)
	}

	// ttyS0 - valid
	ttyS0 := filepath.Join(ttyDir, "ttyS0")
	if err := os.MkdirAll(filepath.Join(ttyS0, "device"), 0755); err != nil {
		t.Fatalf("could not create ttyS0 directory: %v", err)
	}
	resourcesS0 := "io 0x3f8-0x3ff\nirq 4\n"
	if err := os.WriteFile(filepath.Join(ttyS0, "device", "resources"), []byte(resourcesS0), 0644); err != nil {
		t.Fatalf("could not write ttyS0 resources: %v", err)
	}

	// ttyS1 - no resources file
	ttyS1 := filepath.Join(ttyDir, "ttyS1")
	if err := os.MkdirAll(ttyS1, 0755); err != nil {
		t.Fatalf("could not create ttyS1 directory: %v", err)
	}

	// ttyS2 - empty resources
	ttyS2 := filepath.Join(ttyDir, "ttyS2")
	if err := os.MkdirAll(filepath.Join(ttyS2, "device"), 0755); err != nil {
		t.Fatalf("could not create ttyS2 directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(ttyS2, "device", "resources"), []byte(""), 0644); err != nil {
		t.Fatalf("could not write ttyS2 resources: %v", err)
	}

	info, err := New(option.WithChroot(root))
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if len(info.Devices) != 1 {
		t.Errorf("expected 1 device, got %d", len(info.Devices))
	}

	if len(info.Devices) > 0 {
		dev := info.Devices[0]
		if dev.Name != "COM1" {
			t.Errorf("expected name COM1, got %s", dev.Name)
		}
		if dev.Address != "/dev/ttyS0" {
			t.Errorf("expected address /dev/ttyS0, got %s", dev.Address)
		}
		if dev.IO != "3f8-3ff" {
			t.Errorf("expected IO 3f8-3ff, got %s", dev.IO)
		}
		if dev.IRQ != "4" {
			t.Errorf("expected IRQ 4, got %s", dev.IRQ)
		}
	}
}
