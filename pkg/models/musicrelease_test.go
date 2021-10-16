package models_test

// import (
// 	"database/sql"
// 	"database/sql/driver"
// 	"time"

// 	. "github.com/consensys/ugo/pkg/models"
// 	"github.com/jmoiron/sqlx"
// 	"github.com/lib/pq"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/onsi/gomega/gstruct"
// 	gomegatypes "github.com/onsi/gomega/types"
// 	"gopkg.in/DATA-DOG/go-sqlmock.v1"
// )

// var _ = Describe("MusicRelease", func() {
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
// 		Context("with a non-nil byArtist argument", func() {
// 			var byArtist = idPtr(222)
// 			var musicreleaseIDs = []IDType{818, 717}
// 			var musicrecordingIDs = []IDType{1, 2}
// 			var peopleIDs = []IDType{1001, 1002}

// 			var musicreleases []*MusicRelease
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 32, Offset: 1}
// 				expectQuery_MusicRelease_All(mock, byArtist, musicreleaseIDs, musicrecordingIDs, peopleIDs, selectParams)
// 				musicreleases, err = (&MusicRelease{}).All(dbx, byArtist, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected musicreleases", func() {
// 				Expect(musicreleases).To(HaveLen(len(musicreleaseIDs)))

// 				musicreleaseMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, musicreleaseID := range musicreleaseIDs {
// 					musicreleaseMatchers = append(musicreleaseMatchers, PointTo(matchMusicreleaseFixture(musicreleaseID)))
// 				}
// 				Expect(musicreleases).To(ConsistOf(musicreleaseMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with a nil byArtist argument", func() {
// 			var byArtist *IDType = nil
// 			var musicreleaseIDs = []IDType{818, 717}
// 			var musicrecordingIDs = []IDType{1, 2}
// 			var peopleIDs = []IDType{1001, 1002}

// 			var musicreleases []*MusicRelease
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 32, Offset: 1}
// 				expectQuery_MusicRelease_All(mock, byArtist, musicreleaseIDs, musicrecordingIDs, peopleIDs, selectParams)
// 				musicreleases, err = (&MusicRelease{}).All(dbx, byArtist, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected musicreleases", func() {
// 				Expect(musicreleases).To(HaveLen(len(musicreleaseIDs)))

// 				musicreleaseMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, musicreleaseID := range musicreleaseIDs {
// 					musicreleaseMatchers = append(musicreleaseMatchers, PointTo(matchMusicreleaseFixture(musicreleaseID)))
// 				}
// 				Expect(musicreleases).To(ConsistOf(musicreleaseMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})
// 	})

// 	Context("when .Create is called", func() {
// 		var newMusicreleaseID IDType = 126
// 		var mr *MusicRelease
// 		var err error

// 		BeforeEach(func() {
// 			mr = &MusicRelease{
// 				CID:                "_cid",
// 				Type:               "_type",
// 				Context:            "_context",
// 				CreatedAt:          time.Now(),
// 				UpdatedAt:          time.Now(),
// 				Description:        strPtr("_description"),
// 				DatePublished:      strPtr("_date_published"),
// 				CatalogNumber:      strPtr("_catalog_number"),
// 				MusicReleaseFormat: strPtr("_music_release_format"),
// 				Price:              f64Ptr(321.123),
// 				// RecordLabel:
// 				ReleaseOf: &MusicAlbum{
// 					ID:        111,
// 					CID:       "_cid",
// 					Type:      "_type",
// 					Context:   "_context",
// 					CreatedAt: time.Now(),
// 					UpdatedAt: time.Now(),
// 					Name:      strPtr("_name"),
// 					Tracks: []*MusicRecording{
// 						{
// 							ID:            991,
// 							CID:           "_cid",
// 							Type:          "_type",
// 							Context:       "_context",
// 							CreatedAt:     time.Now(),
// 							UpdatedAt:     time.Now(),
// 							Name:          strPtr("_name"),
// 							Duration:      strPtr("_duration"),
// 							Isrc:          strPtr("_isrc"),
// 							Position:      strPtr("_position"),
// 							Genres:        pq.StringArray([]string{"genre1", "genre2"}),
// 							AudioID:       idPtr(225),
// 							ByArtistID:    idPtr(226),
// 							RecordingOfID: idPtr(227),
// 							ImageID:       idPtr(228),
// 							Visibility:    strPtr("_name"),
// 							ByArtist:      &MusicGroup{ID: 4949},
// 							RecordingOf: &MusicComposition{
// 								ID:   9191,
// 								CID:  "_cid",
// 								Name: strPtr("_name"),
// 								Iswc: strPtr("_iswc"),
// 							},
// 							Audio: &AudioObject{
// 								ID:             551,
// 								CID:            "_cid",
// 								Type:           "_type",
// 								Context:        "_context",
// 								CreatedAt:      time.Now(),
// 								UpdatedAt:      time.Now(),
// 								ContentURL:     strPtr("_contentURL"),
// 								EncodingFormat: strPtr("_encodingFormat"),
// 							},
// 							Image: &ImageObject{
// 								ID:             551,
// 								CID:            "_cid",
// 								Type:           "_type",
// 								Context:        "_context",
// 								CreatedAt:      time.Now(),
// 								UpdatedAt:      time.Now(),
// 								ContentURL:     strPtr("_contentURL"),
// 								EncodingFormat: strPtr("_encodingFormat"),
// 							},
// 						},
// 					},
// 					AlbumProductionType: strPtr("_album_production_type"),
// 					AlbumReleaseType:    strPtr("_album_release_type"),
// 					ByArtist:            &MusicGroup{ID: 13213},
// 				},
// 				Image: &ImageObject{
// 					ID:             1582,
// 					CID:            "_cid",
// 					Type:           "_type",
// 					Context:        "_context",
// 					CreatedAt:      time.Now(),
// 					UpdatedAt:      time.Now(),
// 					ContentURL:     strPtr("_contentURL"),
// 					EncodingFormat: strPtr("_encodingFormat"),
// 				},
// 				RecordingLabelID: idPtr(19822),
// 				ReleaseOfID:      idPtr(12935),
// 				ImageID:          idPtr(58271),
// 				Visibility:       strPtr("_visibility"),
// 				Active:           true,
// 			}

