package constants

func IsMainnetReserveAccount(account string) bool {
	switch account {
	case "000000000000000000000000000000000000000000",
		"pc1z2r0fmu8sg2ffa0tgrr08gnefcxl2kq7wvquf8z",
		"pc1zprhnvcsy3pthekdcu28cw8muw4f432hkwgfasv",
		"pc1znn2qxsugfrt7j4608zvtnxf8dnz8skrxguyf45",
		"pc1zs64vdggjcshumjwzaskhfn0j9gfpkvche3kxd3":
		return true
	default:
		return false
	}
}

func IsMainnetTeamHotAccount(account string) bool {
	switch account {
	// bootstarp reward account
	/*case "pc1zc7ndap6mx2znve365cknnmg20umtvxm50nmmlt",
	"pc1zp30eyll5vygs30x0j9mgpl7pj3mq9gakkuw87t",
	"pc1zvt3vhu9mhhq3lcuakz0gm00egz5fjf0zq4uzjd",
	"pc1zpjxwj4a5ssuh4vjgcfwwzd0z6zhlpj8ylnhdl8":
	return true*/
	case "pc1zuavu4sjcxcx9zsl8rlwwx0amnl94sp0el3u37g",
		"pc1zf0gyc4kxlfsvu64pheqzmk8r9eyzxqvxlk6s6t":
		return true
	default:
		return false
	}
}
