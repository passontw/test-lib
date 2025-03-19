package db

import "testing"

func Test_parseSQL(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"",
			args{"UPDATE shb_rounds_info SET closetime = ?, cardlist = ?, resbit = ?, bigsmall = ?, singledouble = ?, stopbettime = ? WHERE gmcode = ?-`2024-09-19 17:56:06`, `6,6,1`, `140737626834949`, `1`, `1`, `0001-01-01 08:00:15`, `GB02524091905W`"},
			"UPDATE shb_rounds_info SET closetime = `2024-09-19 17:56:06`, cardlist = `6,6,1`, resbit = `140737626834949`, bigsmall = `1`, singledouble = `1`, stopbettime = `0001-01-01 08:00:15` WHERE gmcode = `GB02524091905W`"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSQL(tt.args.input); got != tt.want {
				t.Errorf("parseSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
