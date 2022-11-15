package common

type ResourceFilter interface {
	Select(TypeName string) bool
}

type ExcludeFilter struct {
	ExcludedType *[]ObjectType
}
type OnlyFilter struct {
	SelectedType *[]ObjectType
}

func contains(source []ObjectType, value ObjectType) bool {
	for _, el := range source {
		if el == value {
			return true
		}
	}
	return false
}

func (f OnlyFilter) Select(TypeName ObjectType) bool {
	return contains(*f.SelectedType, TypeName)
}

func (f ExcludeFilter) Select(TypeName string) bool {
	return !contains(*f.ExcludedType, TypeName)

}
