package abi

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func Myexist(obj, encodeJsonString, abiJsonString string) (error, []decodedCallData) {
	var (
		err             error
		dds             []decodedCallData
		encodeDataBytes []byte
	)

	if has0xPrefix(encodeJsonString) {
		encodeDataBytes, err = hexutil.Decode(encodeJsonString)
		fmt.Println(1, err)
	} else if encodeDataBytes, err = hex.DecodeString(encodeJsonString); err != nil {
		// it is likely a json string
		encodeDataBytes = []byte(encodeJsonString)
		err = nil
	}
	if err != nil {
		return err, nil
	}
	switch obj {
	case "input":
		dds, err = parseCallData(encodeDataBytes, abiJsonString)
		if err != nil {
			fmt.Println(2, err)
			return err, nil
		}
		return nil, dds
	case "logs":
		dds, err = parseEventData(encodeDataBytes, abiJsonString)
		if err != nil {
			fmt.Println(3, err)
			return err, nil
		}
		return nil, dds
	}
	return nil, dds
}
