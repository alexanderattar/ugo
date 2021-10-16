package models_test

// import (
// 	"database/sql/driver"
// 	"fmt"
// 	"time"

// 	. "github.com/consensys/ugo/pkg/models"
// 	"github.com/jmoiron/sqlx"
// 	. "github.com/onsi/gomega"
// 	. "github.com/onsi/gomega/gstruct"
// 	gomegatypes "github.com/onsi/gomega/types"
// 	"gopkg.in/DATA-DOG/go-sqlmock.v1"
// )

// var (
// 	musicalbums = map[IDType]map[string]driver.Value{
// 		123: {
// 			"id":                    IDType(123),
// 			"cid":                   "_cid",
// 			"type":                  "_type",
// 			"context":               "_context",
// 			"created_at":            time.Now(),
// 			"updated_at":            time.Now(),
// 			"name":                  "_name",
// 			"album_production_type": "_album_production_type",
// 			"album_release_type":    "_album_release_type",
// 			"by_artist_id":          IDType(999),
// 		},
// 	}

// 	musicrecordings = map[IDType]map[string]driver.Value{
// 		1: {
// 			"id":              IDType(1),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"name":            "_name",
// 			"duration":        "_duration",
// 			"isrc":            "_isrc",
// 			"position":        "_position",
// 			"genres":          `{"genre1", "genre2"}`,
// 			"audio_id":        IDType(456),
// 			"by_artist_id":    IDType(998),
// 			"recording_of_id": IDType(333),
// 			"image_id":        IDType(986),
// 			"visibility":      "_visibility",
// 		},
// 		2: {
// 			"id":              IDType(2),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"name":            "_name",
// 			"duration":        "_duration",
// 			"isrc":            "_isrc",
// 			"position":        "_position",
// 			"genres":          `{"genre1", "genre2"}`,
// 			"audio_id":        IDType(457),
// 			"by_artist_id":    IDType(998),
// 			"recording_of_id": IDType(334),
// 			"image_id":        IDType(987),
// 			"visibility":      "_visibility",
// 		},
// 	}

// 	musicgroups = map[IDType]map[string]driver.Value{
// 		998: {
// 			"id":          IDType(998),
// 			"cid":         "_cid",
// 			"cids":        `{"cid1", "cid2"}`,
// 			"type":        "_type",
// 			"context":     "_context",
// 			"created_at":  time.Now(),
// 			"updated_at":  time.Now(),
// 			"name":        "_name",
// 			"description": "_description",
// 			"email":       "_email",
// 			"image_id":    IDType(986),
// 		},
// 		999: {
// 			"id":          IDType(999),
// 			"cid":         "_cid",
// 			"cids":        `{"cid1", "cid2"}`,
// 			"type":        "_type",
// 			"context":     "_context",
// 			"created_at":  time.Now(),
// 			"updated_at":  time.Now(),
// 			"name":        "_name",
// 			"description": "_description",
// 			"email":       "_email",
// 			"image_id":    IDType(987),
// 		},
// 	}

// 	imageobjects = map[IDType]map[string]driver.Value{
// 		986: {
// 			"id":              IDType(986),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"content_url":     "_content_url",
// 			"encoding_format": "_encoding_format",
// 		},
// 		987: {
// 			"id":              IDType(987),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"content_url":     "_content_url",
// 			"encoding_format": "_encoding_format",
// 		},
// 	}

// 	audioobjects = map[IDType]map[string]driver.Value{
// 		456: {
// 			"id":              IDType(456),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"content_url":     "_content_url",
// 			"encoding_format": "_encoding_format",
// 		},
// 		457: {
// 			"id":              IDType(457),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"content_url":     "_content_url",
// 			"encoding_format": "_encoding_format",
// 		},
// 	}

// 	musiccompositions = map[IDType]map[string]driver.Value{
// 		333: {
// 			"id":         IDType(333),
// 			"cid":        "_cid",
// 			"type":       "_type",
// 			"context":    "_context",
// 			"created_at": time.Now(),
// 			"updated_at": time.Now(),
// 			"name":       "_name",
// 			"iswc":       "_iswc",
// 		},
// 	}

// 	people = map[IDType]map[string]driver.Value{
// 		1001: {
// 			"id":                IDType(1001),
// 			"cid":               "_cid",
// 			"type":              "_type",
// 			"context":           "_context",
// 			"created_at":        time.Now(),
// 			"updated_at":        time.Now(),
// 			"ethereum_address":  "0x1",
// 			"given_name":        "_given_name",
// 			"family_name":       "_family_name",
// 			"email":             "_email",
// 			"description":       "_description",
// 			"percentage_shares": float64(0.3),
// 			"musicgroup_admin":  true,
// 		},
// 		1002: {
// 			"id":                IDType(1002),
// 			"cid":               "_cid",
// 			"type":              "_type",
// 			"context":           "_context",
// 			"created_at":        time.Now(),
// 			"updated_at":        time.Now(),
// 			"ethereum_address":  "0x2",
// 			"given_name":        "_given_name",
// 			"family_name":       "_family_name",
// 			"email":             "_email",
// 			"description":       "_description",
// 			"percentage_shares": float64(0.29),
// 			"musicgroup_admin":  true,
// 		},
// 	}

