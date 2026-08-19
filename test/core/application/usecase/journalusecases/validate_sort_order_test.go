package journalusecases_test

import (
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/journalusecases"
)

// asFailResponse xatolikni *response.Response turiga keltiradi yoki testni
// muvaffaqiyatsiz deb belgilaydi.
func asFailResponse(t *testing.T, err error) *response.Response {
	t.Helper()

	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("*response.Response kutilgandi, %T (%v) keldi", err, err)
	}
	return resp
}

func TestValidateSortOrderEmptyBothReturnsNoError(t *testing.T) {
	sortBy, order, err := journalusecases.ValidateSortOrder("", "")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if sortBy != "" || order != "" {
		t.Errorf("bo'sh qiymatlar kutilgandi, sortBy=%q order=%q keldi", sortBy, order)
	}
}

func TestValidateSortOrderEmptySortByWithOrderIsRejected(t *testing.T) {
	_, _, err := journalusecases.ValidateSortOrder("", "desc")
	resp := asFailResponse(t, err)
	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
}

func TestValidateSortOrderEmptyOrderDefaultsToDesc(t *testing.T) {
	sortBy, order, err := journalusecases.ValidateSortOrder("views_count", "")
	if err != nil {
		t.Fatalf("xatolik kutilmagandi: %v", err)
	}
	if sortBy != "views_count" {
		t.Errorf("sortBy %q kutilgandi, %q keldi", "views_count", sortBy)
	}
	if order != "desc" {
		t.Errorf("standart yo'nalish %q kutilgandi, %q keldi", "desc", order)
	}
}

func TestValidateSortOrderNormalizesDirection(t *testing.T) {
	cases := map[string]string{
		"asc":  "asc",
		"ASC":  "asc",
		"desc": "desc",
		"DESC": "desc",
	}
	for input, want := range cases {
		_, order, err := journalusecases.ValidateSortOrder("views_count", input)
		if err != nil {
			t.Errorf("%q: xatolik kutilmagandi: %v", input, err)
			continue
		}
		if order != want {
			t.Errorf("%q: %q kutilgandi, %q keldi", input, want, order)
		}
	}
}

func TestValidateSortOrderRejectsInvalidDirection(t *testing.T) {
	_, _, err := journalusecases.ValidateSortOrder("views_count", "asc; DROP TABLE journals")
	resp := asFailResponse(t, err)
	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
}
