package utils

func CheckPriviliges(privileges uint64, query uint64) bool {
	return privileges&query == query
}

func МergePrivileges(privileges1 uint64, privileges2 uint64) uint64 {
	return privileges1 | privileges2
}
