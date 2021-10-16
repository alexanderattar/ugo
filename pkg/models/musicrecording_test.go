package models_test

// import (
// 	"database/sql"
// 	"database/sql/driver"

// 	. "github.com/consensys/ugo/pkg/models"
// 	"github.com/jmoiron/sqlx"
// 	"github.com/lib/pq"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/onsi/gomega/gstruct"
// 	gomegatypes "github.com/onsi/gomega/types"
// 	"gopkg.in/DATA-DOG/go-sqlmock.v1"
// )

// var _ = Describe("MusicRecording", func() {
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
// 		Context("with a non-nil byArtistID argument", func() {
// 			var byArtistID = idPtr(222)
// 			var musicrecordingIDs = []IDType{1, 2}
// 			var peopleIDs = []IDType{1001, 1002}

// 			var musicrecordings []*MusicRecording
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 32, Offset: 1}
// 				expectQuery_MusicRecording_All(mock, byArtistID, musicrecordingIDs, peopleIDs, selectParams)
// 				musicrecordings, err = (&MusicRecording{}).All(dbx, byArtistID, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected musicrecordings", func() {
// 				Expect(musicrecordings).To(HaveLen(len(musicrecordingIDs)))

// 				musicrecordingMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, musicrecordingID := range musicrecordingIDs {
// 					musicrecordingMatchers = append(musicrecordingMatchers, PointTo(matchMusicrecordingFixture(musicrecordingID)))
// 				}
// 				Expect(musicrecordings).To(ConsistOf(musicrecordingMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with a nil byArtistID argument", func() {
// 			var byArtistID *IDType = nil
// 			var musicrecordingIDs = []IDType{1, 2}
// 			var peopleIDs = []IDType{1001, 1002}

// 			var musicrecordings []*MusicRecording
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 32, Offset: 1}
// 				expectQuery_MusicRecording_All(mock, byArtistID, musicrecordingIDs, peopleIDs, selectParams)
// 				musicrecordings, err = (&MusicRecording{}).All(dbx, byArtistID, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected musicrecordings", func() {
// 				Expect(musicrecordings).To(HaveLen(len(musicrecordingIDs)))

// 				musicrecordingMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, musicrecordingID := range musicrecordingIDs {
// 					musicrecordingMatchers = append(musicrecordingMatchers, PointTo(matchMusicrecordingFixture(musicrecordingID)))
// 				}
// 				Expect(musicrecordings).To(ConsistOf(musicrecordingMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})
// 	})

// 	Context("when .Create is called", func() {
// 		var newMusiccompositionID IDType = 123
// 		var newAudioobjectID IDType = 124
// 		var newImageobjectID IDType = 125
// 		var newMusicrecordingID IDType = 126
// 		var mr *MusicRecording
// 		var err error

// 		BeforeEach(func() {
// 			mr = &MusicRecording{
// 				CID:         "_cid",
// 				Type:        "_type",
// 				Context:     "_context",
// 				Name:        strPtr("_name"),
// 				Duration:    strPtr("_duration"),
// 				Isrc:        strPtr("_isrc"),
// 				Position:    strPtr("_position"),
// 				Genres:      pq.StringArray([]string{"one", "two"}),
// 				Image:       &ImageObject{CID: "_cid", Type: "_type", Context: "_context", ContentURL: strPtr("_contentURL"), EncodingFormat: strPtr("_encodingFormat"), ID: 123},
// 				Audio:       &AudioObject{CID: "_cid", Type: "_type", Context: "_context", ContentURL: strPtr("_contentURL"), EncodingFormat: strPtr("_encodingFormat"), ID: 123},
// 				RecordingOf: &MusicComposition{CID: "_cid", Type: "_type", Context: "_context", Name: strPtr("_name"), ID: 123},
// 				ByArtist:    &MusicGroup{ID: 123},
// 			}

// 			expectQuery_MusicRecording_Create(mock, mr, newMusiccompositionID, newAudioobjectID, newImageobjectID, newMusicrecordingID)
// 			_, err = mr.Create(dbx)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should populate the MusicRecording.RecordingOf.ID field", func() {
// 			Expect(mr.RecordingOf.ID).To(Equal(newMusiccompositionID))
// 		})

// 		It("should populate the MusicRecording.Audio.ID field", func() {
// 			Expect(mr.Audio.ID).To(Equal(newAudioobjectID))
// 		})

