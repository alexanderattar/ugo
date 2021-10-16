package models_test

// import (
// 	"database/sql"
// 	"database/sql/driver"
// 	"time"

// 	"github.com/jmoiron/sqlx"
// 	// "github.com/lib/pq"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/onsi/gomega/gstruct"
// 	gomegatypes "github.com/onsi/gomega/types"
// 	"gopkg.in/DATA-DOG/go-sqlmock.v1"

// 	. "github.com/consensys/ugo/pkg/models"
// )

// var _ = Describe("Person", func() {
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
// 			var peopleIDs = []IDType{1001, 1002}

// 			var people []*Person
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 2, Offset: 3}
// 				expectQuery_Person_All(mock, ethereumAddress, peopleIDs, selectParams)
// 				people, err = (&Person{}).All(dbx, ethereumAddress, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected people", func() {
// 				Expect(people).To(HaveLen(len(peopleIDs)))

// 				peopleMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, personID := range peopleIDs {
// 					peopleMatchers = append(peopleMatchers, PointTo(matchPersonFixture(personID)))
// 				}
// 				Expect(people).To(ConsistOf(peopleMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with an empty ethereumAddress argument", func() {
// 			var ethereumAddress = ""
// 			var peopleIDs = []IDType{1001, 1002}

// 			var people []*Person
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 2, Offset: 3}
// 				expectQuery_Person_All(mock, ethereumAddress, peopleIDs, selectParams)
// 				people, err = (&Person{}).All(dbx, ethereumAddress, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected people", func() {
// 				Expect(people).To(HaveLen(len(peopleIDs)))

// 				peopleMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, personID := range peopleIDs {
// 					peopleMatchers = append(peopleMatchers, PointTo(matchPersonFixture(personID)))
// 				}
// 				Expect(people).To(ConsistOf(peopleMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})
// 	})

// 	Context("when .Get is called", func() {
// 		Context("with a non-empty ethereumAddress argument", func() {
// 			var ethereumAddress = "0xdeadbeef"
// 			var personID IDType = 1001

// 			var person *Person
// 			var err error

// 			BeforeEach(func() {
// 				expectQuery_Person_Get(mock, personID, ethereumAddress)
// 				person, err = (&Person{}).Get(dbx, personID, ethereumAddress)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil person", func() {
// 				Expect(person).ToNot(BeNil())
// 				Expect(*person).To(matchPersonFixture(personID))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with an empty ethereumAddress argument", func() {
// 			var ethereumAddress = ""
// 			var personID IDType = 1001

// 			var person *Person
// 			var err error

// 			BeforeEach(func() {
// 				expectQuery_Person_Get(mock, personID, ethereumAddress)
// 				person, err = (&Person{}).Get(dbx, personID, ethereumAddress)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil person", func() {
// 				Expect(person).ToNot(BeNil())
// 				Expect(*person).To(matchPersonFixture(personID))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})
// 	})

// 	Context("when .Create is called", func() {
// 		Context("with valid values", func() {
// 			var p *Person
// 			var newPersonID IDType = 5557
// 			var returnedPersonID IDType
// 			var err error

// 			BeforeEach(func() {
// 				p = makePersonStruct()

// 				expectQuery_Person_Create(mock, p, newPersonID)
// 				returnedPersonID, err = p.Create(dbx)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return the given personID", func() {
// 				Expect(returnedPersonID).To(Equal(newPersonID))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err := mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with an invalid Ethereum address", func() {
// 			var p *Person
// 			var newPersonID IDType = 5557
// 			var err error

// 			BeforeEach(func() {
// 				p = makePersonStruct()
// 				p.EthereumAddress = "xyzzy"

// 				expectQuery_Person_Create(mock, p, newPersonID)
// 				_, err = p.Create(dbx)
// 			})

// 			It("should return a non-nil error", func() {
// 				Expect(err).ToNot(BeNil())
// 			})
// 		})

// 		Context("with an invalid email address", func() {
// 			var p *Person
// 			var newPersonID IDType = 5557
// 			var err error

// 			BeforeEach(func() {
// 				p = makePersonStruct()
// 				p.Email = strPtr("xyzzy")

// 				expectQuery_Person_Create(mock, p, newPersonID)
// 				_, err = p.Create(dbx)
// 			})

// 			It("should return a non-nil error", func() {
// 				Expect(err).ToNot(BeNil())
// 			})
// 		})
// 	})

// 	Context("when .Update is called", func() {
// 		Context("with valid values", func() {
// 			var p *Person
// 			var personID IDType = 5557
// 			var returnedPersonID IDType
// 			var err error

// 			BeforeEach(func() {
// 				p = makePersonStruct()

