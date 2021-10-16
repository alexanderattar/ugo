package models_test

// import (
// 	"database/sql"
// 	"database/sql/driver"

// 	"github.com/jmoiron/sqlx"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/onsi/gomega/gstruct"
// 	gomegatypes "github.com/onsi/gomega/types"
// 	"gopkg.in/DATA-DOG/go-sqlmock.v1"

// 	. "github.com/consensys/ugo/pkg/models"
// )

// var _ = Describe("Purchase", func() {
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

// 	Context("when .GetMusicRelease is called", func() {
// 		var musicreleaseID IDType = 818
// 		var musicrecordingIDs = []IDType{1, 2}
// 		var peopleIDs = []IDType{1001, 1002}

// 		var musicrelease *MusicRelease
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_Purchase_GetMusicRelease(mock, musicreleaseID, musicrecordingIDs, peopleIDs)
// 			musicrelease, err = (&Purchase{}).GetMusicRelease(dbx, musicreleaseID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil musicrelease", func() {
// 			Expect(musicrelease).ToNot(BeNil())
// 			Expect(*musicrelease).To(matchMusicreleaseFixture(musicreleaseID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .All is called", func() {
// 		Context("with a non-empty ethereumAddress argument", func() {
// 			Context("and a non-nil releaseID argument", func() {
// 				var ethereumAddress = "0xdeadbeef"
// 				var releaseID = idPtr(1234)
// 				var purchaseIDs = []IDType{263, 264}
// 				var musicrecordingIDs = [][]IDType{{1, 2}, {2, 1}}
// 				var peopleIDs = [][]IDType{{1001, 1002}, {1002, 1001}}

// 				var purchases []*Purchase
// 				var err error

// 				BeforeEach(func() {
// 					selectParams := &SelectQuery{Limit: 2, Offset: 3}
// 					expectQuery_Purchase_All(mock, releaseID, ethereumAddress, purchaseIDs, musicrecordingIDs, peopleIDs, selectParams)
// 					purchases, err = (&Purchase{}).All(dbx, releaseID, ethereumAddress, selectParams)
// 				})

// 				It("should return a nil error", func() {
// 					Expect(err).To(BeNil())
// 				})

// 				It("should return a non-nil slice containing the expected purchases", func() {
// 					Expect(purchases).To(HaveLen(len(purchaseIDs)))

// 					purchaseMatchers := []gomegatypes.GomegaMatcher{}
// 					for _, purchaseID := range purchaseIDs {
// 						purchaseMatchers = append(purchaseMatchers, PointTo(matchPurchaseFixture(purchaseID)))
// 					}
// 					Expect(purchases).To(ConsistOf(purchaseMatchers))
// 				})

// 				It("should execute the expected SQL queries", func() {
// 					if err = mock.ExpectationsWereMet(); err != nil {
// 						Fail("there were unfulfilled expectations: " + err.Error())
// 					}
// 				})
// 			})

// 			Context("and a nil releaseID argument", func() {
// 				var ethereumAddress = "0xdeadbeef"
// 				var releaseID *IDType = nil
// 				var purchaseIDs = []IDType{263, 264}
// 				var musicrecordingIDs = [][]IDType{{1, 2}, {2, 1}}
// 				var peopleIDs = [][]IDType{{1001, 1002}, {1002, 1001}}

// 				var purchases []*Purchase
// 				var err error

// 				BeforeEach(func() {
// 					selectParams := &SelectQuery{Limit: 2, Offset: 3}
// 					expectQuery_Purchase_All(mock, releaseID, ethereumAddress, purchaseIDs, musicrecordingIDs, peopleIDs, selectParams)
// 					purchases, err = (&Purchase{}).All(dbx, releaseID, ethereumAddress, selectParams)
// 				})

// 				It("should return a nil error", func() {
// 					Expect(err).To(BeNil())
// 				})

// 				It("should return a non-nil slice containing the expected purchases", func() {
// 					Expect(purchases).To(HaveLen(len(purchaseIDs)))

// 					purchaseMatchers := []gomegatypes.GomegaMatcher{}
// 					for _, purchaseID := range purchaseIDs {
// 						purchaseMatchers = append(purchaseMatchers, PointTo(matchPurchaseFixture(purchaseID)))
// 					}
// 					Expect(purchases).To(ConsistOf(purchaseMatchers))
// 				})

// 				It("should execute the expected SQL queries", func() {
// 					if err = mock.ExpectationsWereMet(); err != nil {
// 						Fail("there were unfulfilled expectations: " + err.Error())
// 					}
// 				})
// 			})
// 		})

// 		Context("with an empty ethereumAddress argument", func() {
// 			var ethereumAddress = ""
// 			var releaseID *IDType = nil
// 			var purchaseIDs = []IDType{263, 264}
// 			var musicrecordingIDs = [][]IDType{{1, 2}, {2, 1}}
// 			var peopleIDs = [][]IDType{{1001, 1002}, {1002, 1001}}

// 			var purchases []*Purchase
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 2, Offset: 3}
// 				expectQuery_Purchase_All(mock, releaseID, ethereumAddress, purchaseIDs, musicrecordingIDs, peopleIDs, selectParams)
// 				purchases, err = (&Purchase{}).All(dbx, releaseID, ethereumAddress, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected purchases", func() {
// 				Expect(purchases).To(HaveLen(len(purchaseIDs)))

// 				purchaseMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, purchaseID := range purchaseIDs {
// 					purchaseMatchers = append(purchaseMatchers, PointTo(matchPurchaseFixture(purchaseID)))
// 				}
// 				Expect(purchases).To(ConsistOf(purchaseMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})
// 	})