// 		It("should populate the MusicRecording.Image.ID field", func() {
// 			Expect(mr.Image.ID).To(Equal(newImageobjectID))
// 		})

// 		It("should populate the MusicRecording.ID field", func() {
// 			Expect(mr.ID).To(Equal(newMusicrecordingID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Update is called", func() {
// 		var musicrecordingID IDType = 1
// 		var mr *MusicRecording
// 		var err error

// 		BeforeEach(func() {
// 			mr = &MusicRecording{
// 				CID:         "_cid",
// 				Type:        "_type",
// 				Context:     "_context",
// 				Name:        strPtr("_name"),
// 				Duration:    strPtr("_duration"),
// 				Isrc:        strPtr("_isrc"),
// 				Position:    strPtr("_position"),
// 				Genres:      pq.StringArray([]string{"one", "two"}),
// 				Image:       &ImageObject{CID: "_cid", Type: "_type", Context: "_context", ContentURL: strPtr("_contentURL"), EncodingFormat: strPtr("_encodingFormat"), ID: 123},
// 				Audio:       &AudioObject{CID: "_cid", Type: "_type", Context: "_context", ContentURL: strPtr("_contentURL"), EncodingFormat: strPtr("_encodingFormat"), ID: 123},
// 				RecordingOf: &MusicComposition{CID: "_cid", Type: "_type", Context: "_context", Name: strPtr("_name"), ID: 123},
// 				ByArtist:    &MusicGroup{ID: 123},
// 			}

// 			expectQuery_MusicRecording_Update(mock, musicrecordingID, mr)
// 			_, err = mr.Update(dbx, musicrecordingID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Delete is called", func() {
// 		var musicrecordingID IDType = 1
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicRecording_Delete(mock, musicrecordingID)
// 			err = (&MusicRecording{}).Delete(dbx, musicrecordingID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .GetAudio is called", func() {
// 		var audioobjectID IDType = 456
// 		var audioobject *AudioObject
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicRecording_GetAudio(mock, audioobjectID)
// 			audioobject, err = (&MusicRecording{}).GetAudio(dbx, audioobjectID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil audioobject", func() {
// 			Expect(audioobject).ToNot(BeNil())
// 			Expect(*audioobject).To(matchAudioobjectFixture(audioobjectID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .GetComposition is called", func() {
// 		var musiccompositionID IDType = 333
// 		var musiccomposition *MusicComposition
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicRecording_GetComposition(mock, musiccompositionID)
// 			musiccomposition, err = (&MusicRecording{}).GetComposition(dbx, musiccompositionID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil musiccomposition", func() {
// 			Expect(musiccomposition).ToNot(BeNil())
// 			Expect(*musiccomposition).To(matchMusiccompositionFixture(musiccompositionID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})
// })

// func expectQuery_MusicRecording_All(mock sqlmock.Sqlmock, byArtistID *IDType, musicrecordingIDs, peopleIDs []IDType, selectParams *SelectQuery) {
// 	musicrecordingRows := make([]map[string]driver.Value, len(musicrecordingIDs))
// 	for i, id := range musicrecordingIDs {
// 		musicrecordingRows[i] = musicrecordings[id]
// 	}

// 	var args []driver.Value
// 	if byArtistID != nil {
// 		args = append(args, *byArtistID)
// 	}
// 	if selectParams.Limit > 0 {
// 		args = append(args, selectParams.Limit, selectParams.Offset)
// 	}

// 	e := mock.ExpectQuery(`SELECT (.+) FROM musicrecording`)
// 	if len(args) > 0 {
// 		e = e.WithArgs(args...)
// 	}
// 	e = e.WillReturnRows(
// 		mockResultRows(
// 			musicrecordingRows...,
// 		),
// 	)

// 	for _, musicrecording := range musicrecordingRows {
// 		expectQuery_MusicRecording_GetAudio(mock, musicrecording["audio_id"].(IDType))
// 		expectQuery_MusicRecording_GetComposition(mock, musicrecording["recording_of_id"].(IDType))
// 		expectQuery_MusicGroup_Get(mock, musicrecording["by_artist_id"].(IDType), peopleIDs)
// 	}
// }

// func expectQuery_MusicRecording_Create(mock sqlmock.Sqlmock, mr *MusicRecording, newMusiccompositionID, newAudioobjectID, newImageobjectID, newMusicrecordingID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`INSERT INTO musiccomposition`).
// 		WithArgs(mr.RecordingOf.CID, mr.RecordingOf.Type, mr.RecordingOf.Context, AnyTime{}, AnyTime{}, mr.RecordingOf.Name).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newMusiccompositionID}))

