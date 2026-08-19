package utils

import "testing"

func TestNormalizeTin_stripsNonDigits(t *testing.T) {
	got := NormalizeTin(" 123-456 789 ")
	if got != "123456789" {
		t.Fatalf("NormalizeTin() = %q, want %q", got, "123456789")
	}
}

func TestIsValidOrganizationTin_acceptsNineDigits(t *testing.T) {
	if !IsValidOrganizationTin("123-456-789") {
		t.Fatal("expected formatted 9-digit TIN to be valid")
	}
}

func TestIsValidOrganizationTin_rejectsWrongLength(t *testing.T) {
	if IsValidOrganizationTin("12345678") {
		t.Fatal("expected 8-digit TIN to be invalid")
	}
	if IsValidOrganizationTin("1234567890") {
		t.Fatal("expected 10-digit TIN to be invalid")
	}
	if IsValidOrganizationTin("") {
		t.Fatal("expected empty TIN to be invalid")
	}
}
