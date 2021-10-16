package models_test

// NOTE: the delete sql method was remvoved from the updatemembers method, so that will need to be reflected in these tests

// import (
// 	"database/sql"
// 	"database/sql/driver"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/lib/pq"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/onsi/gomega/gstruct"
// 	gomegatypes "github.com/onsi/gomega/types"
// 	"gopkg.in/DATA-DOG/go-sqlmock.v1"

// 	. "github.com/consensys/ugo/pkg/models"
// )

// var _ = Describe("MusicGroup", func() {
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
// 		Context("with a non-empty ethereumAddress argument", func() {
// 			var ethereumAddress = "0xdeadbeef"
// 			var musicgroupIDs = []IDType{998, 999}

// 			var musicgroups []*MusicGroup
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 2, Offset: 3}
// 				expectQuery_MusicGroup_All(mock, ethereumAddress, musicgroupIDs, selectParams)
// 				musicgroups, err = (&MusicGroup{}).All(dbx, ethereumAddress, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected musicgroups", func() {
// 				Expect(musicgroups).To(HaveLen(len(musicgroupIDs)))

// 				musicgroupMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, musicgroupID := range musicgroupIDs {
// 					musicgroupMatchers = append(musicgroupMatchers, PointTo(matchMusicgroupFixture(musicgroupID, nil)))
// 				}
// 				Expect(musicgroups).To(ConsistOf(musicgroupMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with an empty ethereumAddress argument", func() {
// 			var ethereumAddress = ""
// 			var musicgroupIDs = []IDType{998, 999}

// 			var musicgroups []*MusicGroup
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 2, Offset: 3}
// 				expectQuery_MusicGroup_All(mock, ethereumAddress, musicgroupIDs, selectParams)
// 				musicgroups, err = (&MusicGroup{}).All(dbx, ethereumAddress, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected musicgroups", func() {
// 				Expect(musicgroups).To(HaveLen(len(musicgroupIDs)))

// 				musicgroupMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, musicgroupID := range musicgroupIDs {
// 					musicgroupMatchers = append(musicgroupMatchers, PointTo(matchMusicgroupFixture(musicgroupID, nil)))
// 				}
// 				Expect(musicgroups).To(ConsistOf(musicgroupMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})
// 	})

// 	Context("when .Get is called", func() {
// 		var musicgroupID IDType = 999
// 		var peopleIDs = []IDType{1001, 1002}

// 		var musicgroup *MusicGroup
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicGroup_Get(mock, musicgroupID, peopleIDs)

// 			musicgroup, err = (&MusicGroup{}).Get(dbx, musicgroupID)
// 			if err != nil {
// 				Fail(err.Error())
// 			}
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil musicgroup", func() {
// 			Expect(musicgroup).ToNot(BeNil())
// 			Expect(*musicgroup).To(matchMusicgroupFixture(musicgroupID, peopleIDs))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .GetByCID is called", func() {
// 		var cids = []string{"cid1", "cid2"}
// 		var musicgroupID IDType = 999
// 		var peopleIDs = []IDType{1001, 1002}

// 		var musicgroup *MusicGroup
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicGroup_GetByCID(mock, cids, musicgroupID, peopleIDs)
// 			musicgroup, err = (&MusicGroup{}).GetByCID(dbx, cids...)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil musicgroup", func() {
// 			Expect(musicgroup).ToNot(BeNil())
// 			Expect(*musicgroup).To(matchMusicgroupFixture(musicgroupID, peopleIDs))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Create is called", func() {
// 		var mg *MusicGroup
// 		var newImageobjectID IDType = 4446
// 		var newMusicgroupID IDType = 5557
// 		var err error

// 		BeforeEach(func() {
// 			mg = &MusicGroup{
// 				CID:         "_cid",
// 				CIDs:        pq.StringArray([]string{"cid1", "cid2"}),
// 				Type:        "_type",
// 				Context:     "_context",
// 				Name:        strPtr("_name"),
// 				Description: strPtr("_description"),
// 				Email:       strPtr("_email"),
// 				Image:       &ImageObject{CID: "_cid", Type: "_type", Context: "_context", ContentURL: strPtr("_contentURL"), EncodingFormat: strPtr("_encodingFormat"), ID: 123},
// 				Members: []*Person{
// 					{ID: 888, Description: strPtr("_description"), PercentageShares: f64Ptr(0.321), MusicGroupAdmin: boolPtr(true)},
// 					{ID: 999, Description: strPtr("_description"), PercentageShares: f64Ptr(0.456), MusicGroupAdmin: boolPtr(false)},
// 				},
// 			}

// 			expectQuery_MusicGroup_Create(mock, mg, newMusicgroupID, newImageobjectID)
// 			_, err = mg.Create(dbx)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should populate the MusicGroup.ID field", func() {
// 			Expect(mg.ID).To(Equal(newMusicgroupID))
// 		})