// 	mock.ExpectQuery(`INSERT INTO audioobject`).
// 		WithArgs(mr.Audio.CID, mr.Audio.Type, mr.Context, AnyTime{}, AnyTime{}, mr.Audio.ContentURL, mr.Audio.EncodingFormat).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newAudioobjectID}))

// 	mock.ExpectQuery(`INSERT INTO imageobject`).
// 		WithArgs(mr.Image.CID, mr.Image.Type, mr.Context, AnyTime{}, AnyTime{}, mr.Image.ContentURL, mr.Image.EncodingFormat).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newImageobjectID}))

// 	mock.ExpectQuery(`INSERT INTO musicrecording`).
// 		WithArgs(mr.CID, mr.Type, mr.Context, AnyTime{}, AnyTime{}, mr.Name, mr.Duration, mr.Isrc, mr.Position, pq.Array(mr.Genres), mr.ByArtist.ID, newImageobjectID).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newMusicrecordingID}))

// 	mock.ExpectCommit()
// }

// func expectQuery_MusicRecording_Update(mock sqlmock.Sqlmock, musicrecordingID IDType, mr *MusicRecording) {
// 	mock.ExpectBegin()

// 	mock.ExpectExec(`UPDATE imageobject`).
// 		WithArgs(mr.Image.CID, mr.Image.Type, mr.Image.Context, AnyTime{},
// 			mr.Image.ContentURL, mr.Image.EncodingFormat, mr.Image.ID).
// 		WillReturnResult(sqlmock.NewResult(mr.Image.ID, 1))

// 	mock.ExpectExec(`UPDATE audioobject`).
// 		WithArgs(mr.Audio.CID, mr.Audio.Type, mr.Context, AnyTime{},
// 			mr.Audio.ContentURL, mr.Audio.EncodingFormat, mr.Audio.ID).
// 		WillReturnResult(sqlmock.NewResult(mr.Audio.ID, 1))

// 	mock.ExpectExec(`UPDATE musiccomposition`).
// 		WithArgs(mr.RecordingOf.CID, mr.RecordingOf.Type, mr.Context, AnyTime{},
// 			mr.RecordingOf.Name, mr.RecordingOf.ID).
// 		WillReturnResult(sqlmock.NewResult(mr.RecordingOf.ID, 1))

// 	mock.ExpectExec(`UPDATE musicrecording`).
// 		WithArgs(mr.CID, AnyTime{}, mr.Type,
// 			mr.Context, mr.Name, mr.Duration, mr.Isrc, mr.Position,
// 			mr.Genres, mr.Audio.ID, mr.ByArtist.ID, mr.RecordingOf.ID,
// 			mr.Image.ID, musicrecordingID).
// 		WillReturnResult(sqlmock.NewResult(musicrecordingID, 1))

// 	mock.ExpectCommit()
// }

// func expectQuery_MusicRecording_Delete(mock sqlmock.Sqlmock, musicrecordingID IDType) {
// 	mock.ExpectBegin()
// 	mock.ExpectExec(`DELETE FROM audioobject`).
// 		WithArgs(musicrecordingID).
// 		WillReturnResult(sqlmock.NewResult(-1, 1))
// 	mock.ExpectExec(`DELETE FROM musiccomposition`).
// 		WithArgs(musicrecordingID).
// 		WillReturnResult(sqlmock.NewResult(-1, 1))
// 	mock.ExpectExec(`DELETE FROM musicrecording`).
// 		WithArgs(musicrecordingID).
// 		WillReturnResult(sqlmock.NewResult(-1, 1))
// 	mock.ExpectCommit()
// }

// func expectQuery_MusicRecording_GetAudio(mock sqlmock.Sqlmock, audioobjectID IDType) {
// 	mock.ExpectQuery(`SELECT (.+) FROM audioobject`).
// 		WithArgs(audioobjectID).
// 		WillReturnRows(
// 			mockResultRows(
// 				audioobjects[audioobjectID],
// 			),
// 		)
// }

// func expectQuery_MusicRecording_GetComposition(mock sqlmock.Sqlmock, musiccompositionID IDType) {
// 	mock.ExpectQuery(`SELECT (.+) FROM musiccomposition`).
// 		WithArgs(musiccompositionID).
// 		WillReturnRows(
// 			mockResultRows(
// 				musiccompositions[musiccompositionID],
// 			),
// 		)
// }
