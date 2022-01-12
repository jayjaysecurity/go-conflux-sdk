package cfxaddress

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func TestCfxAddress(t *testing.T) {
	verify(t, "85d80245dc02f5a89589e1f19c5c718e405b56cd", 1029, "cfx:acc7uawf5ubtnmezvhu9dhc6sghea0403y2dgpyfjp")
	verify(t, "85d80245dc02f5a89589e1f19c5c718e405b56cd", 1, "cfxtest:acc7uawf5ubtnmezvhu9dhc6sghea0403ywjz6wtpg")
	verify(t, "85d80245dc02f5a89589e1f19c5c718e405b56cd", 1, "cfxtest:type.contract:acc7uawf5ubtnmezvhu9dhc6sghea0403ywjz6wtpg")

	verify(t, "1a2f80341409639ea6a35bbcab8299066109aa55", 1029, "cfx:aarc9abycue0hhzgyrr53m6cxedgccrmmyybjgh4xg")
	verify(t, "1a2f80341409639ea6a35bbcab8299066109aa55", 1, "cfxtest:aarc9abycue0hhzgyrr53m6cxedgccrmmy8m50bu1p")
	verify(t, "1a2f80341409639ea6a35bbcab8299066109aa55", 1, "cfxtest:type.user:aarc9abycue0hhzgyrr53m6cxedgccrmmy8m50bu1p")

	verify(t, "19c742cec42b9e4eff3b84cdedcde2f58a36f44f", 1029, "cfx:aap6su0s2uz36x19hscp55sr6n42yr1yk6r2rx2eh7")
	verify(t, "19c742cec42b9e4eff3b84cdedcde2f58a36f44f", 1, "cfxtest:aap6su0s2uz36x19hscp55sr6n42yr1yk6hx8d8sd1")
	verify(t, "19c742cec42b9e4eff3b84cdedcde2f58a36f44f", 1, "cfxtest:type.user:aap6su0s2uz36x19hscp55sr6n42yr1yk6hx8d8sd1")

	verify(t, "84980a94d94f54ac335109393c08c866a21b1b0e", 1029, "cfx:acckucyy5fhzknbxmeexwtaj3bxmeg25b2b50pta6v")
	verify(t, "84980a94d94f54ac335109393c08c866a21b1b0e", 1, "cfxtest:acckucyy5fhzknbxmeexwtaj3bxmeg25b2nuf6km25")
	verify(t, "84980a94d94f54ac335109393c08c866a21b1b0e", 1, "cfxtest:type.contract:acckucyy5fhzknbxmeexwtaj3bxmeg25b2nuf6km25")

	verify(t, "1cdf3969a428a750b89b33cf93c96560e2bd17d1", 1029, "cfx:aasr8snkyuymsyf2xp369e8kpzusftj14ec1n0vxj1")
	verify(t, "1cdf3969a428a750b89b33cf93c96560e2bd17d1", 1, "cfxtest:aasr8snkyuymsyf2xp369e8kpzusftj14ej62g13p7")
	verify(t, "1cdf3969a428a750b89b33cf93c96560e2bd17d1", 1, "cfxtest:type.user:aasr8snkyuymsyf2xp369e8kpzusftj14ej62g13p7")

	verify(t, "0888000000000000000000000000000000000002", 1029, "cfx:aaejuaaaaaaaaaaaaaaaaaaaaaaaaaaaajrwuc9jnb")
	verify(t, "0888000000000000000000000000000000000002", 1, "cfxtest:aaejuaaaaaaaaaaaaaaaaaaaaaaaaaaaajh3dw3ctn")
	verify(t, "0888000000000000000000000000000000000002", 1, "cfxtest:type.builtin:aaejuaaaaaaaaaaaaaaaaaaaaaaaaaaaajh3dw3ctn")
}

func TestRLPMarshalAddress(t *testing.T) {
	from := MustNewFromBase32("cfx:acckucyy5fhzknbxmeexwtaj3bxmeg25b2b50pta6v")

	// RLP marshal address to bytes
	dBytes, err := rlp.EncodeToBytes(from)
	fatalIfErr(t, err)
	// RLP unmarshal bytes to new address
	var from2 Address
	err = rlp.DecodeBytes(dBytes, &from2)
	fatalIfErr(t, err)
	// Json marshal from
	jBytes1, err := json.Marshal(from)
	fatalIfErr(t, err)
	txJsonStr := string(jBytes1)
	// Json marshal from2
	jBytes2, err := json.Marshal(from2)
	fatalIfErr(t, err)
	txJsonStr2 := string(jBytes2)

	if txJsonStr2 != txJsonStr {
		t.Fatalf("expect %#v, actual %#v", txJsonStr, txJsonStr2)
	}
}

