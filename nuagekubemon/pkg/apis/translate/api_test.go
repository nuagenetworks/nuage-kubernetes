package translate

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestInitVSDObjsMap(t *testing.T) {
	g := NewGomegaWithT(t)
	p := InitVSDObjsMap()
	g.Expect(p).NotTo(BeNil())
	g.Expect(p.PGMap).NotTo(BeNil())
	g.Expect(p.NSLabelsMap).NotTo(BeNil())
	g.Expect(p.NWMacroMap).NotTo(BeNil())
}
