package test_tools

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"time"
)

func DiffEq(v interface{}, opts ...cmp.Option) gomock.Matcher {
	var lopts cmp.Options
	if len(opts) == 0 {
		lopts = append(lopts, cmpopts.EquateApproxTime(1*time.Second))
	} else {
		lopts = append(lopts, opts...)
	}
	return &diffMatcher{want: v, opts: lopts}
}

type diffMatcher struct {
	want interface{}
	diff string
	opts cmp.Options
}

// Matches implements golang/mock/gomock#Matcher interface.
func (d *diffMatcher) Matches(x interface{}) bool {
	d.diff = cmp.Diff(x, d.want, d.opts...)
	return len(d.diff) == 0
}

// String implements golang/mock/gomock#Matcher interface.
func (d *diffMatcher) String() string {
	if d.diff == "" {
		return ""
	}
	return fmt.Sprintf("diff(-got +want) is %s", d.diff)
}
