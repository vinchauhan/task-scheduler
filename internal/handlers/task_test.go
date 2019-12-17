package handlers_test

import (
	"github.com/matryer/way"
	"net/http/httptest"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCompleteTask(t *testing.T) {
	Convey("Given a HTTP request for /task", t, func() {
		req := httptest.NewRequest("POST", "/task", nil)
		resp := httptest.NewRecorder()
		Convey("When the request is handled by the Router", func() {
			way.NewRouter().ServeHTTP(resp, req)
			Convey("Then the response should be a 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
		})

	})
}
