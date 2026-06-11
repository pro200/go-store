package store

import (
	"bytes"
	"testing"
)

func testKey() []byte {
	return deriveKey("test-machine-id")
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := testKey()
	plaintext := []byte("hello, go-store")

	encrypted, err := encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	decrypted, err := decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("roundtrip mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptRandomIV(t *testing.T) {
	key := testKey()
	plaintext := []byte("same plaintext")

	first, err := encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	second, err := encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	if bytes.Equal(first, second) {
		t.Fatal("same plaintext produced identical ciphertext; IV is not random")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	encrypted, err := encrypt(testKey(), []byte("secret"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	wrongKey := deriveKey("another-machine-id")
	decrypted, err := decrypt(wrongKey, encrypted)
	// CBC+PKCS7 has no authentication: a wrong key usually fails padding
	// validation, but may rarely yield garbage. It must never succeed
	// with the original plaintext and must never panic.
	if err == nil && bytes.Equal(decrypted, []byte("secret")) {
		t.Fatal("decrypt with wrong key returned original plaintext")
	}
}

func TestDecryptTooShort(t *testing.T) {
	if _, err := decrypt(testKey(), []byte("short")); err == nil {
		t.Fatal("decrypt of data shorter than IV size should fail")
	}
}

func TestDecryptTampered(t *testing.T) {
	encrypted, err := encrypt(testKey(), []byte("integrity check"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	encrypted[len(encrypted)-1] ^= 0xff
	// Must not panic; error or garbage are both acceptable for CBC.
	_, _ = decrypt(testKey(), encrypted)
}

func TestDeriveKeyStable(t *testing.T) {
	a := deriveKey("id")
	b := deriveKey("id")
	if !bytes.Equal(a, b) {
		t.Fatal("deriveKey not deterministic")
	}
	if len(a) != 32 {
		t.Fatalf("key length = %d, want 32", len(a))
	}
	if bytes.Equal(a, deriveKey("other")) {
		t.Fatal("different machine ids produced the same key")
	}
}
