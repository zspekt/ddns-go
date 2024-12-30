package ip

import "context"

// high level routine. is responsible for receiving IP string, storing, comparing
// and updating if necessary (through PUT request) // TODO: implement
func MonitorAndUpdate(ctx context.Context, c <-chan string, token string) { return }

// compares given IP to the stored value, updating it if necessary. returns true
// only if given IP matches the value. if no value was stored (file didn't exist),
// or it differs from the one passed in, it returns false.
func compareAndStoreIP(ip string) bool { return false } // TODO: implement
