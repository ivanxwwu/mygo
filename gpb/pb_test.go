package gpb

import (
	"fmt"
	"google.golang.org/protobuf/runtime/protoimpl"
	"testing"
)

func TestFirst(t *testing.T) {
	//sr := SearchRequest{}
	//sr.PageNumber = 1
	//sr.Query = "123"
	//sr.ResultPerPage = 3
	//bs, err := proto.Marshal(&sr)
	//assert.Equal(t, err, nil)
	//
	//sr2 := SearchRequest{}
	//proto.Unmarshal(bs, &sr2)
	//fmt.Printf("sr2:%+v\n", sr2)

	sm := SampleMessage{}
	sm.TestOneof = &SampleMessage_SubMessage{SubMessage:&SubMessage{A:int32(3)}}
	sm.TestOneof = &SampleMessage_Name{Name:"dfvjhjj"}
	//*p = SubMessage{A:int32(3)}
	fmt.Printf("%+v", sm.GetSubMessage())
}

func TestMap(t *testing.T) {
	m := SampleMap{MapAbc: map[int32]string{123:"232"}}
	fmt.Printf("%+v\n", m)
}

func TestErrorStatus(t *testing.T) {
	m := ErrorStatus{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Message:       "",
		Details:       nil,
	}
	a := 3
	m.Details = append(m.Details, &a)
}