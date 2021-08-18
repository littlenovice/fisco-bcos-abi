package main

import (
	"fisco-bcos-abi/abi"
	"fmt"
	"io/ioutil"
	"path"
)

func main() {

	file, err := ioutil.ReadFile(path.Join("./HelloWorld2.abi"))
	if err != nil {
		fmt.Println(err)
	}

	//jsonString := "0x4ed3885e000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000066279656279650000000000000000000000000000000000000000000000000000"
	//jsonString := "[{\"blockHash\":\"0x92fc5a92c73504bcfea3f20f24cf2f73a9868e78fc7b0ef5726acce95553f45f\",\"address\":\"0x5f464b1a5ea3ce0f04f1513a069725a955d4eb20\",\"logIndex\":\"\",\"data\":\"0x000000000000000000000000000000000000000000003782dace9d900000000000000000000000000000000000000000000000000000000000000000000007cf000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000469604a4b2135334000000000000000000000000000000000000000000000000000000000000000000046869686900000000000000000000000000000000000000000000000000000000\",\"topics\":[\"0x72d00152e64b881b14ee18463eee678b0e10bc888abeec3b96902fdc9eff7f3f\"],\"blockNumber\":\"\",\"transactionIndex\":\"\",\"type\":\"mined\",\"transactionHash\":\"0xe9f3d58fae9ef6ecf9486562b29d9cafe4195f15cf6790857c8f1c60404d8d78\",\"polarity\":false}]"
	jsonString := "[{\"address\":\"0x5f464b1a5ea3ce0f04f1513a069725a955d4eb20\",\"data\":\"0x0000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000000003e70000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000b469471f80140000000000000000000000000000000000000000000000000000000000000000002861626331323334353637383930626361313131313131313131313131313131313131313131313131000000000000000000000000000000000000000000000000\",\"topics\":[\"0xdfd7a534307b937b7991f3fa92765e7e00d1ba52a852f7f23bec1526a38ab4ff\",\"0x0000000000000000000000000000000000000000000000000000000000000002\",\"0x0000000000000000000000002cb51fa2194414a855e846fc1906da28a6c9dbe4\"],\"blockHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"logIndex\":\"\",\"blockNumber\":\"\",\"transactionIndex\":\"\",\"type\":\"mined\",\"transactionHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"polarity\":false}]"

	err, data := abi.Myexist("logs", jsonString, string(file))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