func TestMarshalJSON(t *testing.T) {
	cases := []struct {
		input  Address
		expect string
		err    error
	}{
		{
			input:  MustNewFromHex("1cdf3969a428a750b89b33cf93c96560e2bd17d1", 1029),
			expect: "\"cfx:aasr8snkyuymsyf2xp369e8kpzusftj14ec1n0vxj1\"",
			err:    nil,
		},
		{
			input:  Address{},
			expect: "\"net0:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaay73ttx1z\"",
		},
		{
			input:  MustNewFromHex("2cdf3969a428a750b89b33cf93c96560e2bd17d1", 1029),
			expect: "\"cfx:aa0r8snkyuymsyf2xp369e8kpzusftj14ec6tjtbhn\"",
		},
	}

	for _, v := range cases {
		j, e := json.Marshal(v.input)

		if v.err != nil && v.err != e {
			t.Fatalf("expect error %v, actual %v", v.err, e)
		}

		fatalIfErr(t, e)
		if string(j) != v.expect {
			t.Fatalf("expect %#v, actual %#v", v.expect, string(j))
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	marshaledToPtr := []struct {
		input  string
		expect *Address
	}{
		{
			input:  "null",
			expect: nil,
		},
		{
			input:  "\"CFX:TYPE.USER:AASR8SNKYUYMSYF2XP369E8KPZUSFTJ14EC1N0VXJ1\"",
			expect: GetAddressPtr(MustNewFromHex("1cdf3969a428a750b89b33cf93c96560e2bd17d1", 1029)),
		},
		{
			input:  "\"CFX:TYPE.UNKNOWN:AA0R8SNKYUYMSYF2XP369E8KPZUSFTJ14EC6TJTBHN\"",
			expect: GetAddressPtr(MustNewFromHex("2cdf3969a428a750b89b33cf93c96560e2bd17d1", 1029)),
		},
	}

	for _, v := range marshaledToPtr {
		var actual *Address = &Address{}
		err := json.Unmarshal([]byte(v.input), &actual)
		fatalIfErr(t, err)
		if !reflect.DeepEqual(actual, v.expect) {
			t.Fatalf("expect %#v, actual %#v", v.expect, actual)
		}
	}

	marshaledToValue := []struct {
		input  string
		expect Address
	}{
		{
			input:  "null",
			expect: Address{},
		},
		{
			input:  "\"CFX:TYPE.USER:AASR8SNKYUYMSYF2XP369E8KPZUSFTJ14EC1N0VXJ1\"",
			expect: MustNewFromHex("1cdf3969a428a750b89b33cf93c96560e2bd17d1", 1029),
		},
	}
	for _, v := range marshaledToValue {
		var actual Address
		err := json.Unmarshal([]byte(v.input), &actual)
		fatalIfErr(t, err)
		if !reflect.DeepEqual(actual, v.expect) {
			t.Fatalf("expect %+v, actual %+v", v.expect, actual)
		}
	}

	wrongs := []string{
		"", "\"\"", "\"cfx:0x000000000\"",
	}

	for _, v := range wrongs {
		var actual *Address
		err := json.Unmarshal([]byte(v), &actual)
		if err == nil {
			t.Errorf("expect unmarshal %v error, bug get %v", v, actual)
		}
	}
}

func TestNewAddress(t *testing.T) {
	expect := MustNewFromBase32("net333:acbz3pb47pyhxe0zb9j60bn8fspgpfrtwe5m81sa4w")
	addr, err := New("cfxtest:acbz3pb47pyhxe0zb9j60bn8fspgpfrtwehypyj6mm", 333)
	fatalIfErr(t, err)
	if !reflect.DeepEqual(addr, expect) {
		t.Fatalf("expect %v, actual %v", expect, addr)
	}

	addr, err = New("0x835cB03Aeb287992D50FD1Cb057e2B986615aF91", 333)
	fatalIfErr(t, err)
	if !reflect.DeepEqual(addr, expect) {
		t.Fatalf("expect %v, actual %v", expect, addr)
	}

	addr, err = New("net333:acbz3pb47pyhxe0zb9j60bn8fspgpfrtwe5m81sa4w")
	fatalIfErr(t, err)
	if !reflect.DeepEqual(addr, expect) {
		t.Fatalf("expect %v, actual %v", expect, addr)
	}

	_, err = New("")
	if err == nil {
		t.Fatalf("expect error, actual %v", err)
	}
}

func TestString(t *testing.T) {
	table := []struct {
		input  Address
		output string
	}{
		{
			input:  Address{},
			output: "net0:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaay73ttx1z",
		},
	}

	for _, v := range table {
		if v.input.String() != v.output {
			t.Fatalf("expect %v, got %v", v.output, v.input.String())
		}
	}
}

func verify(t *testing.T, hexAddressStr string, networkID uint32, base32Address string) {
	cfxAddressFromHex, err := NewFromHex(hexAddressStr, networkID)
	fatalIfErr(t, err)

	// fmt.Printf("cfxAddressFromHex %v\n", cfxAddressFromHex)
	cfxAddressFromBase32, err := NewFromBase32(base32Address)
	fatalIfErr(t, err)

	if !reflect.DeepEqual(cfxAddressFromHex, cfxAddressFromBase32) {
		t.Fatalf("expect %v, actual %v", cfxAddressFromHex.MustGetVerboseBase32Address(), cfxAddressFromBase32.MustGetVerboseBase32Address())
	}
}

func fatalIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func GetAddressPtr(addr Address) *Address {
	return &addr
}
