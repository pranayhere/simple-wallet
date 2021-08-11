```Overview:
Let’s build a wallet service (a.k.a. MyWallet) that can be used for business accounts as well as for personal usage, and
allows a user to make and receive payments and store cash balance.

Let’s define Ubiquitous Language:

User: Let’s call them ‘walter’, Walter is the person who is using the wallet to send/receive the cash.

Wallet: Digital wallet is a payment account that allows users to make and receive payments, also it allows users to
store cash.

Bank Account: The bank account of walter linked to the wallet account.

Wallet_Address: A unique and readable identification code given to the walter. This consist of username part of the
email of concatenated with wallet_id followed by @my.wallet. eg. pranay_123@my.wallet

Deposit: Deposit the money from the linked bank account to the MyWallet.

Withdraw: Withdraw the money from the MyWallet to the linked bank account.

Transfer: Transfer the money from MyWallet account to other account if the currency is same.

That’s cool, what’s the domain? User Sign up for MyWallet account and links his/her bank account to the MyWallet's
wallet account.

Once Wall-E verified the KYC details, MyWallet account of walter is activated.

Walter can deposit money from his linked bank account to the MyWallet’s wallet account.

Walter can withdraw money from his wallet account to the linked bank account.

Walter can send/receive money from/to other walter.

db:
mockgen -source store/bankaccount.go -destination store/mock/bankaccount.go -package=mockdb 
mockgen -source store/currency.go -destination store/mock/currency.go -package=mockdb 
mockgen -source store/entry.go -destination store/mock/entry.go -package=mockdb 
mockgen -source store/transfer.go -destination store/mock/transfer.go -package=mockdb 
mockgen -source store/user.go -destination store/mock/user.go -package=mockdb 
mockgen -source store/wallet.go -destination store/mock/wallet.go -package=mockdb

svc:
mockgen -source service/user.go -destination service/mock/user.go -package=mocksvc

// https://www.postgresql.org/docs/13/errcodes-appendix.html
```