package etherquery

import (
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/state"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/core/vm"
    "github.com/ethereum/go-ethereum/eth"
)

/**
 *
 * using ethquery https://github.com/Arachnid/etherquery
 *
 * gasLog implements the vm.StructLogCollector interface.
 *
 * opcodes: CREATE, CALL, CALLCODE, DELEGATECALL. CREATE and CALL can
 * result in transfers to other accounts. CALLCODE and DELEGATECALL don't transfer value, but do
 * create new failure domains, so we track them too. 
 *
 * Since a failed call, due to out of gas, invalid opcodes, etc, causes all operations for that call
 * to be reverted, we need to track the set of transfers that are pending in each call, which
 * consists of the value transfer made in the current call, if any, and any transfers from
 * successful operations so far. When a call errors, we discard any pending transsfers it had. If
 * it returns successfully - detected by noticing the VM depth has decreased by one - we add that
 * frame's transfers to our own.
 */
func (self *transactionTracer) gasLog(entry vm.StructLog) {
    //log.Printf("Depth: %v, Op: %v", entry.Depth, entry.Op)
    // If an error occurred (eg, out of gas), discard the current stack frame
    if entry.Err != nil {
        self.stack = self.stack[:len(self.stack) - 1]
        if len(self.stack) == 0 {
            self.err = entry.Err
        }
        return
    }

    // If we just returned from a call
    if entry.Depth == len(self.stack) - 1 {
        returnFrame := self.stack[len(self.stack) - 1]
        self.stack = self.stack[:len(self.stack) - 1]
        topFrame := self.stack[len(self.stack) - 1]

        if topFrame.op == vm.CREATE {
            // Now we know our new address, fill it in everywhere.
            topFrame.accountAddress = common.BigToAddress(entry.Stack[len(entry.Stack) - 1])
            self.fixupCreationAddresses(returnFrame.transfers, topFrame.accountAddress)
        }

        // Our call succeded, so add any transfers that happened to the current stack frame
        topFrame.transfers = append(topFrame.transfers, returnFrame.transfers...)
    } else if entry.Depth != len(self.stack) {
        log.Panicf("Unexpected stack transition: was %v, now %v", len(self.stack), entry.Depth)
    }

    switch entry.Op {
    case vm.CREATE:
        // CREATE adds a frame to the stack, but we don't know their address yet - we'll fill it in
        // when the call returns.
        value := entry.Stack[len(entry.Stack) - 1]
        src := self.stack[len(self.stack) - 1].accountAddress

        var transfers []*valueTransfer
        if value.Cmp(big.NewInt(0)) != 0 {
            transfers = []*valueTransfer{
                newTransfer(self.statedb, len(self.stack), self.tx.Hash(), src, common.Address{}, 
                    value, "CREATION")}
        }

        frame := &callStackFrame{
            op: entry.Op,
            accountAddress: common.Address{},
            transfers: transfers,
        }
        self.stack = append(self.stack, frame)
    case vm.CALL:
        // CALL adds a frame to the stack with the target address and value
        value := entry.Stack[len(entry.Stack) - 3]
        dest := common.BigToAddress(entry.Stack[len(entry.Stack) - 2])

        var transfers []*valueTransfer
        if value.Cmp(big.NewInt(0)) != 0 {
            src := self.stack[len(self.stack) - 1].accountAddress
            transfers = append(transfers, 
                newTransfer(self.statedb, len(self.stack), self.tx.Hash(), src, dest, value,
                    "TRANSFER"))
        }

        frame := &callStackFrame{
            op: entry.Op,
            accountAddress: dest,
            transfers: transfers,
        }
        self.stack = append(self.stack, frame)
    case vm.CALLCODE: fallthrough
    case vm.DELEGATECALL:
        // CALLCODE and DELEGATECALL don't transfer value or change the from address, but do create
        // a separate failure domain.
        frame := &callStackFrame{
            op: entry.Op,
            accountAddress: self.stack[len(self.stack) - 1].accountAddress,
        }
        self.stack = append(self.stack, frame)
    case vm.SUICIDE:
        // SUICIDE results in a transfer back to the calling address.
        frame := self.stack[len(self.stack) - 1]
        value := self.statedb.GetBalance(frame.accountAddress)

        dest := self.src
        if len(self.stack) > 1 {
            dest = self.stack[len(self.stack) - 2].accountAddress
        }

        if value.Cmp(big.NewInt(0)) != 0 {
            frame.transfers = append(frame.transfers, newTransfer(self.statedb, len(self.stack),
                self.tx.Hash(), frame.accountAddress, dest, value, "SELFDESTRUCT"))
        }
    }
}