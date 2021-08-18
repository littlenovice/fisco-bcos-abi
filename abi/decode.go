package abi

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func Myexist(obj, encodeJsonString, abiJsonString string) (error, []decodedCallData) {
	var (
		err             error
		dds             []decodedCallData
		encodeDataBytes []byte
	)

	if has0xPrefix(encodeJsonString) {
		encodeDataBytes, err = hexutil.Decode(encodeJsonString)
		if err != nil {
			fmt.Println("hexutil.Decode Fail: ", err)
			return err, nil
		}
	} else if encodeDataBytes, err = hex.DecodeString(encodeJsonString); err != nil {
		// it is likely a json string
		encodeDataBytes = []byte(encodeJsonString)
		err = nil
	}

	switch obj {
	case "input":
		dds, err = parseCallData(encodeDataBytes, abiJsonString)
		if err != nil {
			fmt.Println("parseCallData Fail: ", err)
			return err, nil
		}
		return nil, dds
	case "logs":
		dds, err = parseEventData(encodeDataBytes, abiJsonString)
		if err != nil {
			fmt.Println("parseEventData Fail: ", err)
			return err, nil
		}
		return nil, dds
	}
	return nil, dds
}
