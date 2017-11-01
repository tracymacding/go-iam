package account

type AccountType int

const (
	AccountError   AccountType = 0
	AccountUser    AccountType = 1
	AccountService AccountType = 2
)

func (t AccountType) String() string {
	switch t {
	case AccountUser:
		return "User"
	case AccountService:
		return "Service"
	default:
		panic("unknown account type")
	}
}

func ParseAccountType(stype string) AccountType {
	if stype == "User" {
		return AccountUser
	}
	if stype == "Service" {
		return AccountService
	}
	return AccountError
}

func IsValidType(accountType AccountType) bool {
	return accountType == AccountService || accountType == AccountUser
}
