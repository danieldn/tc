// +build mock

package mockpkg

import "testing"

func TestMockSucceed(t *testing.T) {
	t.Log("Checking if mocking succeed works")

	expected := 2
	ans := PlusOne(1)

	if ans != expected {
		t.Errorf("Got %v expected %v", ans, expected)
	}
	t.Log("Mocking succeed confirmed")
}

func TestMockFail(t *testing.T) {
	expected := 2
	ans := BuggyPlusOne(1)
	if ans != expected {
		t.Errorf("Got %v expected %v", ans, expected)
	}
}

func TestMockSkip(t *testing.T) {
	t.Skip("skipping this test...")

}

func TestMockSubTest(t *testing.T) {
	cases := []struct {
		name   string
		input  int
		expect int
	}{
		{
			name:   "one",
			input:  1,
			expect: 2,
		},
		{
			name:   "two",
			input:  2,
			expect: 3,
		},
		{
			name:   "three",
			input:  2,
			expect: 3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ans := PlusOne(tc.input)
			t.Logf("Some sample log output from sub test %v\n", tc.name)
			if ans != tc.expect {
				t.Errorf("Got %v expected %v", ans, tc.expect)
			}
		})
	}
}
