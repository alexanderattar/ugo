package utils

import "testing"

func TestConcatOrderByClause(t *testing.T) {
	query := ""
	orderBy := "created_at"
	descending := ""
	ConcatOrderByClause(&query, orderBy, descending)

	if query != " ORDER BY created_at" {
		t.Fatalf("expected \" ORDER BY created_at\" but got: %v", query)
	}

	query = ""
	descending = "true"
	ConcatOrderByClause(&query, orderBy, descending)

	if query != " ORDER BY created_at DESC" {
		t.Fatalf("expected \" ORDER BY created_at DESC\" but got: %v", query)
	}
}

func TestParseMusicReleaseOrderBy(t *testing.T) {
	//some test conditions
	returned := ParseMusicReleaseOrderBy("createdAt")
	if returned != "created_at" {
		t.Errorf("expected \"created_at\" but got: %v", returned)
	}

	returned = ParseMusicReleaseOrderBy("evil things")
	if returned != "" {
		t.Fatalf("expected an empty string but got: %v", returned)
	}
}