// 				expectQuery_Person_Update(mock, p, personID)
// 				returnedPersonID, err = p.Update(dbx, personID)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return the given personID", func() {
// 				Expect(returnedPersonID).To(Equal(personID))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err := mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with an invalid Ethereum address", func() {
// 			var p *Person
// 			var personID IDType = 5557
// 			var err error

// 			BeforeEach(func() {
// 				p = makePersonStruct()
// 				p.EthereumAddress = "xyzzy"

// 				expectQuery_Person_Update(mock, p, personID)
// 				_, err = p.Update(dbx, personID)
// 			})

// 			It("should return a non-nil error", func() {
// 				Expect(err).ToNot(BeNil())
// 			})
// 		})

// 		Context("with an invalid email address", func() {
// 			var p *Person
// 			var personID IDType = 5557
// 			var err error

// 			BeforeEach(func() {
// 				p = makePersonStruct()
// 				p.Email = strPtr("xyzzy")

// 				expectQuery_Person_Update(mock, p, personID)
// 				_, err = p.Update(dbx, personID)
// 			})

// 			It("should return a non-nil error", func() {
// 				Expect(err).ToNot(BeNil())
// 			})
// 		})
// 	})

// 	Context("when .Delete is called", func() {
// 		var p *Person
// 		var personID IDType = 5557
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_Person_Delete(mock, p, personID)
// 			err = p.Delete(dbx, personID)
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

// func expectQuery_Person_All(mock sqlmock.Sqlmock, ethereumAddress string, peopleIDs []IDType, selectParams *SelectQuery) {
// 	personRows := make([]map[string]driver.Value, len(peopleIDs))
// 	for i, id := range peopleIDs {
// 		personRows[i] = people[id]
// 	}

// 	var args []driver.Value
// 	if ethereumAddress != "" {
// 		args = append(args, ethereumAddress)
// 	}
// 	if selectParams.Limit > 0 {
// 		args = append(args, selectParams.Limit, selectParams.Offset)
// 	}

// 	e := mock.ExpectQuery(`SELECT (.+) FROM person`)
// 	if len(args) > 0 {
// 		e = e.WithArgs(args...)
// 	}
// 	e = e.WillReturnRows(
// 		mockResultRows(
// 			personRows...,
// 		),
// 	)
// }

// func expectQuery_Person_Get(mock sqlmock.Sqlmock, personID IDType, ethereumAddress string) {
// 	var args []driver.Value
// 	if ethereumAddress != "" {
// 		args = []driver.Value{ethereumAddress}
// 	} else {
// 		args = []driver.Value{personID}
// 	}

// 	mock.ExpectQuery(`SELECT (.+) FROM person`).
// 		WithArgs(args...).
// 		WillReturnRows(
// 			mockResultRows(
// 				people[personID],
// 			),
// 		)
// }

// func expectQuery_Person_Create(mock sqlmock.Sqlmock, p *Person, newPersonID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`INSERT INTO person`).
// 		WithArgs(p.CID, p.Type, p.Context, AnyTime{}, AnyTime{}, p.EthereumAddress, p.GivenName, p.FamilyName, p.Email).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newPersonID}))

// 	mock.ExpectCommit()
// }

// func expectQuery_Person_Update(mock sqlmock.Sqlmock, p *Person, personID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`UPDATE person`).
// 		WithArgs(p.CID, p.Type, p.Context, AnyTime{}, p.EthereumAddress, p.GivenName, p.FamilyName, p.Email, personID).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": personID}))

// 	mock.ExpectCommit()
// }

// func expectQuery_Person_Delete(mock sqlmock.Sqlmock, p *Person, personID IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectExec(`DELETE FROM person`).
// 		WithArgs(personID).
// 		WillReturnResult(sqlmock.NewResult(123, 1))

// 	mock.ExpectCommit()
// }

// func makePersonStruct() *Person {
// 	return &Person{
// 		CID:             "_cid",
// 		Type:            "_type",
// 		Context:         "_context",
// 		EthereumAddress: "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
// 		GivenName:       strPtr("_given_name"),
// 		FamilyName:      strPtr("_family_name"),
// 		Email:           strPtr("xyzzy@zork.org"),
// 		Image: &ImageObject{
// 			ID:             551,
// 			CID:            "_cid",
// 			Type:           "_type",
// 			Context:        "_context",
// 			CreatedAt:      time.Now(),
// 			UpdatedAt:      time.Now(),
// 			ContentURL:     strPtr("_contentURL"),
// 			EncodingFormat: strPtr("_encodingFormat"),
// 		},
// 		Description:      strPtr("_description"),
// 		PercentageShares: f64Ptr(0.125),
// 		MusicGroupAdmin:  boolPtr(true),
// 	}
// }
