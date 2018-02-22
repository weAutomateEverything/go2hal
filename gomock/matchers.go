package gomock

import "github.com/golang/mock/gomock"

type errorMsgMatcher struct {
	x interface{}
}

func ErrorMsgMatches(x interface{}) gomock.Matcher{
	return errorMsgMatcher{x:x}
}

func (e errorMsgMatcher) Matches(x interface{}) bool {
	xerr, ok := e.x.(error)
	if !ok {
		return false
	}

	ierr, ok := x.(error)
	if !ok{
		return false
	}

	return xerr.Error() == ierr.Error()

}

func (e errorMsgMatcher) String() string {
	return "error messages match"
}


