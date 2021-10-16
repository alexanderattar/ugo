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

// var _ = Describe("Report", func() {
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
// 		Context("with a non-nil musicreleaseID argument", func() {
// 			var musicreleaseID = idPtr(818)
// 			var reportIDs = []IDType{582, 592}
// 			var musicrecordingIDs = []IDType{1, 2}
// 			var peopleIDs = []IDType{1001, 1002}

// 			var reports []*Report
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 32, Offset: 1}
// 				expectQuery_Report_All(mock, musicreleaseID, reportIDs, musicrecordingIDs, peopleIDs, selectParams)
// 				reports, err = (&Report{}).All(dbx, musicreleaseID, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected reports", func() {
// 				Expect(reports).To(HaveLen(len(reportIDs)))

// 				reportMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, reportID := range reportIDs {
// 					reportMatchers = append(reportMatchers, PointTo(matchReportFixture(reportID)))
// 				}
// 				Expect(reports).To(ConsistOf(reportMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})

// 		Context("with a nil musicreleaseID argument", func() {
// 			var musicreleaseID *IDType = nil
// 			var reportIDs = []IDType{582, 592}
// 			var musicrecordingIDs = []IDType{1, 2}
// 			var peopleIDs = []IDType{1001, 1002}

// 			var reports []*Report
// 			var err error

// 			BeforeEach(func() {
// 				selectParams := &SelectQuery{Limit: 32, Offset: 1}
// 				expectQuery_Report_All(mock, musicreleaseID, reportIDs, musicrecordingIDs, peopleIDs, selectParams)
// 				reports, err = (&Report{}).All(dbx, musicreleaseID, selectParams)
// 			})

// 			It("should return a nil error", func() {
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a non-nil slice containing the expected reports", func() {
// 				Expect(reports).To(HaveLen(len(reportIDs)))

// 				reportMatchers := []gomegatypes.GomegaMatcher{}
// 				for _, reportID := range reportIDs {
// 					reportMatchers = append(reportMatchers, PointTo(matchReportFixture(reportID)))
// 				}
// 				Expect(reports).To(ConsistOf(reportMatchers))
// 			})

// 			It("should execute the expected SQL queries", func() {
// 				if err = mock.ExpectationsWereMet(); err != nil {
// 					Fail("there were unfulfilled expectations: " + err.Error())
// 				}
// 			})
// 		})
// 	})

// 	Context("when .Get is called", func() {
// 		var reportID IDType = 582
// 		var musicrecordingIDs = []IDType{1, 2}
// 		var peopleIDs = []IDType{1001, 1002}
// 		var report *Report
// 		var err error

// 		BeforeEach(func() {
// 			expectQuery_Report_Get(mock, reportID, musicrecordingIDs, peopleIDs)
// 			report, err = (&Report{}).Get(dbx, reportID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should return a non-nil report", func() {
// 			Expect(report).ToNot(BeNil())
// 			Expect(*report).To(matchReportFixture(reportID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Create is called", func() {
// 		var newReportID IDType = 126
// 		var musicrecordingIDs = []IDType{1, 2}
// 		var peopleIDs = []IDType{1001, 1002}
// 		var returnedReportID IDType
// 		var r *Report
// 		var err error

// 		BeforeEach(func() {
// 			r = &Report{
// 				State:          "",
// 				Reason:         "_reason",
// 				Message:        strPtr("_message"),
// 				Email:          strPtr("_email"),
// 				MusicReleaseID: 818,
// 				ReporterID:     44,
// 			}

// 			expectQuery_Report_Create(mock, r, newReportID, musicrecordingIDs, peopleIDs)
// 			returnedReportID, err = r.Create(dbx)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should set the Report.State field to 'unreviewed'", func() {
// 			Expect(r.State).To(Equal("unreviewed"))
// 		})

// 		It("should populate the Report.MusicRelease field", func() {
// 			Expect(r.MusicRelease).To(PointTo(matchMusicreleaseFixture(r.MusicReleaseID)))
// 		})

// 		It("should return the new report ID", func() {
// 			Expect(returnedReportID).To(Equal(newReportID))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Update is called", func() {
// 		var reportID IDType = 818
// 		var musicrecordingIDs = []IDType{1, 2}
// 		var peopleIDs = []IDType{1001, 1002}
// 		var r *Report
// 		var err error

// 		BeforeEach(func() {
// 			r = &Report{
// 				State:          "",
// 				Reason:         "_reason",
// 				Message:        strPtr("_message"),
// 				Email:          strPtr("_email"),
// 				MusicReleaseID: 818,
// 				ReporterID:     44,
// 			}

// 			expectQuery_Report_Update(mock, reportID, r, r.State, musicrecordingIDs, peopleIDs)
// 			_, err = r.Update(dbx, reportID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should populate the Report.MusicRelease field", func() {
// 			Expect(r.MusicRelease).To(PointTo(matchMusicreleaseFixture(r.MusicReleaseID)))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	Context("when .Resolve is called", func() {
// 		var reportID IDType = 818
// 		var musicrecordingIDs = []IDType{1, 2}
// 		var peopleIDs = []IDType{1001, 1002}
// 		var r *Report
// 		var err error

// 		BeforeEach(func() {
// 			r = &Report{
// 				State:          "",
// 				Reason:         "_reason",
// 				Message:        strPtr("_message"),
// 				Email:          strPtr("_email"),
// 				MusicReleaseID: 818,
// 				ReporterID:     44,
// 			}

