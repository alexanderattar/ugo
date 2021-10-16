package models_test

// import (
// 	"database/sql/driver"
// 	"testing"
// 	"time"

// 	"github.com/brynbellomy/ginkgo-reporter"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	"github.com/onsi/gomega/format"

// 	"github.com/sirupsen/logrus"

// 	"github.com/consensys/ugo/pkg/lg"
// 	"github.com/consensys/ugo/pkg/models"
// )

// func TestModels(t *testing.T) {
// 	logger := logrus.New()
// 	lg.RedirectStdlogOutput(logger)
// 	lg.DefaultLogger = logger

// 	RegisterFailHandler(Fail)
// 	RunSpecsWithCustomReporters(t, "Models Suite", []Reporter{
// 		&reporter.TerseReporter{Logger: &reporter.DefaultLogger{}},
// 	})
// }

// func strPtr(s string) *string               { return &s }
// func f64Ptr(f float64) *float64             { return &f }
// func boolPtr(b bool) *bool                  { return &b }
// func idPtr(id models.IDType) *models.IDType { return &id }

// //
// // AnyTime is a matcher for an argument to a SQL query that should be a time.Time.  It satisfies the
// // sqlmock.Argument interface.
// //
// type AnyTime struct{}

// func (a AnyTime) Match(v driver.Value) bool {
// 	_, ok := v.(time.Time)
// 	return ok
// }

// //
// // AnyMatcher is a gomega matcher that matches anything you pass to it.  We use this with the
// // gomega/gstruct.Field matcher to ensure that the tests explicitly specify all struct fields while
// // not requiring all of them to contain a particular value.
// //
// type AnyMatcher struct{}

// func (matcher *AnyMatcher) Match(actual interface{}) (success bool, err error) {
// 	return true, nil
// }

// func (matcher *AnyMatcher) FailureMessage(actual interface{}) string {
// 	return format.Message(actual, "to equal anything")
// }

// func (matcher *AnyMatcher) NegatedFailureMessage(actual interface{}) string {
// 	return format.Message(actual, "not to equal anything (??!?!)")
// }

// func BeAnything() *AnyMatcher {
// 	return &AnyMatcher{}
// }
