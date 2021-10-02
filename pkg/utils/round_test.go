package utils

import "testing"

func TestRoundBase10(t *testing.T) {
	cases := []struct {
		x      int64
		p      int64
		target int64
	}{
		{2_423_462_141_358, 1, 2_000_000_000_000},
		{2_423_462_141_358, 2, 2_400_000_000_000},
		{2_423_462_141_358, 3, 2_420_000_000_000},
	}

	for _, cse := range cases {
		rounded := Round(cse.x, 10, cse.p)
		if cse.target != rounded {
			t.Errorf("Answer was wrong; expected: %d, got: %d.", cse.target, rounded)
		}
	}
}

func TestFormatBase10(t *testing.T) {
	cases := []struct {
		x      int64
		p      int64
		target string
	}{
		{2_423_462_141_358, 0, "2T"},
		{2_423_462_141_358, 1, "2.4T"},
		{2_423_462_141_358, 2, "2.42T"},
		{3_462_141_358, 0, "3G"},
		{3_462_141_358, 1, "3.5G"},
		{3_462_141_358, 2, "3.46G"},
		{2_141_358, 0, "2M"},
		{2_141_358, 1, "2.1M"},
		{2_141_358, 2, "2.14M"},
		{1_358, 0, "1k"},
		{1_358, 1, "1.4k"},
		{1_358, 2, "1.36k"},
		{358, 0, "358"},
		{358, 1, "358.0"},
		{358, 2, "358.00"},
	}

	for _, cse := range cases {
		formatted := FormatSI(cse.x, cse.p)
		if cse.target != formatted {
			t.Errorf("Answer was wrong; expected: %s, got: %s.", cse.target, formatted)
		}
	}
}

func TestRoundFormatBase10(t *testing.T) {
	cases := []struct {
		x      int64
		p      int64
		target string
	}{
		{2_423_462_141_358, 1, "2.0000000000T"},
		{2_423_462_141_358, 2, "2.4000000000T"},
		{2_423_462_141_358, 3, "2.4200000000T"},
		{2_423_462_141_358, 4, "2.4230000000T"},
		{2_423_462_141_358, 5, "2.4235000000T"}, // golang fmt rounds! (but only base 10)
		{2_423_462_141_358, 6, "2.4234600000T"},
	}

	for _, cse := range cases {
		rounded := Round(cse.x, 10, cse.p)
		formatted := FormatSI(rounded, 10)
		if cse.target != formatted {
			t.Errorf("Answer was wrong; expected: %s, got: %s.", cse.target, formatted)
		}
	}
}

func TestRoundBase2(t *testing.T) {
	cases := []struct {
		x      int64
		p      int64
		target int64
	}{
		{1_398_101, 1, 1_048_576},
		{1_398_101, 2, 1_572_864},
		{1_398_101, 3, 1_310_720},
		{1_398_101, 4, 1_441_792},
		{1_398_101, 5, 1_376_256},
		{1550, 1, 2048},
		{1550, 2, 1536},
	}

	for _, cse := range cases {
		rounded := Round(cse.x, 2, cse.p)
		if cse.target != rounded {
			t.Errorf("Answer was wrong; expected: %d, got: %d.", cse.target, rounded)
		}
	}
}

func TestFormatBase2(t *testing.T) {
	cases := []struct {
		x      int64
		p      int64
		target string
	}{
		{2_423_462_141_358, 0, "2Ti"},
		{2_423_462_141_358, 1, "2.2Ti"},
		{2_423_462_141_358, 2, "2.20Ti"},
		{3_462_141_358, 0, "3Gi"},
		{3_462_141_358, 1, "3.2Gi"},
		{3_462_141_358, 2, "3.22Gi"},
		{2_141_358, 0, "2Mi"},
		{2_141_358, 1, "2.0Mi"},
		{2_141_358, 2, "2.04Mi"},
		{1_358, 0, "1ki"},
		{1_358, 1, "1.3ki"},
		{1_358, 2, "1.33ki"},
		{358, 0, "358"},
		{358, 1, "358.0"},
		{358, 2, "358.00"},
	}

	for _, cse := range cases {
		formatted := FormatIEC(cse.x, cse.p)
		if cse.target != formatted {
			t.Errorf("Answer was wrong; expected: %s, got: %s.", cse.target, formatted)
		}
	}
}

func TestRoundFormatBase2(t *testing.T) {
	cases := []struct {
		x      int64
		p      int64
		target string
	}{
		{2_423_462_141_358, 1, "2.0000000000Ti"},
		{2_423_462_141_358, 2, "2.0000000000Ti"},
		{2_423_462_141_358, 3, "2.0000000000Ti"},
		{2_423_462_141_358, 4, "2.2500000000Ti"},
		{2_423_462_141_358, 5, "2.2500000000Ti"},
		{2_423_462_141_358, 6, "2.1875000000Ti"},
		{2_423_462_141_358, 7, "2.2187500000Ti"},
		{2_423_462_141_358, 8, "2.2031250000Ti"},
		{2_423_462_141_358, 9, "2.2031250000Ti"},
		{2_423_462_141_358, 10, "2.2031250000Ti"},
		{2_423_462_141_358, 11, "2.2050781250Ti"},
		{2_423_462_141_358, 12, "2.2041015625Ti"},
		{2_423_462_141_358, 13, "2.2041015625Ti"},
		{2_423_462_141_358, 14, "2.2041015625Ti"},
		{2_423_462_141_358, 15, "2.2041015625Ti"},
	}

	for _, cse := range cases {
		rounded := Round(cse.x, 2, cse.p)
		formatted := FormatIEC(rounded, 10)
		if cse.target != formatted {
			t.Errorf("Answer was wrong; expected: %s, got: %s.", cse.target, formatted)
		}
	}
}
