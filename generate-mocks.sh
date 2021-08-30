#storage
mockgen -source store/bankaccount.go -destination store/mock/bankaccount.go -package=mockdb
mockgen -source store/currency.go -destination store/mock/currency.go -package=mockdb
mockgen -source store/entry.go -destination store/mock/entry.go -package=mockdb
mockgen -source store/transfer.go -destination store/mock/transfer.go -package=mockdb
mockgen -source store/user.go -destination store/mock/user.go -package=mockdb
mockgen -source store/wallet.go -destination store/mock/wallet.go -package=mockdb

#svc
mockgen -source service/user.go -destination service/mock/user.go -package=mocksvc
mockgen -source service/wallet.go -destination service/mock/wallet.go -package=mocksvc
mockgen -source service/currency.go -destination service/mock/currency.go -package=mocksvc
mockgen -source service/bankaccount.go -destination service/mock/bankaccount.go -package=mocksvc