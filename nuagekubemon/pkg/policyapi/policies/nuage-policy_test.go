package policies

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewNuagePolicy(t *testing.T) {
	g := NewGomegaWithT(t)

	enterprise := "test-enterprise"
	domain := "test-domain"
	name := "test-name"
	id := "some-random-id"
	priority := 10

	p, err := NewNuagePolicy(enterprise, domain, name, id, priority)

	g.Expect(err).To(BeNil())
	g.Expect(p.Enterprise).To(Equal(enterprise))
	g.Expect(p.Domain).To(Equal(domain))
	g.Expect(p.Name).To(Equal(name))
	g.Expect(p.ID).To(Equal(id))
	g.Expect(p.Priority).To(Equal(priority))

	p, err = NewNuagePolicy("", domain, name, id, priority)
	g.Expect(err).To(MatchError("enterprise is empty"))

	p, err = NewNuagePolicy(enterprise, "", name, id, priority)
	g.Expect(err).To(MatchError("domain is empty"))

	p, err = NewNuagePolicy(enterprise, domain, "", id, priority)
	g.Expect(err).To(MatchError("name is empty"))

	p, err = NewNuagePolicy(enterprise, domain, name, "", priority)
	g.Expect(err).To(MatchError("id is empty"))

	p, err = NewNuagePolicy(enterprise, domain, name, id, -1)
	g.Expect(err).To(MatchError("priority cannot be negative"))
}
