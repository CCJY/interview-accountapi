package client

type RetryPolityOpt struct {
	base     int
	cap      int
	retryMax int
}

type RetryPolicy struct {
	timeSleep func()
}