// 			expectQuery_MusicRelease_Create(mock, mr, newMusicreleaseID)
// 			_, err = mr.Create(dbx)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should populate the MusicRelease.ID field", func() {
// 			Expect(mr.ID).To(Equal(newMusicreleaseID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Update is called", func() {
// 		var musicreleaseID IDType = 818
// 		var mr *MusicRelease
// 		var err error

// 		BeforeEach(func() {
// 			mr = &MusicRelease{
// 				CID:                "_cid",
// 				Type:               "_type",
// 				Context:            "_context",
// 				CreatedAt:          time.Now(),
// 				UpdatedAt:          time.Now(),
// 				Description:        strPtr("_description"),
// 				DatePublished:      strPtr("_date_published"),
// 				CatalogNumber:      strPtr("_catalog_number"),
// 				MusicReleaseFormat: strPtr("_music_release_format"),
// 				Price:              f64Ptr(321.123),
// 				// RecordLabel:
// 				ReleaseOf: &MusicAlbum{
// 					ID:        111,
// 					CID:       "_cid",
// 					Type:      "_type",
// 					Context:   "_context",
// 					CreatedAt: time.Now(),
// 					UpdatedAt: time.Now(),
// 					Name:      strPtr("_name"),
// 					Tracks: []*MusicRecording{
// 						{
// 							ID:            991,
// 							CID:           "_cid",
// 							Type:          "_type",
// 							Context:       "_context",
// 							CreatedAt:     time.Now(),
// 							UpdatedAt:     time.Now(),
// 							Name:          strPtr("_name"),
// 							Duration:      strPtr("_duration"),
// 							Isrc:          strPtr("_isrc"),
// 							Position:      strPtr("_position"),
// 							Genres:        pq.StringArray([]string{"genre1", "genre2"}),
// 							AudioID:       idPtr(225),
// 							ByArtistID:    idPtr(226),
// 							RecordingOfID: idPtr(227),
// 							ImageID:       idPtr(228),
// 							Visibility:    strPtr("_name"),
// 							ByArtist:      &MusicGroup{ID: 4949},
// 							RecordingOf: &MusicComposition{
// 								ID:   9191,
// 								CID:  "_cid",
// 								Name: strPtr("_name"),
// 								Iswc: strPtr("_iswc"),
// 							},
// 							Audio: &AudioObject{
// 								ID:             551,
// 								CID:            "_cid",
// 								Type:           "_type",
// 								Context:        "_context",
// 								CreatedAt:      time.Now(),
// 								UpdatedAt:      time.Now(),
// 								ContentURL:     strPtr("_contentURL"),
// 								EncodingFormat: strPtr("_encodingFormat"),
// 							},
// 							Image: &ImageObject{
// 								ID:             551,
// 								CID:            "_cid",
// 								Type:           "_type",
// 								Context:        "_context",
// 								CreatedAt:      time.Now(),
// 								UpdatedAt:      time.Now(),
// 								ContentURL:     strPtr("_contentURL"),
// 								EncodingFormat: strPtr("_encodingFormat"),
// 							},
// 						},
// 					},
// 					AlbumProductionType: strPtr("_album_production_type"),
// 					AlbumReleaseType:    strPtr("_album_release_type"),
// 					ByArtist:            &MusicGroup{ID: 13213},
// 				},
// 				Image: &ImageObject{
// 					ID:             1582,
// 					CID:            "_cid",
// 					Type:           "_type",
// 					Context:        "_context",
// 					CreatedAt:      time.Now(),
// 					UpdatedAt:      time.Now(),
// 					ContentURL:     strPtr("_contentURL"),
// 					EncodingFormat: strPtr("_encodingFormat"),
// 				},
// 				RecordingLabelID: idPtr(19822),
// 				ReleaseOfID:      idPtr(12935),
// 				ImageID:          idPtr(58271),
// 				Visibility:       strPtr("_visibility"),
// 				Active:           true,
// 			}

