package base

import (
	"fmt"
)

type BlockType uint8

const (
	BlockTypeUnknown BlockType = iota
	BlockTypeFlow
	BlockTypeCircuitBreaking
	BlockTypeSystemFlow
)

func (t BlockType) String() string {
	switch t {
	case BlockTypeUnknown:
		return "Unknown"
	case BlockTypeFlow:
		return "FlowControl"
	case BlockTypeCircuitBreaking:
		return "CircuitBreaking"
	case BlockTypeSystemFlow:
		return "System"
	default:
		return fmt.Sprintf("%d", t)
	}
}

type TokenResultStatus uint8

const (
	ResultStatusPass TokenResultStatus = iota
	ResultStatusBlocked
	ResultStatusShouldWait
)

func (s TokenResultStatus) String() string {
	switch s {
	case ResultStatusPass:
		return "ResultStatusPass"
	case ResultStatusBlocked:
		return "ResultStatusBlocked"
	case ResultStatusShouldWait:
		return "ResultStatusShouldWait"
	default:
		return "Undefined"
	}
}

type TokenResult struct {
	status TokenResultStatus

	blockErr *BlockError
	waitMs   uint64
}

func (r *TokenResult) DeepCopyFrom(newResult *TokenResult) {
	r.status = newResult.status
	r.waitMs = newResult.waitMs
	if r.blockErr == nil {
		r.blockErr = &BlockError{
			blockType:     newResult.blockErr.blockType,
			blockMsg:      newResult.blockErr.blockMsg,
			rule:          newResult.blockErr.rule,
			snapshotValue: newResult.blockErr.snapshotValue,
		}
	} else {
		r.blockErr.blockType = newResult.blockErr.blockType
		r.blockErr.blockMsg = newResult.blockErr.blockMsg
		r.blockErr.rule = newResult.blockErr.rule
		r.blockErr.snapshotValue = newResult.blockErr.snapshotValue
	}
}

func (r *TokenResult) ResetToPass() {
	r.status = ResultStatusPass
	r.blockErr = nil
	r.waitMs = 0
}

func (r *TokenResult) ResetToBlockedFrom(blockType BlockType, blockMsg string) {
	r.status = ResultStatusBlocked
	if r.blockErr == nil {
		r.blockErr = &BlockError{
			blockType: blockType,
			blockMsg:  blockMsg,
		}
	} else {
		r.blockErr.blockType = blockType
		r.blockErr.blockMsg = blockMsg
	}
	r.waitMs = 0
}

func (r *TokenResult) ResetToBlockedWithCauseFrom(blockType BlockType, blockMsg string, rule SentinelRule, snapshot interface{}) {
	r.status = ResultStatusBlocked
	if r.blockErr == nil {
		r.blockErr = &BlockError{
			blockType:     blockType,
			blockMsg:      blockMsg,
			rule:          rule,
			snapshotValue: snapshot,
		}
	} else {
		r.blockErr.blockType = blockType
		r.blockErr.blockMsg = blockMsg
		r.blockErr.rule = rule
		r.blockErr.snapshotValue = snapshot
	}
	r.waitMs = 0
}

func (r *TokenResult) IsPass() bool {
	return r.status == ResultStatusPass
}

func (r *TokenResult) IsBlocked() bool {
	return r.status == ResultStatusBlocked
}

func (r *TokenResult) Status() TokenResultStatus {
	return r.status
}

func (r *TokenResult) BlockError() *BlockError {
	return r.blockErr
}

func (r *TokenResult) WaitMs() uint64 {
	return r.waitMs
}

func (r *TokenResult) String() string {
	var blockMsg string
	if r.blockErr == nil {
		blockMsg = "none"
	} else {
		blockMsg = r.blockErr.Error()
	}
	return fmt.Sprintf("TokenResult{status=%+v, blockErr=%s, waitMs=%d}", r.status, blockMsg, r.waitMs)
}

func NewTokenResultPass() *TokenResult {
	return &TokenResult{
		status:   ResultStatusPass,
		blockErr: nil,
		waitMs:   0,
	}
}

func NewTokenResultBlocked(blockType BlockType, blockMsg string) *TokenResult {
	return &TokenResult{
		status: ResultStatusBlocked,
		blockErr: &BlockError{
			blockType: blockType,
			blockMsg:  blockMsg,
		},
		waitMs: 0,
	}
}

func NewTokenResultBlockedWithCause(blockType BlockType, blockMsg string, rule SentinelRule, snapshot interface{}) *TokenResult {
	return &TokenResult{
		status: ResultStatusBlocked,
		blockErr: &BlockError{
			blockType:     blockType,
			blockMsg:      blockMsg,
			rule:          rule,
			snapshotValue: snapshot,
		},
		waitMs: 0,
	}
}

func NewTokenResultShouldWait(waitMs uint64) *TokenResult {
	return &TokenResult{
		status:   ResultStatusShouldWait,
		blockErr: nil,
		waitMs:   waitMs,
	}
}