// 	musicreleases = map[IDType]map[string]driver.Value{
// 		818: {
// 			"id":                   IDType(818),
// 			"cid":                  "_cid",
// 			"type":                 "_type",
// 			"context":              "_context",
// 			"created_at":           time.Now(),
// 			"updated_at":           time.Now(),
// 			"description":          "_description",
// 			"date_published":       "_date_published",
// 			"catalog_number":       "_catalog_number",
// 			"music_release_format": "_music_release_format",
// 			"price":                float64(1.23),
// 			"record_label_id":      IDType(-1),
// 			"release_of_id":        IDType(123),
// 			"image_id":             IDType(986),
// 			"visibility":           "_visibility",
// 			"active":               true,
// 		},
// 		717: {
// 			"id":                   IDType(717),
// 			"cid":                  "_cid",
// 			"type":                 "_type",
// 			"context":              "_context",
// 			"created_at":           time.Now(),
// 			"updated_at":           time.Now(),
// 			"description":          "_description",
// 			"date_published":       "_date_published",
// 			"catalog_number":       "_catalog_number",
// 			"music_release_format": "_music_release_format",
// 			"price":                float64(1.24),
// 			"record_label_id":      IDType(-1),
// 			"release_of_id":        IDType(123),
// 			"image_id":             IDType(987),
// 			"visibility":           "_visibility",
// 			"active":               true,
// 		},
// 	}

// 	purchases = map[IDType]map[string]driver.Value{
// 		263: {
// 			"id":              IDType(263),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"tx_hash":         "0xdeadbeef",
// 			"buyer_id":        IDType(1001),
// 			"musicrelease_id": IDType(818),
// 		},
// 		264: {
// 			"id":              IDType(264),
// 			"cid":             "_cid",
// 			"type":            "_type",
// 			"context":         "_context",
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"tx_hash":         "0xdeadbeef",
// 			"buyer_id":        IDType(1002),
// 			"musicrelease_id": IDType(717),
// 		},
// 	}

// 	reports = map[IDType]map[string]driver.Value{
// 		582: {
// 			"id":              IDType(582),
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"state":           "_state",
// 			"response":        "_response",
// 			"reason":          "_reason",
// 			"message":         "_message",
// 			"email":           "_email",
// 			"reporter_id":     IDType(41),
// 			"musicrelease_id": IDType(818),
// 		},
// 		592: {
// 			"id":              IDType(592),
// 			"created_at":      time.Now(),
// 			"updated_at":      time.Now(),
// 			"state":           "_state",
// 			"response":        "_response",
// 			"reason":          "_reason",
// 			"message":         "_message",
// 			"email":           "_email",
// 			"reporter_id":     IDType(23),
// 			"musicrelease_id": IDType(717),
// 		},
// 	}
// )

// func matchMusicgroupFixture(id IDType, memberIDs []IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := musicgroups[id]
// 	if !ok {
// 		panic(fmt.Sprintf("musicgroup fixture '%v' doesn't exist", id))
// 	}

// 	var members gomegatypes.GomegaMatcher
// 	if memberIDs == nil {
// 		members = BeNil()
// 	} else {
// 		memberMatchers := []gomegatypes.GomegaMatcher{}
// 		for _, memberID := range memberIDs {
// 			memberMatchers = append(memberMatchers, PointTo(matchPersonFixture(memberID)))
// 		}
// 		members = ConsistOf(memberMatchers)
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":          Equal(fixture["id"]),
// 		"CID":         Equal(fixture["cid"]),
// 		"CIDs":        BeAnything(), // @@TODO?
// 		"Type":        Equal(fixture["type"]),
// 		"Context":     Equal(fixture["context"]),
// 		"CreatedAt":   BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":   BeAssignableToTypeOf(time.Time{}),
// 		"Name":        PointTo(Equal(fixture["name"])),
// 		"Description": PointTo(Equal(fixture["description"])),
// 		"Email":       PointTo(Equal(fixture["email"])),
// 		"ImageID":     PointTo(Equal(fixture["image_id"])),
// 		"Members":     members,
// 		"Image":       PointTo(matchImageobjectFixture(fixture["image_id"].(IDType))),
// 	})
// }

