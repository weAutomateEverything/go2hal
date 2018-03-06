package halmock

import (
	"github.com/golang/mock/gomock"
	"strings"
)

type errorMsgMatcher struct {
	x interface{}
}

//ErrorMsgMatches returns a new go mock matcher that returns true if both items are errors and their messages match
func ErrorMsgMatches(x interface{}) gomock.Matcher {
	return errorMsgMatcher{x: x}
}

//Matches returns true if the two objects are both errors and their error messages match
func (e errorMsgMatcher) Matches(x interface{}) bool {
	xerr, ok := e.x.(error)
	if !ok {
		panic("Can only be used with error messages")
	}

	ierr, ok := x.(error)
	if !ok {
		panic("Can only be used with error messages ")
	}

	return strings.Compare(xerr.Error(), ierr.Error()) == 0

}

func (e errorMsgMatcher) String() string {
	xerr, ok := e.x.(error)
	if !ok {
		panic("Can only be used with error messages ")
	}
	return xerr.Error()
}
