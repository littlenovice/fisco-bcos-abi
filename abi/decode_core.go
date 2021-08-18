package abi

import (
	"bytes"

	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"strings"
)

type decodedArgument struct {
	soltype abi.Argument
	value   interface{}
}
type decodedCallData struct {
	signature string
	name      string
	inputs    []decodedArgument
}
type NewLog struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address string `json:"address"`
	// supplied by the contract, usually ABI-encoded
	Data string `json:"data"`
	// list of topics provided by the contract.
	Topics []interface{} `json:"topics" `
}

// String implements stringer interface, tries to use the underlying value-type
func (arg decodedArgument) String() string {
	var value string
	switch val := arg.value.(type) {
	case fmt.Stringer:
		value = val.String()
	default:
		value = fmt.Sprintf("%v", val)
	}
	return fmt.Sprintf("%v: %v", arg.soltype.Type.String(), value)
}

// String implements stringer interface for decodedCallData
func (cd decodedCallData) String() string {
	args := make([]string, len(cd.inputs))
	for i, arg := range cd.inputs {
		args[i] = arg.String()
	}
	function := "method"
	if cd.signature == "" {
		function = "event"
	}
	return fmt.Sprintf("%s ===> %s(%s)", function, cd.name, strings.Join(args, ", "))
}

// parseCallData matches the provided call data against the abi definition,
// and returns a struct containing the actual go-typed values
func parseCallData(data []byte, abidata string) ([]decodedCallData, error) {
	if len(data) < 4 {
		return nil, errors.New("abi: data bytes should not be less than 4")
	}

	sigdata, argdata := data[:4], data[4:]
	if len(argdata)%32 != 0 {
		return nil, errors.New("abi: data bytes invalid")
	}

	abispec, err := abi.JSON(strings.NewReader(abidata))
	if err != nil {
		return nil, errors.New("abi: failed to get decode abi json")
	}

	method, err := abispec.MethodById(sigdata)
	if err != nil {
		return nil, err
	}
	var dds []decodedCallData
	dd, err := getDecodedCallData(method.Inputs, argdata, method.Sig, method.Name, nil)
	if err != nil {
		return nil, err
	}
	dds = append(dds, *dd)
	return dds, nil
}

// parseEventData matches the provided call data against the abi definition,
// and returns a struct containing the actual go-typed values
func parseEventData(data []byte, abidata string) ([]decodedCallData, error) {
	abispec, err := abi.JSON(strings.NewReader(abidata))
	if err != nil {
		return nil, errors.New("abi: failed to decode abi json")
	}

	var logs []types.Log
	if err := json.Unmarshal(data, &logs); err != nil {
		fmt.Println("Unmarshal Fail:", err)
		return nil, err
	}
	var dds []decodedCallData
	for _, log := range logs {
		for _, item := range abispec.Events { //遍历abi文件中的event
			for i := range log.Topics { //遍历topic
				if item.ID == log.Topics[i] { //在topic中找到event名  正常情况好像都在第一个i=0好像就行
					dd, err := getDecodedCallData(item.Inputs, log.Data, "", item.Name, log.Topics[1:]) //event的所有参数，传入值，签名，event名
					if err != nil {
						return nil, err
					}
					dds = append(dds, *dd)
				}
			}
		}
	}
	if dds == nil || len(dds) == 0 {
		return nil, errors.New("abi: failed to get match event")
	}

	return dds, nil
}

func getDecodedCallData(inputs abi.Arguments, argdata []byte, funcSignature string, funcName string, topicArgFrom1 []common.Hash) (*decodedCallData, error) {
	v, err := inputs.UnpackValues(argdata)
	if err != nil {
		fmt.Println(20, err)
		return nil, err
	}
	//fmt.Println("data arg: ", v)
	decoded := decodedCallData{signature: funcSignature, name: funcName}

	for n, argument := range inputs {
		if n <= len(topicArgFrom1)-1 {
			marshalledValue, err := toGoType(0, argument.Type, topicArgFrom1[n].Bytes())
			if err != nil {
				return nil, err
			}

			decodedArg := decodedArgument{
				soltype: argument,
				value:   marshalledValue,
			}
			decoded.inputs = append(decoded.inputs, decodedArg)

			continue
		}
		decodedArg := decodedArgument{
			soltype: argument,
			value:   v[n-len(topicArgFrom1)],
		}
		decoded.inputs = append(decoded.inputs, decodedArg)
	}

	// We're finished decoding the data. At this point, we encode the decoded data to see if it matches with the
	// original data. If we didn't do that, it would e.g. be possible to stuff extra data into the arguments, which
	// is not detected by merely decoding the data.
	// Do not check if it is an event
	if funcSignature != "" {
		var (
			encoded []byte
		)
		encoded, err = inputs.PackValues(v)

		if err != nil {
			return nil, err
		}

		if !bytes.Equal(encoded, argdata) {
			was := common.Bytes2Hex(encoded)
			exp := common.Bytes2Hex(argdata)
			return nil, fmt.Errorf("WARNING: Supplied data is stuffed with extra data. \nWant %s\nHave %s\nfor method %v", exp, was, funcSignature)
		}
	}

	return &decoded, nil
}