// 			expectQuery_Report_Update(mock, reportID, r, "resolved", musicrecordingIDs, peopleIDs)
// 			_, err = r.Resolve(dbx, reportID)
// 		})

// 		It("should return a nil error", func() {
// 			Expect(err).To(BeNil())
// 		})

// 		It("should set the Report.State field to 'resolved'", func() {
// 			Expect(r.State).To(Equal("resolved"))
// 		})

// 		It("should execute the expected SQL queries", func() {
// 			if err = mock.ExpectationsWereMet(); err != nil {
// 				Fail("there were unfulfilled expectations: " + err.Error())
// 			}
// 		})
// 	})

// 	// Context("when .Deactivate is called", func() {
// 	//  var reportID IDType = 818
// 	//  var musicrecordingIDs = []IDType{1, 2}
// 	//  var peopleIDs = []IDType{1001, 1002}
// 	//  var r *Report
// 	//  var err error

// 	//  BeforeEach(func() {
// 	//      r = &Report{
// 	//          State:          "",
// 	//          Reason:         "_reason",
// 	//          Message:        strPtr("_message"),
// 	//          Email:          strPtr("_email"),
// 	//          MusicReleaseID: 818,
// 	//          ReporterID:     44,
// 	//      }
// 	//      mr, err := musicrelease_GetByID_simulate(r.MusicReleaseID, musicrecordingIDs, peopleIDs)
// 	//      if err != nil {
// 	//          Fail(err.Error())
// 	//      }

// 	//      expectQuery_MusicRelease_GetByID(mock, r.MusicReleaseID, musicrecordingIDs, peopleIDs)
// 	//      expectQuery_MusicRelease_Update(mock, r.MusicReleaseID, mr)

// 	//      // expectQuery_Report_Update(mock, reportID, r, "deactivated", musicrecordingIDs, peopleIDs)
// 	//      _, err = r.Deactivate(dbx, reportID)
// 	//  })

// 	//  It("should return a nil error", func() {
// 	//      Expect(err).To(BeNil())
// 	//  })

// 	//  It("should set the Report.State field to 'deactivated'", func() {
// 	//      Expect(r.State).To(Equal("deactivated"))
// 	//  })

// 	//  It("should execute the expected SQL queries", func() {
// 	//      if err = mock.ExpectationsWereMet(); err != nil {
// 	//          Fail("there were unfulfilled expectations: " + err.Error())
// 	//      }
// 	//  })
// 	// })
// })

// func expectQuery_Report_All(mock sqlmock.Sqlmock, musicreleaseID *IDType, reportIDs, musicrecordingIDs, peopleIDs []IDType, selectParams *SelectQuery) {
// 	reportRows := make([]map[string]driver.Value, len(reportIDs))
// 	for i, id := range reportIDs {
// 		reportRows[i] = reports[id]
// 	}

// 	var args []driver.Value
// 	if musicreleaseID != nil {
// 		args = append(args, *musicreleaseID)
// 	}
// 	if selectParams.Limit > 0 {
// 		args = append(args, selectParams.Limit, selectParams.Offset)
// 	}

// 	e := mock.ExpectQuery(`SELECT (.+) FROM report`)
// 	if len(args) > 0 {
// 		e = e.WithArgs(args...)
// 	}
// 	e = e.WillReturnRows(
// 		mockResultRows(
// 			reportRows...,
// 		),
// 	)

// 	for _, reportID := range reportIDs {
// 		report := reports[reportID]
// 		expectQuery_MusicRelease_GetByID(mock, report["musicrelease_id"].(IDType), musicrecordingIDs, peopleIDs)
// 	}
// }

// func expectQuery_Report_Get(mock sqlmock.Sqlmock, reportID IDType, musicrecordingIDs, peopleIDs []IDType) {
// 	var report = reports[reportID]

// 	mock.ExpectQuery("SELECT (.+) FROM report").
// 		WithArgs(reportID).
// 		WillReturnRows(
// 			mockResultRows(
// 				report,
// 			),
// 		)

// 	expectQuery_MusicRelease_GetByID(mock, report["musicrelease_id"].(IDType), musicrecordingIDs, peopleIDs)
// }

// func expectQuery_Report_Create(mock sqlmock.Sqlmock, r *Report, newReportID IDType, musicrecordingIDs, peopleIDs []IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`INSERT INTO report`).
// 		WithArgs(AnyTime{}, AnyTime{}, "unreviewed", r.Reason, r.Message, r.Email, r.MusicReleaseID, r.ReporterID).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": newReportID}))

// 	expectQuery_MusicRelease_GetByID(mock, r.MusicReleaseID, musicrecordingIDs, peopleIDs)

// 	mock.ExpectCommit()
// }

// func expectQuery_Report_Update(mock sqlmock.Sqlmock, reportID IDType, r *Report, state string, musicrecordingIDs, peopleIDs []IDType) {
// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`UPDATE report`).
// 		WithArgs(AnyTime{}, AnyTime{}, state, r.Reason, r.Message, r.Response, r.Email, r.MusicReleaseID, r.ReporterID, reportID).
// 		WillReturnRows(mockResultRows(map[string]driver.Value{"id": reportID}))

// 	expectQuery_MusicRelease_GetByID(mock, r.MusicReleaseID, musicrecordingIDs, peopleIDs)

// 	mock.ExpectCommit()
// }