// 			expectQuery_MusicRelease_Update(mock, musicreleaseID, mr)
// 			_, err = mr.Update(dbx, musicreleaseID)
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
// 		var musicreleaseID IDType = 1
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicRelease_Delete(mock, musicreleaseID)
// 			err = (&MusicRelease{}).Delete(dbx, musicreleaseID)
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

// 	Context("when .GetImage is called", func() {
// 		var imageobjectID IDType = 986
// 		var imageobject *ImageObject
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicRelease_GetImage(mock, imageobjectID)
// 			imageobject, err = (&MusicRelease{}).GetImage(dbx, imageobjectID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil imageobject", func() {
// 			Expect(imageobject).ToNot(BeNil())
// 			Expect(*imageobject).To(matchImageobjectFixture(imageobjectID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})
// })

// func expectQuery_MusicRelease_All(mock sqlmock.Sqlmock, byArtist *IDType, musicreleaseIDs, musicrecordingIDs, peopleIDs []IDType, selectParams *SelectQuery) {
// 	musicreleaseRows := make([]map[string]driver.Value, len(musicreleaseIDs))
// 	for i, id := range musicreleaseIDs {
// 		musicreleaseRows[i] = musicreleases[id]
// 	}

// 	var args []driver.Value
// 	if byArtist != nil {
// 		args = append(args, byArtist)
// 	}
// 	if selectParams.Limit > 0 {
// 		args = append(args, selectParams.Limit, selectParams.Offset)
// 	}

// 	e := mock.ExpectQuery(`SELECT (.+) FROM musicrelease`)
// 	if len(args) > 0 {
// 		e = e.WithArgs(args...)
// 	}
// 	e = e.WillReturnRows(
// 		mockResultRows(
// 			musicreleaseRows...,
// 		),
// 	)

// 	for _, musicreleaseID := range musicreleaseIDs {
// 		musicrelease := musicreleases[musicreleaseID]
// 		expectQuery_MusicAlbum_Get(mock, musicrelease["release_of_id"].(IDType), musicrecordingIDs, peopleIDs)
// 		expectQuery_MusicRelease_GetImage(mock, musicrelease["image_id"].(IDType))
// 	}
// }

// func expectQuery_MusicRelease_GetByID(mock sqlmock.Sqlmock, musicreleaseID IDType, musicrecordingIDs []IDType, peopleIDs []IDType) {
// 	var musicrelease = musicreleases[musicreleaseID]

// 	mock.ExpectQuery(`SELECT (.+) FROM musicrelease`).
// 		WithArgs(musicreleaseID).
// 		WillReturnRows(
// 			mockResultRows(
// 				musicrelease,
// 			),
// 		)

// 	expectQuery_MusicAlbum_Get(mock, musicrelease["release_of_id"].(IDType), musicrecordingIDs, peopleIDs)
// 	expectQuery_MusicRelease_GetImage(mock, musicrelease["image_id"].(IDType))
// }

// func expectQuery_MusicRelease_GetImage(mock sqlmock.Sqlmock, imageobjectID IDType) {
// 	mock.ExpectQuery(`SELECT (.+) FROM imageobject`).
// 		WithArgs(imageobjectID).
// 		WillReturnRows(
// 			mockResultRows(
// 				imageobjects[imageobjectID],
// 			),
// 		)
// }

// func expectQuery_MusicRelease_Create(mock sqlmock.Sqlmock, mr *MusicRelease, newMusicreleaseID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectExec(`INSERT INTO musicalbum`).
// 		WithArgs(mr.ReleaseOf.CID, mr.ReleaseOf.Type, mr.Context, AnyTime{}, AnyTime{}, mr.ReleaseOf.Name, mr.ReleaseOf.AlbumProductionType, mr.ReleaseOf.AlbumReleaseType, mr.ReleaseOf.ByArtist.ID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	for _, t := range mr.ReleaseOf.Tracks {
// 		mock.ExpectExec(`INSERT INTO musiccomposition`).
// 			WithArgs(t.RecordingOf.CID, t.RecordingOf.Type, t.Context, AnyTime{}, AnyTime{}, t.RecordingOf.Name, t.RecordingOf.Iswc).
// 			WillReturnResult(sqlmock.NewResult(123, 1))

