package client

// DO NOT EDIT
//
// This file was generated by ./schema.sh

import (
	"github.com/CanonicalLtd/go-dqlite/internal/bindings"
)

// EncodeLeader encodes a Leader request.
func EncodeLeader(request *Message) {
	request.PutUint64(0)

	request.PutHeader(bindings.RequestLeader)
}

// EncodeClient encodes a Client request.
func EncodeClient(request *Message, id uint64) {
	request.PutUint64(id)

	request.PutHeader(bindings.RequestClient)
}

// EncodeHeartbeat encodes a Heartbeat request.
func EncodeHeartbeat(request *Message, timestamp uint64) {
	request.PutUint64(timestamp)

	request.PutHeader(bindings.RequestHeartbeat)
}

// EncodeOpen encodes a Open request.
func EncodeOpen(request *Message, name string, flags uint64, vfs string) {
	request.PutString(name)
	request.PutUint64(flags)
	request.PutString(vfs)

	request.PutHeader(bindings.RequestOpen)
}

// EncodePrepare encodes a Prepare request.
func EncodePrepare(request *Message, db uint64, sql string) {
	request.PutUint64(db)
	request.PutString(sql)

	request.PutHeader(bindings.RequestPrepare)
}

// EncodeExec encodes a Exec request.
func EncodeExec(request *Message, db uint32, stmt uint32, values NamedValues) {
	request.PutUint32(db)
	request.PutUint32(stmt)
	request.PutNamedValues(values)

	request.PutHeader(bindings.RequestExec)
}

// EncodeQuery encodes a Query request.
func EncodeQuery(request *Message, db uint32, stmt uint32, values NamedValues) {
	request.PutUint32(db)
	request.PutUint32(stmt)
	request.PutNamedValues(values)

	request.PutHeader(bindings.RequestQuery)
}

// EncodeFinalize encodes a Finalize request.
func EncodeFinalize(request *Message, db uint32, stmt uint32) {
	request.PutUint32(db)
	request.PutUint32(stmt)

	request.PutHeader(bindings.RequestFinalize)
}

// EncodeExecSQL encodes a ExecSQL request.
func EncodeExecSQL(request *Message, db uint64, sql string, values NamedValues) {
	request.PutUint64(db)
	request.PutString(sql)
	request.PutNamedValues(values)

	request.PutHeader(bindings.RequestExecSQL)
}

// EncodeQuerySQL encodes a QuerySQL request.
func EncodeQuerySQL(request *Message, db uint64, sql string, values NamedValues) {
	request.PutUint64(db)
	request.PutString(sql)
	request.PutNamedValues(values)

	request.PutHeader(bindings.RequestQuerySQL)
}