// func matchPersonFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := people[id]
// 	if !ok {
// 		panic(fmt.Sprintf("person fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":               Equal(fixture["id"]),
// 		"CID":              Equal(fixture["cid"]),
// 		"Type":             Equal(fixture["type"]),
// 		"Context":          Equal(fixture["context"]),
// 		"CreatedAt":        BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":        BeAssignableToTypeOf(time.Time{}),
// 		"EthereumAddress":  Equal(fixture["ethereum_address"]),
// 		"GivenName":        PointTo(Equal(fixture["given_name"])),
// 		"FamilyName":       PointTo(Equal(fixture["family_name"])),
// 		"Email":            PointTo(Equal(fixture["email"])),
// 		"Image":            BeNil(),
// 		"Description":      PointTo(Equal(fixture["description"])),
// 		"PercentageShares": PointTo(Equal(fixture["percentage_shares"])),
// 		"MusicGroupAdmin":  PointTo(Equal(fixture["musicgroup_admin"])),
// 	})
// }

// func matchImageobjectFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := imageobjects[id]
// 	if !ok {
// 		panic(fmt.Sprintf("imageobject fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":             Equal(fixture["id"]),
// 		"CID":            Equal(fixture["cid"]),
// 		"Type":           Equal(fixture["type"]),
// 		"Context":        Equal(fixture["context"]),
// 		"CreatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"ContentURL":     PointTo(Equal(fixture["content_url"])),
// 		"EncodingFormat": PointTo(Equal(fixture["encoding_format"])),
// 	})
// }

// func matchMusiccompositionFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := musiccompositions[id]
// 	if !ok {
// 		panic(fmt.Sprintf("musiccomposition fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":        Equal(fixture["id"]),
// 		"CID":       Equal(fixture["cid"]),
// 		"Type":      Equal(fixture["type"]),
// 		"Context":   Equal(fixture["context"]),
// 		"CreatedAt": BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt": BeAssignableToTypeOf(time.Time{}),
// 		"Name":      PointTo(Equal(fixture["name"])),
// 		"Iswc":      PointTo(Equal(fixture["iswc"])),
// 	})
// }

// func matchAudioobjectFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := audioobjects[id]
// 	if !ok {
// 		panic(fmt.Sprintf("audioobject fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":             Equal(fixture["id"]),
// 		"CID":            Equal(fixture["cid"]),
// 		"Type":           Equal(fixture["type"]),
// 		"Context":        Equal(fixture["context"]),
// 		"CreatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"ContentURL":     PointTo(Equal(fixture["content_url"])),
// 		"EncodingFormat": PointTo(Equal(fixture["encoding_format"])),
// 	})
// }

// func matchMusicrecordingFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := musicrecordings[id]
// 	if !ok {
// 		panic(fmt.Sprintf("musicrecording fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":            Equal(fixture["id"]),
// 		"CID":           Equal(fixture["cid"]),
// 		"Type":          Equal(fixture["type"]),
// 		"Context":       Equal(fixture["context"]),
// 		"CreatedAt":     BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":     BeAssignableToTypeOf(time.Time{}),
// 		"Name":          PointTo(Equal(fixture["name"])),
// 		"RecordingOf":   BeAnything(), //BeNil(),
// 		"ByArtist":      BeAnything(), // BeNil(),
// 		"Duration":      PointTo(Equal(fixture["duration"])),
// 		"Isrc":          PointTo(Equal(fixture["isrc"])),
// 		"Position":      PointTo(Equal(fixture["position"])),
// 		"Genres":        BeAnything(), // @@TODO?
// 		"Audio":         BeAnything(), // BeNil(),
// 		"Image":         BeAnything(), // BeNil(),
// 		"AudioID":       PointTo(Equal(fixture["audio_id"])),
// 		"ByArtistID":    PointTo(Equal(fixture["by_artist_id"])),
// 		"RecordingOfID": PointTo(Equal(fixture["recording_of_id"])),
// 		"ImageID":       PointTo(Equal(fixture["image_id"])),
// 		"Visibility":    PointTo(Equal(fixture["visibility"])),
// 	})
// }

