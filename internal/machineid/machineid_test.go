package machineid

import "testing"

func TestID(t *testing.T) {
	id1, err := ID()
	if err != nil {
		t.Fatalf("ID() error: %v", err)
	}
	if id1 == "" {
		t.Fatal("ID() returned empty string without error")
	}

	id2, err := ID()
	if err != nil {
		t.Fatalf("ID() second call error: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("ID() not stable: %q != %q", id1, id2)
	}
}

func TestFirstMAC(t *testing.T) {
	mac, err := firstMAC()
	if err == nil && mac == "" {
		t.Fatal("firstMAC() returned empty string without error")
	}
	// Either a MAC or an error is acceptable depending on the host;
	// the call must simply never panic and never return both empty.
}

func TestPlatformID(t *testing.T) {
	id, err := platformID()
	if err != nil {
		t.Skipf("platformID unavailable on this host: %v", err)
	}
	if id == "" {
		t.Fatal("platformID() returned empty string without error")
	}
}
