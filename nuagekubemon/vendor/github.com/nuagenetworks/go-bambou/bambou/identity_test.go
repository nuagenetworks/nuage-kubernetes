package bambou

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIdentity_AllIdentity(t *testing.T) {

	Convey("Given I retrieve the AllIdentity", t, func() {
		i := AllIdentity

		Convey("Then Name should __all__", func() {
			So(i.Name, ShouldEqual, "__all__")
		})

		Convey("Then Category should __all__", func() {
			So(i.Category, ShouldEqual, "__all__")
		})
	})
}

func TestIdentity_NewIdentity(t *testing.T) {

	Convey("Given I create a new identity", t, func() {
		i := Identity{"n", "c"}

		Convey("Then Name should n", func() {
			So(i.Name, ShouldEqual, "n")
		})

		Convey("Then Category should c", func() {
			So(i.Category, ShouldEqual, "c")
		})
	})
}

func TestIdentity_String(t *testing.T) {

	Convey("Given I create a new identity", t, func() {
		i := Identity{"n", "c"}

		Convey("Then String should <Identity n|c>", func() {
			So(i.String(), ShouldEqual, "<Identity n|c>")
		})
	})
}
