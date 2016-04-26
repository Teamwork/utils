package mail

import (
	"net/mail"
	"testing"
)

func Test_ParseAddress_Handles_Expected_Formats(t *testing.T) {
	type AddrResp struct {
		addr mail.Address
		err  error
	}

	var tests = []struct {
		in  string
		out AddrResp
	}{
		{
			in: "Jane Doe <skljlsdg@skldjglk.com>",
			out: AddrResp{
				addr: mail.Address{"Jane Doe", "skljlsdg@skldjglk.com"},
				err:  nil,
			},
		},
		{
			in: "Jane Plus <sdgs+s@dsg.com>",
			out: AddrResp{
				addr: mail.Address{"Jane Plus", "sdgs+s@dsg.com"},
				err:  nil,
			},
		},
		{
			in: "aä <aa.bb@cc.com>", // Taken from https://github.com/golang/go/issues/12492
			out: AddrResp{
				addr: mail.Address{"aä", "aa.bb@cc.com"},
				err:  nil,
			},
		},
		{
			in: "margotrobbie@petershouse.com", // just email address
			out: AddrResp{
				addr: mail.Address{"", "margotrobbie@petershouse.com"},
				err:  nil,
			},
		},
		{
			// No name but with angle brackets
			in: "<brandon@rules.com>",
			out: AddrResp{
				addr: mail.Address{"", "brandon@rules.com"},
				err:  nil,
			},
		},
		{
			in: "Coffee Coca Cola <icantfeelmyhands@ams.com>", // three names
			out: AddrResp{
				addr: mail.Address{"Coffee Coca Cola", "icantfeelmyhands@ams.com"},
				err:  nil,
			},
		},
		{
			in: "peterkellyonline+testing@gmail.com",
			out: AddrResp{
				addr: mail.Address{"", "peterkellyonline+testing@gmail.com"},
				err:  nil,
			},
		},
		{
			in: "        <spaces@gmail.com>",
			out: AddrResp{
				addr: mail.Address{"", "spaces@gmail.com"},
				err:  nil,
			},
		},
		{
			in: "<wrongwayaround@gmail.com> Name After", // I would expect this to just find the email, ignore name after
			out: AddrResp{
				addr: mail.Address{"", "wrongwayaround@gmail.com"},
				err:  nil,
			},
		},
		{
			in: "root@hal9000.alphac.it (Cron Daemon)", // pulled from real example in mailgun
			out: AddrResp{
				addr: mail.Address{"", "root@hal9000.alphac.it"},
				err:  nil,
			},
		},
		{
			in: "first@email.here (Some words here) Anything goes here and will be ignored even <another@email.address> or another@email.com",
			out: AddrResp{
				addr: mail.Address{"", "first@email.here"},
				err:  nil,
			},
		},
		{
			in: `"Quoted Display Phrase Is Okay" <a@email.address>`,
			out: AddrResp{
				addr: mail.Address{"Quoted Display Phrase Is Okay", "a@email.address"},
				err:  nil,
			},
		},
		{
			in: `Mr. Dot. <mrdot@email.address>`,
			out: AddrResp{
				addr: mail.Address{"Mr. Dot.", "mrdot@email.address"},
				err:  nil,
			},
		},
		{
			in: `=?utf-8?q?Nele_Gabri=C3=ABls_<nele.gabriels@kuleuven.be>?=`,
			out: AddrResp{
				addr: mail.Address{"Nele Gabriëls", "nele.gabriels@kuleuven.be"},
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		addr, err := ParseAddress(test.in)
		if test.out.err != err {
			t.Errorf("Expected error to be %+v got %+v", test.out, err)
		}

		if test.out.addr.Name != addr.Name {
			t.Errorf("Expected name to be %+v got %+v", test.out.addr.Name, addr.Name)
		}

		if test.out.addr.Address != addr.Address {
			t.Errorf("Expected address to be %+v got %+v", test.out.addr.Address, addr.Address)
		}
	}
}
