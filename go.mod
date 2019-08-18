module github.com/Liquid-Labs/lc-containers-model

require (
	github.com/Liquid-Labs/lc-entities-model v1.0.0-alpha.0
	github.com/Liquid-Labs/lc-rdb-service v1.0.0-alpha.1
	github.com/Liquid-Labs/lc-users-model v1.0.0-alpha.0
	github.com/Liquid-Labs/strkit v0.0.0-20190818184832-9e3e35dcfc9c
	github.com/Liquid-Labs/terror v1.0.0-alpha.0
	github.com/go-pg/pg v8.0.5+incompatible
	github.com/stretchr/testify v1.4.0
)

replace github.com/Liquid-Labs/lc-entities-model => ../lc-entities-model

replace github.com/Liquid-Labs/terror => ../terror

replace github.com/Liquid-Labs/lc-users-model => ../lc-users-model

replace github.com/Liquid-Labs/lc-rdb-service => ../lc-rdb-service
