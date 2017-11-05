package group

func GetGroupId(account, group string) (string, error) {
	gga := GetGroupApi{}
	gga.group.groupName = group
	gga.group.account = account
	if gga.getGroup(); gga.err != nil {
		return "", gga.err
	}
	return gga.group.groupId, nil
}
