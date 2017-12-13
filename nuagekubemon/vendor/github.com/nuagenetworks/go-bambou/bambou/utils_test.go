package bambou

import "fmt"

/*
   Fake Exposed
*/

var FakeIdentity = Identity{"fake", "fakes"}

type FakeObjectsList []*FakeObject

type FakeObject struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func NewFakeObject(ID string) *FakeObject {

	return &FakeObject{ID: ID}
}

func (o *FakeObject) Identity() Identity      { return FakeIdentity }
func (o *FakeObject) Identifier() string      { return o.ID }
func (o *FakeObject) SetIdentifier(ID string) { o.ID = ID }

/*
   Fake Rootable
*/
var FakeRootIdentity = Identity{"root", "root"}

type FakeRootObject struct {
	FakeObject

	Token string `json:"APIKey,omitempty"`
}

func NewFakeRootObject() *FakeRootObject {
	return &FakeRootObject{}
}

func (o *FakeRootObject) Identity() Identity   { return FakeRootIdentity }
func (o *FakeRootObject) APIKey() string       { return o.Token }
func (o *FakeRootObject) SetAPIKey(key string) { o.Token = key }

/*
   Unmarshalable object
*/

var FakeUnmarshalableObjectIdentity = Identity{"unmarshalable", "unmarshalable"}

type UnmarshalableFakeObject struct {
	FakeObject
}

func NewUnmarshalableFakeObject(ID string) *UnmarshalableFakeObject {
	return &UnmarshalableFakeObject{
		FakeObject: FakeObject{
			ID: ID,
		},
	}
}

func (o *UnmarshalableFakeObject) Identity() Identity { return FakeUnmarshalableObjectIdentity }

func (o *UnmarshalableFakeObject) UnmarshalJSON([]byte) error {
	return fmt.Errorf("error unmarshalling")
}

func (o *UnmarshalableFakeObject) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("error marshalling")
}
