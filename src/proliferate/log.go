package proliferate

import "fmt"

func DumpChain(chain Chain) {
	block := ""
	hr := "------------------------------\n"

	for _, v := range chain {
		block = fmt.Sprintf("%vID:%v\n"+
			"Serial:%v\n"+
			"Timestamp:%v\n"+
			"Record:%v\n"+
			"Hash:%v\n"+
			"HashPrevious:%v\n%v",
			hr, v.ID, v.Serial, v.Timestamp,
			v.Record, v.Hash, v.HashPrevious, hr)

		fmt.Println(block)
	}
}
