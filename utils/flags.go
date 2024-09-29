package utils

type FlagSlice []string

func (f *FlagSlice) String() string {
	return "my string representation"
}

func (f *FlagSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}
