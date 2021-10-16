package utils

// ConcatOrderByClause adds an ORDER BY clause to a sql query
func ConcatOrderByClause(query *string, orderBy string, descending string) {
	if orderBy != "" {
		*query += ` ORDER BY ` + orderBy

		if descending == "true" {
			*query += ` DESC`
		}
	}
}

// We let consumers pass in orderBy values that correspond to the json field they'd like to order by.
// parseOrderBy looks to match the value passed to a db field that goes in the sql query.
func ParseMusicReleaseOrderBy(orderBy string) string {
	var match string
	validOrderByKeywordsMap := map[string]string{
		"createdAt":               "created_at",
		"datePublished":           "date_published",
		"releaseOf.name":          "musicalbum.name",
		"releaseOf.byArtist.name": "musicgroup.name",
	}

	for k, v := range validOrderByKeywordsMap {
		if k == orderBy {
			match = v
		}
	}

	return match
}