// func matchMusicreleaseFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := musicreleases[id]
// 	if !ok {
// 		panic(fmt.Sprintf("musicrelease fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":                 Equal(fixture["id"]),
// 		"CID":                Equal(fixture["cid"]),
// 		"Type":               Equal(fixture["type"]),
// 		"Context":            Equal(fixture["context"]),
// 		"CreatedAt":          BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":          BeAssignableToTypeOf(time.Time{}),
// 		"Description":        PointTo(Equal(fixture["description"])),
// 		"DatePublished":      PointTo(Equal(fixture["date_published"])),
// 		"CatalogNumber":      PointTo(Equal(fixture["catalog_number"])),
// 		"MusicReleaseFormat": PointTo(Equal(fixture["music_release_format"])),
// 		"Price":              PointTo(Equal(fixture["price"])),
// 		"RecordingLabelID":   PointTo(Equal(fixture["record_label_id"])),
// 		"ReleaseOfID":        PointTo(Equal(fixture["release_of_id"])),
// 		"ImageID":            PointTo(Equal(fixture["image_id"])),
// 		"Visibility":         PointTo(Equal(fixture["visibility"])),
// 		"Active":             Equal(fixture["active"]),
// 		"RecordLabel":        BeNil(),
// 		"ReleaseOf":          PointTo(matchMusicalbumFixture(fixture["release_of_id"].(IDType))),
// 		"Image":              PointTo(matchImageobjectFixture(fixture["image_id"].(IDType))),
// 	})
// }

// func matchMusicalbumFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := musicalbums[id]
// 	if !ok {
// 		panic(fmt.Sprintf("musicalbum fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":                  Equal(fixture["id"]),
// 		"CID":                 Equal(fixture["cid"]),
// 		"Type":                Equal(fixture["type"]),
// 		"Context":             Equal(fixture["context"]),
// 		"CreatedAt":           BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":           BeAssignableToTypeOf(time.Time{}),
// 		"Name":                PointTo(Equal(fixture["name"])),
// 		"Tracks":              BeAnything(), // @@TODO
// 		"AlbumProductionType": PointTo(Equal(fixture["album_production_type"])),
// 		"AlbumReleaseType":    PointTo(Equal(fixture["album_release_type"])),
// 		"ByArtist":            BeAnything(), // @@TODO
// 		"ByArtistID":          PointTo(Equal(fixture["by_artist_id"])),
// 	})
// }

// func matchPurchaseFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := purchases[id]
// 	if !ok {
// 		panic(fmt.Sprintf("purchase fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":             Equal(fixture["id"]),
// 		"CID":            Equal(fixture["cid"]),
// 		"Type":           Equal(fixture["type"]),
// 		"Context":        Equal(fixture["context"]),
// 		"CreatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"TxHash":         Equal(fixture["tx_hash"]),
// 		"BuyerID":        PointTo(Equal(fixture["buyer_id"])),
// 		"Buyer":          PointTo(matchPersonFixture(fixture["buyer_id"].(IDType))),
// 		"MusicReleaseID": PointTo(Equal(fixture["musicrelease_id"])),
// 		"MusicRelease":   PointTo(matchMusicreleaseFixture(fixture["musicrelease_id"].(IDType))),
// 	})
// }

// func matchReportFixture(id IDType) gomegatypes.GomegaMatcher {
// 	fixture, ok := reports[id]
// 	if !ok {
// 		panic(fmt.Sprintf("report fixture '%v' doesn't exist", id))
// 	}

// 	return MatchFields(0, Fields{
// 		"ID":             Equal(fixture["id"]),
// 		"CreatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"UpdatedAt":      BeAssignableToTypeOf(time.Time{}),
// 		"State":          Equal(fixture["state"]),
// 		"Response":       PointTo(Equal(fixture["response"])),
// 		"Reason":         Equal(fixture["reason"]),
// 		"Message":        PointTo(Equal(fixture["message"])),
// 		"Email":          PointTo(Equal(fixture["email"])),
// 		"ReporterID":     Equal(fixture["reporter_id"]),
// 		"MusicRelease":   PointTo(matchMusicreleaseFixture(fixture["musicrelease_id"].(IDType))),
// 		"MusicReleaseID": Equal(fixture["musicrelease_id"]),
// 	})
// }

// func mockResultRows(rows ...map[string]driver.Value) *sqlmock.Rows {
// 	cols := []string{}

// 	if len(rows) == 0 {
// 		return sqlmock.NewRows(cols)
// 	}

// 	for col := range rows[0] {
// 		cols = append(cols, col)
// 	}

// 	result := sqlmock.NewRows(cols)
// 	for i := range rows {
// 		result = result.AddRow(values(cols, rows[i])...)
// 	}
// 	return result
// }

// func values(cols []string, obj map[string]driver.Value) []driver.Value {
// 	vals := make([]driver.Value, len(cols))
// 	for i, col := range cols {
// 		vals[i] = obj[col]
// 	}
// 	return vals
// }

// func musicrelease_GetByID_simulate(id IDType, musicrecordingIDs, peopleIDs []IDType) (*MusicRelease, error) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	dbx := sqlx.NewDb(db, "postgres")
// 	defer dbx.Close()

// 	expectQuery_MusicRelease_GetByID(mock, id, musicrecordingIDs, peopleIDs)
// 	return (&MusicRelease{}).GetByID(dbx, id)
// }
