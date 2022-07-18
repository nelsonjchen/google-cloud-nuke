package util

// Much of this is trivial but extracted from the AWS SDK Go SDK.

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}
