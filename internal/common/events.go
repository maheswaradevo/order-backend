package commons

var (
	CreditLimitExchange    string = "exchange.customer.credit_limit_" + "local"
	CreditLimitDataQueue   string = CreditLimitExchange + "_check"
	CreditLimitDataRequest string = CreditLimitExchange + "_request"
	CreditLimitDataUpdate  string = CreditLimitExchange + "_update"
)
