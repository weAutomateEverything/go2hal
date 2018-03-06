package halmock

import (
	"github.com/golang/mock/gomock"
	"strings"
)

type errorMsgMatcher struct {
	x interface{}
}

func ErrorMsgMatches(x interface{}) gomock.Matcher{
	return errorMsgMatcher{x:x}
}

func (e errorMsgMatcher) Matches(x interface{}) bool {
	xerr, ok := e.x.(error)
	if !ok {
		panic("Can only be used with error messages")
	}

	ierr, ok := x.(error)
	if !ok{
		panic("Can only be used with error messages ")
	}


	return strings.Compare(xerr.Error(), ierr.Error()) == 0

}

func (e errorMsgMatcher) String() string {
	xerr, ok := e.x.(error)
	if !ok{
		panic("Can only be used with error messages ")
	}
	return xerr.Error()
}


