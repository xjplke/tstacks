package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetMapPath(t *testing.T) {
	Convey("Test Register", t, func() {
		tmp := map[string]interface{}{
			"aaa": "123",
			"bbb": 123,
			"ccc": map[string]interface{}{
				"aaa": "bbb",
				"bbb": 123,
				"ccc": map[string]interface{}{
					"aaa": "aaa",
					"bbb": "bbb",
				},
			},
		}

		a, err := GetMapPath(tmp, "aaa")
		So(a, ShouldNotEqual, nil)
		So(err, ShouldEqual, nil)
		b, err := GetMapPath(tmp, "ccc", "bbb")
		So(b, ShouldNotEqual, nil)
		So(err, ShouldEqual, nil)
		c, err := GetMapPath(tmp, "ccc", "ccc", "aaa")
		So(err, ShouldEqual, nil)
		So(c, ShouldNotEqual, nil)
		So(c, ShouldEqual, "aaa")

		d, err := GetMapPath(tmp, "")
		So(d, ShouldEqual, nil)
		So(err, ShouldNotEqual, nil)

	})
}
