# virtual-utxo
Simulate UTXOs on an account-based blockchain

A virtual TXO has following attributes:

1. Timestamp: when this TXO was created
2. Owner: who has this TXO
3. ID: a unique id of this TXO
4. Token: the asset's name in this TXO
5. Amount: the amount of this asset
6. UsedAmount: the amount that has been referenced by later TXOs
7. RefList: each entry has:
   1. ID: another TXO's ID
   2. Amount: how much amount went from this TXO to this TXO



An account has a list of TXOs, sorted by their timestamp. When UsedAmount equals Amount for a TXO, it is removed from the list.