// 		It("should populate the MusicGroup.Image.ID field", func() {
// 			Expect(mg.Image.ID).To(Equal(newImageobjectID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Update is called", func() {
// 		var mg *MusicGroup
// 		var musicgroupID IDType = 5557
// 		var newMembers = []IDType{0}
// 		var oldMembers = []IDType{1002}
// 		var err error

// 		BeforeEach(func() {
// 			mg = &MusicGroup{
// 				CID:         "_cid",
// 				CIDs:        pq.StringArray([]string{"cid1", "cid2"}),
// 				Type:        "_type",
// 				Context:     "_context",
// 				Name:        strPtr("_name"),
// 				Description: strPtr("_description"),
// 				Email:       strPtr("_email"),
// 				Image:       &ImageObject{CID: "_cid", Type: "_type", Context: "_context", ContentURL: strPtr("_contentURL"), EncodingFormat: strPtr("_encodingFormat"), ID: 123},
// 				Members: []*Person{
// 					{ID: 0, CID: "_cid", GivenName: strPtr("_familyname"), Type: "type", Context: "context", EthereumAddress: "_familyname", FamilyName: strPtr("_familyname"), Description: strPtr("_description"), PercentageShares: f64Ptr(0.321), MusicGroupAdmin: boolPtr(true)},
// 				},
// 			}

// 			expectQuery_MusicGroup_Update(mock, mg, musicgroupID, newMembers, oldMembers)
// 			_, err = mg.Update(dbx, musicgroupID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Delete is called", func() {
// 		var mg *MusicGroup
// 		var musicgroupID IDType = 5557
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_MusicGroup_Delete(mock, mg, musicgroupID)
// 			err = mg.Delete(dbx, musicgroupID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})
// })

// func expectQuery_MusicGroup_All(mock sqlmock.Sqlmock, ethereumAddress string, musicgroupIDs []IDType, selectParams *SelectQuery) {
// 	musicgroupRows := make([]map[string]driver.Value, len(musicgroupIDs))
// 	for i, id := range musicgroupIDs {
// 		musicgroupRows[i] = musicgroups[id]
// 	}

// 	var args []driver.Value
// 	if ethereumAddress != "" {
// 		args = append(args, ethereumAddress)
// 	}
// 	if selectParams.Limit > 0 {
// 		args = append(args, selectParams.Limit, selectParams.Offset)
// 	}

// 	e := mock.ExpectQuery(`SELECT (.+) FROM musicgroup`)
// 	if len(args) > 0 {
// 		e = e.WithArgs(args...)
// 	}
// 	e = e.WillReturnRows(
// 		mockResultRows(
// 			musicgroupRows...,
// 		),
// 	)

// 	for _, id := range musicgroupIDs {
// 		expectQuery_MusicGroup_GetImage(mock, musicgroups[id]["image_id"].(IDType))
// 	}
// }

// func expectQuery_MusicGroup_Get(mock sqlmock.Sqlmock, musicgroupID IDType, peopleIDs []IDType) {
// 	mock.ExpectQuery(`SELECT (.+) FROM musicgroup`).
// 		WithArgs(musicgroupID).
// 		WillReturnRows(
// 			mockResultRows(
// 				musicgroups[musicgroupID],
// 			),
// 		)

// 	expectQuery_MusicGroup_GetMembers(mock, musicgroupID, peopleIDs)
// 	expectQuery_MusicGroup_GetImage(mock, musicgroups[musicgroupID]["image_id"].(IDType))
// }

// func expectQuery_MusicGroup_GetByCID(mock sqlmock.Sqlmock, cids []string, musicgroupID IDType, peopleIDs []IDType) {
// 	mock.ExpectQuery(`SELECT (.+) FROM musicgroup`).
// 		WithArgs(pq.Array(cids)).
// 		WillReturnRows(
// 			mockResultRows(
// 				musicgroups[musicgroupID],
// 			),
// 		)

// 	expectQuery_MusicGroup_GetMembers(mock, musicgroupID, peopleIDs)
// 	expectQuery_MusicGroup_GetImage(mock, musicgroups[musicgroupID]["image_id"].(IDType))
// }

// func expectQuery_MusicGroup_GetMembers(mock sqlmock.Sqlmock, musicgroupID IDType, peopleIDs []IDType) {
// 	peopleRows := make([]map[string]driver.Value, len(peopleIDs))
// 	for i, id := range peopleIDs {
// 		peopleRows[i] = people[id]
// 	}

// 	mock.ExpectQuery(`SELECT (.+) FROM person`).
// 		WithArgs(musicgroupID).
// 		WillReturnRows(
// 			mockResultRows(
// 				peopleRows...,
// 			),
// 		)
// }

