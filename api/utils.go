package api

import "github.com/google/uuid"

func convertUUIDToString(uuid *uuid.UUID) string {
	if uuid == nil {
		return ""
	}

	return uuid.String()
}

func convertStringToUUID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}