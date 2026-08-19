package validation

import (
	"fmt"

	"slib.uz/src/core/application/response"
)

func ValidateIDs(field string, ids []uint, getExistingIDs func([]uint) ([]uint, error)) error {
	if len(ids) == 0 {
		return nil
	}

	existing, err := getExistingIDs(ids)
	if err != nil {
		return err
	}
	if len(existing) == len(ids) {
		return nil
	}

	existingSet := make(map[uint]bool, len(existing))
	for _, id := range existing {
		existingSet[id] = true
	}

	var missing []uint
	for _, id := range ids {
		if !existingSet[id] {
			missing = append(missing, id)
		}
	}
	return response.NewFailResponse(400, fmt.Sprintf("Invalid %s IDs: %v", field, missing))
}
