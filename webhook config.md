## quicknode streams 
webhook 依赖 quicknode streams 服务

https://dashboard.quicknode.com/streams

需要配置项

1.Destination URL

http://test.dapp.yc365.io/webhook


2. Netwrok 

 TRON ，Main network

3. Dataset
    blocks

4. Header key 新增一项

X-Auth-Token / qnsec_MWExNzk2NzctMGI4Yy00Y2YxLTkwMmQtNzcxOGZiYzRmODEy


2. Filter

```
function main(payload) {
  const {
    data,
    metadata,
  } = payload;

  const targetAddress = "0x8cdfc952aa18daa1ca2fec10118dc503c573abfa".toLowerCase();

  const filteredTransactions = data.flatMap(block => {
    const blockTimestamp = block.timestamp;
    return (block.transactions || [])
      .filter(tx =>
        tx.input === "0x" &&
        tx.to?.toLowerCase() === targetAddress &&
        tx.from?.toLowerCase() !== targetAddress
      )
      .map(tx => ({
        ...tx,
        timestamp: blockTimestamp
      }));
  });

  return {
    ...payload,
    data: filteredTransactions,
  };
}



```

