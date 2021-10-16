package models_test

// import (
// 	"database/sql"
// 	"database/sql/driver"

// 	"github.com/jmoiron/sqlx"
// 	. "github.com/onsi/ginkgo"
// 	"gopkg.in/DATA-DOG/go-sqlmock.v1"

// 	. "github.com/consensys/ugo/pkg/models"
// )

// var _ = Describe("MusicAlbum", func() {
// 	var dbx *sqlx.DB
// 	var mock sqlmock.Sqlmock

// 	BeforeEach(func() {
// 		var db *sql.DB
// 		var err error
// 		db, mock, err = sqlmock.New()
// 		if err != nil {
// 			Fail(err.Error())
// 		}

// 		dbx = sqlx.NewDb(db, "postgres")
// 	})

// 	AfterEach(func() { dbx.Close() })

// 	Context("when .All is called", func() {
// 		It("should return an error about not being implemented yet", func() {
// 			_, err := (&MusicAlbum{}).All(dbx)
// 			if err == nil || err.Error() != "Not implemented" {
// 				Fail("expected to receive an error saying 'Not implemented'")
// 			}
// 		})
// 	})

// 	Context("when .Get is called", func() {
// 		It("should do some SQL things", func() {
// 			var musicalbumID IDType = 123
// 			var musicrecordingIDs = []IDType{1, 2}
// 			var peopleIDs = []IDType{1001, 1002}

// 			expectQuery_MusicAlbum_Get(mock, musicalbumID, musicrecordingIDs, peopleIDs)

// 			_, err := (&MusicAlbum{}).Get(dbx, musicalbumID)
// 			if err != nil {
// 				Fail(err.Error())
// 			}

// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .GetTracks is called", func() {
// 		It("should do some SQL things", func() {

// 			var musicalbumID IDType = 123
// 			var musicrecordingIDs = []IDType{1, 2}
// 			expectQuery_MusicAlbum_GetTracks(mock, musicalbumID, musicrecordingIDs)

// 			_, err := (&MusicAlbum{}).GetTracks(dbx, musicalbumID)
// 			if err != nil {
// 				Fail(err.Error())
// 			}

// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .GetArtist is called", func() {
// 		It("should do some SQL things", func() {

// 			var musicgroupID IDType = 999
// 			expectQuery_MusicAlbum_GetArtist(mock, musicgroupID)

// 			_, err := (&MusicAlbum{}).GetArtist(dbx, musicgroupID)
// 			if err != nil {
// 				Fail(err.Error())
// 			}

// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})
// })

// func expectQuery_MusicAlbum_Get(mock sqlmock.Sqlmock, musicalbumID IDType, musicrecordingIDs []IDType, peopleIDs []IDType) {
// 	mock.ExpectQuery(`SELECT (.+) FROM musicalbum`).
// 		WithArgs(musicalbumID).
// 		WillReturnRows(
// 			mockResultRows(
// 				musicalbums[musicalbumID],
// 			),
// 		)

// 	var musicgroupID = musicalbums[musicalbumID]["by_artist_id"].(IDType)

// 	expectQuery_MusicAlbum_GetTracks(mock, musicalbumID, musicrecordingIDs)
// 	expectQuery_MusicAlbum_GetArtist(mock, musicgroupID)
// 	expectQuery_MusicGroup_GetMembers(mock, musicgroupID, peopleIDs)
// }

// func expectQuery_MusicAlbum_GetTracks(mock sqlmock.Sqlmock, musicalbumID IDType, musicrecordingIDs []IDType) {
// 	musicrecordingRows := make([]map[string]driver.Value, len(musicrecordingIDs))
// 	for i, id := range musicrecordingIDs {
// 		musicrecordingRows[i] = musicrecordings[id]
// 	}

// 	// Main query
// 	mock.ExpectQuery(`SELECT (.+) FROM musicrecording`).
// 		WithArgs(musicalbumID).
// 		WillReturnRows(
// 			mockResultRows(
// 				musicrecordingRows...,
// 			),
// 		)

// 	for _, musicrecordingID := range musicrecordingIDs {
// 		// .GetAudio query
// 		var audioobjectID = musicrecordings[musicrecordingID]["audio_id"].(IDType)
// 		expectQuery_MusicRecording_GetAudio(mock, audioobjectID)

// 		// .GetComposition query
// 		var musiccompositionID = musicrecordings[musicrecordingID]["recording_of_id"].(IDType)
// 		expectQuery_MusicRecording_GetComposition(mock, musiccompositionID)
// 	}
// }

// func expectQuery_MusicAlbum_GetArtist(mock sqlmock.Sqlmock, musicgroupID IDType) {
// 	var imageobjectID IDType = musicgroups[musicgroupID]["image_id"].(IDType)

// 	mock.ExpectQuery(`SELECT (.+) FROM musicgroup`).
// 		WithArgs(musicgroupID).
// 		WillReturnRows(
// 			mockResultRows(musicgroups[musicgroupID]),
// 		)

// 	mock.ExpectQuery(`SELECT (.+) FROM imageobject`).
// 		WithArgs(imageobjectID).
// 		WillReturnRows(
// 			mockResultRows(imageobjects[imageobjectID]),
// 		)
// }
