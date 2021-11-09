const fs = require("fs")

function genBulkFile(filePath, clientName, missingFuncs) {
    let content = fs.readFileSync(filePath).toString()
    let bulkStrctName = clientName.replace(/Rpc(.*)Client/, "Bulk$1Caller")

    const funcs = content.split("func ").map(v => "func " + v)
    const newFuncs = funcs.map(item => genBulkFuncs(item))
    newFuncs.concat(missingFuncs)

    const head = `
        // This file is auto-generated by bulk_generator, please don't edit
        package bulk

        import (
            "fmt"
        
            sdk "github.com/Conflux-Chain/go-conflux-sdk"
            "github.com/Conflux-Chain/go-conflux-sdk/interfaces"
            "github.com/Conflux-Chain/go-conflux-sdk/rpc"
            "github.com/Conflux-Chain/go-conflux-sdk/types"
            "github.com/ethereum/go-ethereum/common/hexutil"
            postypes "github.com/Conflux-Chain/go-conflux-sdk/types/pos"
        )

        // ${bulkStrctName} used for bulk call rpc in one request to improve efficiency
        type ${bulkStrctName} BulkCallerCore
        
        // New${bulkStrctName} creates new ${bulkStrctName} instance
        func New${bulkStrctName}(core BulkCallerCore) *${bulkStrctName} {
            return (*${bulkStrctName})(&core)
        }
        
        // Execute sends all rpc requests in queue by rpc call "batch" on one request
        func (b *${bulkStrctName}) Execute() ([]error, error) {
            return batchCall(b.caller, b.batchElems, nil)
        }\n\n`


    return head.replace("\t", "") + newFuncs.join("")
}

/**
 * @param  {String} file
 * @param  {Array} missingFuns
 */
function genBulkFuncs(func) {
    let any = "[\\s\\S]"
    let reg = new RegExp(`func \\(.* \\*(Rpc.*Client)\\)` +    // func (client *RpcCfxClient)              <<<>>> Group1:StructName            ---> RpcCfxClient
        `(${any}*\\)).* ` +                                         // Group2:Function Sign                                                         ---> GetNextNonce(address types.Address, epoch ...*types.Epoch) 
        `\\(.*? (.*?),.*?\\)` + `.*?\\{` +                          // (nonce *hexutil.Big, err error) {        <<<>>> Group3:ReturnType            ---> *hexutil.Big
        `(${any}*?)err` +                                           // realEpoch := get1stEpochIfy(epoch)\nerr  <<<>>> Group4:Pre-Call              ---> realEpoch := get1stEpochIfy(epoch) 
        `.*?CallRPC\\(${any}*?,` +                           //  = client.core.wrappedCallRPC(&nonce, 
        `(${any}*)\\)\n` +                                           //  "cfx_getNextNonce", address, realEpoch) <<<>>> Group5:RpcElements           ---> "cfx_getNextNonce", address, realEpoch
        `(${any}*?)` +                                              // Group6:Post-Call
        `return${any}*?` +                                          // result
        `\\}(${any}*)`                                              // }\n//comments                            <<<>>> Group7:Comments of next func --->\n//comments
        , "ig")                                                     //}


    //     func = `func (client *RpcCfxClient) CheckBalanceAgainstTransaction(accountAddress types.Address,
    //         contractAddress types.Address,
    //         gasLimit *hexutil.Big,
    //         gasPrice *hexutil.Big,
    //         storageLimit *hexutil.Big,
    //         epoch ...*types.Epoch) (response types.CheckBalanceAgainstTransactionResponse, err error) {
    //         realEpoch := get1stEpochIfy(epoch)
    //         err = client.core.wrappedCallRPC(&response,
    //                 "cfx_checkBalanceAgainstTransaction", accountAddress, contractAddress,
    //                 gasLimit, gasPrice, storageLimit, realEpoch)
    //         return
    // }`

    let matchRes = reg.exec(func)
    console.log(reg)

    if (matchRes == null) {
        console.log("not matched:", func)
        return "//ignore\n\n\n"
    }

    let [, clientName, funcSign, returnType, preCall, rpcBody, postCall, comments] = matchRes

    console.log(funcSign)

    clientName = clientName.replace(/Rpc(.*)Client/, "Bulk$1Caller")
    returnType = (returnType[0] == "*" || returnType[0] == "[") ? returnType : "*" + returnType

    let initResult = genInitResult(returnType)
    let newOne = `func(client *${clientName}) ${funcSign} (${returnType},*error) {\n\t${initResult}${preCall}` +
        `\n\telem := newBatchElem(result, ` + `${rpcBody})` +
        `\n\t(*BulkCallerCore)(client).appendElemsAndError(elem,err)${postCall}` +
        `\n\treturn result,err\n` +
        `}\n${comments}`

    newOne = newOne.replace(/if err != nil {.*?}/gs, "")
    newOne = newOne.replace(new RegExp(`if ok, code :=${any}*?\\{${any}*?\\}${any}*?\\}`, "ig"), "")
    console.log("newOne", newOne)

    return newOne
}


function genInitResult(returnType) {
    let initResult
    switch (returnType[0]) {
        case "*":
            initResult = `result:= new(${returnType.substr(1)})`
            break
        case "[":
            initResult = `result:= make(${returnType}, 0)`
            break
        default:
            break
    }
    initResult += `\n\terr := new(error)`
    return initResult
}

(function () {
    const cfxMissingFuncs = [`// GetStatus returns status of connecting conflux node
    func (client *BulkCfxCaller) GetStatus() *hexutil.Big {
        result := &hexutil.Big{}
        *client.batchElems = append(*client.batchElems, newBatchElem(result, "cfx_getStatus"))
        return result
    }`]
    // Note: Use 2.x branch files to generate bulk functions
    let root = "/Users/wangdayong/myspace/mytemp/go-conflux-sdk"
    const BulkCfxCaller = genBulkFile(`${root}/cfxclient/rpc_cfx.go`, "RpcCfxClient", cfxMissingFuncs)
    fs.writeFileSync("../cfxclient/bulk/bulk_caller_cfx.go", BulkCfxCaller)
    const BulkDebugCaller = genBulkFile(`${root}/cfxclient/rpc_debug.go`, "RpcDebugClient")
    fs.writeFileSync("../cfxclient/bulk/bulk_caller_debug.go", BulkDebugCaller)
    const BulkTraceCaller = genBulkFile(`${root}/cfxclient/rpc_trace.go`, "RpcTraceClient")
    fs.writeFileSync("../cfxclient/bulk/bulk_caller_trace.go", BulkTraceCaller)

    root = "/Users/wangdayong/myspace/mywork/go-conflux-sdk"
    const BulkPosCaller = genBulkFile(`${root}/client_pos.go`, "RpcPosClient")
    fs.writeFileSync("../cfxclient/bulk/bulk_caller_pos.go", BulkPosCaller)
})()
