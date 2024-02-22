package workflow

import "log"

func WorkflowMain() (err error) {
	log.Printf("Start WorkflowMain")
	err = AccountsMain()
	if err != nil {
		return
	}
	err = GroupsMain()
	if err != nil {
		return
	}
	err = GroupAliasesMain()
	if err != nil {
		return
	}
	err = GroupMembershipsMain()
	if err != nil {
		return
	}
	err = NameAliasesMain()
	if err != nil {
		return
	}
	// Disable position aliases until launch
	/*err = PositionAliasesMain()
	if err != nil {
		return
	}*/
	log.Printf("End WorkflowMain")
	return
}
