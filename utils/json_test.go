package utils

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestStructToJson(t *testing.T) {
	type testStruct struct {
		UserId core.UserId `json:"user_id"`
		Price  core.Price  `json:"price"`
	}
	type testCase struct {
		data   testStruct
		expect string
	}

	testCases := []*testCase{
		{
			data: testStruct{
				UserId: "test_user_id",
				Price:  123456,
			},
			expect: `{"user_id":"test_user_id","price":123456}`,
		},
	}

	for _, v := range testCases {
		res, err := StructToJson(&v.data)
		if err != nil {
			t.Fatalf("error: %s", err.Error())
		}
		if *res != v.expect {
			t.Errorf("expect: %s, got: %s", v.expect, *res)
		}
	}
}

func TestJsonToStruct(t *testing.T) {
	type testStruct struct {
		UserId core.UserId `json:"user_id"`
		Price  core.Price  `json:"price"`
	}
	type testCase struct {
		expect testStruct
		text   string
	}

	testCases := []*testCase{
		{
			expect: testStruct{
				UserId: "test_user_id",
				Price:  123456,
			},
			text: `{"user_id":"test_user_id","price":123456}`,
		},
	}
	for _, v := range testCases {
		res, err := JsonToStruct[testStruct](v.text)
		if err != nil {
			t.Fatalf("error: %s", err.Error())
		}
		if test.DeepEqual(res, v.expect) {
			t.Errorf("expect: %+v, got: %+v", v.expect, res)
		}

	}
}
