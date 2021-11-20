package controllers

import (
	"testing"
)

func TestGet(t *testing.T) {
	if ret, err := Get(3, 5); nil != err {
		t.Error(err)
	} else {
		if len(ret) < 5 {
			t.Error("Can not get result")
		}
	}
}

func TestGetMostLike(t *testing.T) {
	if ret, err := GetMostLike(6, 5); nil != err {
		t.Error(err)
	} else {
		if len(ret) != 5 {
			t.Error("Can not get enough result:", len(ret))
		}
		for k, v := range ret {
			if k == 0 {
				continue
			}
			if ret[k-1].MessageCount.Count < v.MessageCount.Count {
				t.Error("cannot get most like list k-1=", ret[k-1].MessageCount.Count, " v=,", v.MessageCount.Count)
			}
		}
	}
}