// 	Context("when .Get is called", func() {
// 		var purchaseID IDType = 264
// 		var musicrecordingIDs = []IDType{1, 2}
// 		var peopleIDs = []IDType{1001, 1002}
// 		var purchase *Purchase
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_Purchase_Get(mock, purchaseID, musicrecordingIDs, peopleIDs)
// 			purchase, err = (&Purchase{}).Get(dbx, purchaseID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil purchase", func() {
// 			Expect(purchase).ToNot(BeNil())
// 			Expect(*purchase).To(matchPurchaseFixture(purchaseID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Create is called", func() {
// 		var p *Purchase
// 		var newPurchaseID IDType = 5557
// 		var returnedPurchase *Purchase
// 		var err error

// 		BeforeEach(func() {
// 			p = &Purchase{
// 				CID:          "_cid",
// 				Type:         "_type",
// 				Context:      "_context",
// 				TxHash:       "0xdeadbeef",
// 				Buyer:        &Person{ID: 321},
// 				MusicRelease: &MusicRelease{ID: 123},
// 			}

// 			expectQuery_Purchase_Create(mock, p, newPurchaseID)
// 			returnedPurchase, err = p.Create(dbx)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil person with its ID field set", func() {
// 			Expect(returnedPurchase).ToNot(BeNil())
// 			Expect(returnedPurchase.ID).To(Equal(newPurchaseID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Update is called", func() {
// 		var p *Purchase
// 		var err error

// 		BeforeEach(func() {
// 			err = p.Update(dbx)
// 		})

// 		It("should return an error saying that the method is not implemented", func() {
// 			Expect(err).ToNot(BeNil())
// 		})
// 	})

// 	Context("when .Delete is called", func() {
// 		var p *Purchase
// 		var purchaseID IDType = 5557
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_Purchase_Delete(mock, p, purchaseID)
// 			err = p.Delete(dbx, purchaseID)
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

// func expectQuery_Purchase_GetMusicRelease(mock sqlmock.Sqlmock, musicreleaseID IDType, musicrecordingIDs []IDType, peopleIDs []IDType) {
// 	var musicrelease = musicreleases[musicreleaseID]

// 	mock.ExpectQuery(`SELECT (.+) FROM musicrelease`).
// 		WithArgs(musicreleaseID).
// 		WillReturnRows(
// 			mockResultRows(
// 				musicrelease,
// 			),
// 		)

// 	var musicalbumID = musicrelease["release_of_id"].(IDType)
// 	expectQuery_MusicAlbum_Get(mock, musicalbumID, musicrecordingIDs, peopleIDs)
// 	expectQuery_MusicRelease_GetImage(mock, musicrelease["image_id"].(IDType))
// }

// func expectQuery_Purchase_All(mock sqlmock.Sqlmock, releaseID *IDType, ethereumAddress string, purchaseIDs []IDType, musicrecordingIDs [][]IDType, peopleIDs [][]IDType, selectParams *SelectQuery) {
// 	purchaseRows := make([]map[string]driver.Value, len(purchaseIDs))
// 	for i, id := range purchaseIDs {
// 		purchaseRows[i] = purchases[id]
// 	}

// 	var args []driver.Value
// 	if ethereumAddress != "" && releaseID != nil {
// 		args = append(args, *releaseID, ethereumAddress)
// 	} else if ethereumAddress != "" && releaseID == nil {
// 		args = append(args, ethereumAddress)
// 	}

// 	if selectParams.Limit > 0 {
// 		args = append(args, selectParams.Limit, selectParams.Offset)
// 	}

// 	e := mock.ExpectQuery(`SELECT (.+) FROM purchase`)
// 	if len(args) > 0 {
// 		e = e.WithArgs(args...)
// 	}
// 	e = e.WillReturnRows(
// 		mockResultRows(
// 			purchaseRows...,
// 		),
// 	)

// 	for i, purchase := range purchaseRows {
// 		expectQuery_Purchase_GetMusicRelease(mock, purchase["musicrelease_id"].(IDType), musicrecordingIDs[i], peopleIDs[i])
// 		expectQuery_Person_Get(mock, purchase["buyer_id"].(IDType), "")
// 	}
// }

// func expectQuery_Purchase_Get(mock sqlmock.Sqlmock, purchaseID IDType, musicrecordingIDs []IDType, peopleIDs []IDType) {
// 	var purchase = purchases[purchaseID]

// 	mock.ExpectQuery(`SELECT (.+) FROM purchase`).
// 		WithArgs(purchaseID).
// 		WillReturnRows(
// 			mockResultRows(
// 				purchase,
// 			),
// 		)

// 	expectQuery_Purchase_GetMusicRelease(mock, purchase["musicrelease_id"].(IDType), musicrecordingIDs, peopleIDs)
// 	expectQuery_Person_Get(mock, purchase["buyer_id"].(IDType), "")
// }

// func expectQuery_Purchase_Create(mock sqlmock.Sqlmock, p *Purchase, newPurchaseID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`INSERT INTO purchase`).
// 		WithArgs(p.CID, p.Type, p.Context, AnyTime{}, AnyTime{}, p.TxHash, p.Buyer.ID, p.MusicRelease.ID).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newPurchaseID}))

// 	mock.ExpectCommit()
// }

// func expectQuery_Purchase_Delete(mock sqlmock.Sqlmock, p *Purchase, purchaseID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectExec(`DELETE FROM purchase`).
// 		WithArgs(purchaseID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectCommit()
// }