// 		mock.ExpectExec(`INSERT INTO audioobject`).
// 			WithArgs(t.Audio.CID, t.Audio.Type, t.Context, AnyTime{}, AnyTime{}, t.Audio.ContentURL, t.Audio.EncodingFormat).
// 			WillReturnResult(sqlmock.NewResult(123, 1))

// 		mock.ExpectExec(`INSERT INTO musicrecording`).
// 			WithArgs(t.CID, t.Type, t.Context, AnyTime{}, AnyTime{}, t.Name, t.Duration, t.Isrc, t.Position, pq.Array(t.Genres), mr.ReleaseOf.ByArtist.ID).
// 			WillReturnResult(sqlmock.NewResult(123, 1))

// 		mock.ExpectExec(`INSERT INTO musicalbum_tracks`).
// 			WillReturnResult(sqlmock.NewResult(123, 1))
// 	}

// 	mock.ExpectExec(`INSERT INTO imageobject`).
// 		WithArgs(mr.Image.CID, mr.Image.Type, mr.Image.Context, AnyTime{}, AnyTime{}, mr.Image.ContentURL, mr.Image.EncodingFormat).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectQuery(`INSERT INTO musicrelease`).
// 		WithArgs(mr.CID, mr.Type, mr.Context, AnyTime{}, AnyTime{}, true, mr.Description, mr.DatePublished, mr.CatalogNumber, mr.MusicReleaseFormat, mr.Price, nil).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newMusicreleaseID}))

// 	mock.ExpectCommit()
// }

// func expectQuery_MusicRelease_Update(mock sqlmock.Sqlmock, musicreleaseID IDType, mr *MusicRelease) {
// 	mock.ExpectBegin()

// 	mock.ExpectExec(`UPDATE musicalbum`).
// 		WithArgs(mr.ReleaseOf.CID, mr.ReleaseOf.Type, mr.Context, AnyTime{}, mr.ReleaseOf.Name, mr.ReleaseOf.AlbumProductionType, mr.ReleaseOf.AlbumReleaseType, mr.ReleaseOf.ByArtist.ID, mr.ReleaseOf.ID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	for _, t := range mr.ReleaseOf.Tracks {
// 		mock.ExpectExec(`UPDATE musiccomposition`).
// 			WithArgs(t.CID, t.RecordingOf.Type, t.Context, AnyTime{}, t.RecordingOf.Name, t.RecordingOf.Iswc, t.RecordingOf.ID).
// 			WillReturnResult(sqlmock.NewResult(123, 1))

// 		mock.ExpectExec(`UPDATE audioobject`).
// 			WithArgs(t.CID, t.Audio.Type, t.Audio.Context, AnyTime{}, t.Audio.ContentURL, t.Audio.EncodingFormat, t.Audio.ID).
// 			WillReturnResult(sqlmock.NewResult(123, 1))

// 		mock.ExpectExec(`UPDATE musicrecording`).
// 			WithArgs(t.CID, t.Type, t.Context, AnyTime{}, t.Name, t.Duration, t.Isrc, t.Position, pq.Array(t.Genres), mr.ReleaseOf.ByArtist.ID, t.Audio.ID, t.RecordingOf.ID, t.ID).
// 			WillReturnResult(sqlmock.NewResult(123, 1))
// 	}

// 	mock.ExpectExec(`UPDATE musicrelease`).
// 		WithArgs(mr.CID, mr.Type, mr.Context, AnyTime{}, mr.Active, mr.Description, mr.DatePublished, mr.CatalogNumber, mr.MusicReleaseFormat, mr.Price, nil, mr.ReleaseOf.ID, musicreleaseID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectExec(`UPDATE imageobject`).
// 		WithArgs(mr.Image.CID, mr.Image.Type, mr.Image.Context, AnyTime{}, mr.Image.ContentURL, mr.Image.EncodingFormat, mr.Image.ID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectCommit()
// }

// func expectQuery_MusicRelease_Delete(mock sqlmock.Sqlmock, musicreleaseID IDType) {
// 	mock.ExpectBegin()
// 	mock.ExpectExec(`DELETE FROM audioobject`).
// 		WithArgs(musicreleaseID).
// 		WillReturnResult(sqlmock.NewResult(-1, 1))
// 	mock.ExpectExec(`DELETE FROM musiccomposition`).
// 		WithArgs(musicreleaseID).
// 		WillReturnResult(sqlmock.NewResult(-1, 1))
// 	mock.ExpectExec(`DELETE FROM musicrecording`).
// 		WithArgs(musicreleaseID).
// 		WillReturnResult(sqlmock.NewResult(-1, 1))
// 	mock.ExpectExec(`DELETE from musicalbum`).
// 		WithArgs(musicreleaseID).
// 		WillReturnResult(sqlmock.NewResult(-1, 1))
// 	mock.ExpectCommit()
// }