// func expectQuery_MusicGroup_GetImage(mock sqlmock.Sqlmock, imageobjectID IDType) {
// 	mock.ExpectQuery(`SELECT (.+) FROM imageobject`).
// 		WithArgs(imageobjectID).
// 		WillReturnRows(
// 			mockResultRows(
// 				imageobjects[imageobjectID],
// 			),
// 		)
// }

// func expectQuery_MusicGroup_Create(mock sqlmock.Sqlmock, mg *MusicGroup, newMusicgroupID, newImageobjectID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`INSERT INTO imageobject`).
// 		WithArgs(mg.Image.CID, mg.Image.Type, mg.Context, AnyTime{}, AnyTime{}, mg.Image.ContentURL, mg.Image.EncodingFormat).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newImageobjectID}))

// 	mock.ExpectQuery(`INSERT INTO musicgroup`).
// 		WithArgs(mg.CID, pq.Array([]string{mg.CID}), mg.Type, mg.Context, AnyTime{}, AnyTime{}, mg.Name, mg.Description, mg.Email).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newMusicgroupID}))

// 	for _, member := range mg.Members {
// 		mock.ExpectExec(`INSERT INTO musicgroup_members`).
// 			WithArgs(member.ID, member.Description, member.PercentageShares, member.MusicGroupAdmin).
// 			WillReturnResult(sqlmock.NewResult(123, 1))
// 	}

// 	mock.ExpectCommit()
// }

// func expectQuery_MusicGroup_Update(mock sqlmock.Sqlmock, mg *MusicGroup, musicgroupID IDType, newMembers []IDType, oldMembers []IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectExec(`UPDATE musicgroup`).
// 		WithArgs(mg.CID, mg.CID, AnyTime{}, mg.Type, mg.Context, mg.Name, mg.Description, mg.Email, musicgroupID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectExec(`UPDATE imageobject`).
// 		WithArgs(mg.Image.CID, mg.Type, mg.Context, AnyTime{}, mg.Image.ContentURL, mg.Image.EncodingFormat, mg.Image.ID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	expectQuery_MusicGroup_UpdateMembers(mock, mg, newMembers, musicgroupID, oldMembers)

// 	mock.ExpectCommit()
// }

// func expectQuery_MusicGroup_UpdateMembers(mock sqlmock.Sqlmock, mg *MusicGroup, newMembers []IDType, musicgroupID IDType, oldMembers []IDType) {
// 	expectQuery_MusicGroup_GetMembers(mock, musicgroupID, oldMembers)

// 	for _, p := range oldMembers {
// 		mock.ExpectExec(`DELETE FROM person`).
// 			WithArgs(p).
// 			WillReturnResult(sqlmock.NewResult(123, 1))

// 		mock.ExpectExec(`DELETE FROM musicgroup_members`).
// 			WithArgs(p, musicgroupID).
// 			WillReturnResult(sqlmock.NewResult(123, 1))
// 	}

// 	var id IDType = 123
// 	for _, p := range mg.Members {
// 		if p.ID == 0 {
// 			mock.ExpectQuery(`INSERT INTO person`).
// 				WithArgs(p.CID, p.Type, p.Context, AnyTime{}, AnyTime{},
// 					p.EthereumAddress, p.GivenName, p.FamilyName, p.Email).
// 				WillReturnRows(mockResultRows(map[string]driver.Value{"id": id}))

// 			mock.ExpectExec(`INSERT INTO musicgroup_members`).
// 				WithArgs(musicgroupID, id, p.Description, p.PercentageShares, p.MusicGroupAdmin).
// 				WillReturnResult(sqlmock.NewResult(123, 1))
// 		} else {
// 			mock.ExpectExec(`UPDATE person`).
// 				WithArgs(p.CID, p.Type, p.Context, AnyTime{},
// 					p.EthereumAddress, p.GivenName, p.FamilyName, p.Email, p.ID).
// 				WillReturnResult(sqlmock.NewResult(123, 1))

// 			mock.ExpectQuery(`SELECT (.+) FROM person`).
// 				WithArgs(musicgroupID, p.ID).
// 				WillReturnRows(mockResultRows())

// 			mock.ExpectExec(`INSERT INTO musicgroup_members`).
// 				WithArgs(musicgroupID, id, p.Description, p.PercentageShares, p.MusicGroupAdmin).
// 				WillReturnResult(sqlmock.NewResult(123, 1))
// 		}
// 	}
// }

// func expectQuery_MusicGroup_Delete(mock sqlmock.Sqlmock, mg *MusicGroup, musicgroupID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectExec(`DELETE FROM person`).
// 		WithArgs(musicgroupID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectExec(`DELETE FROM musicgroup`).
// 		WithArgs(musicgroupID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectCommit()
// }
