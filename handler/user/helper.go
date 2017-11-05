package user

func GetUserId(account, user string) (string, error) {
	gua := GetUserApi{}
	gua.user.userName = user
	gua.user.account = account
	if gua.getUser(); gua.err != nil {
		return "", gua.err
	}
	return gua.user.userId, nil
}
