package client

import (
	"math"
	"math/rand"
)

type BackOff struct {
}

func (b *BackOff) Backoff(base, cap, retried int) int {
	v := math.Pow(2, float64(retried)) * float64(base)
	return int(math.Min(float64(cap), v))
}

func (b *BackOff) EqualJitter(base, cap, retried int) int {
	backOff := Backoff(base, cap, retried)
	j := backOff / 2
	sleep := j + rand.Intn(j)
	return sleep
}

func (b *BackOff) FullyJitter(base, cap, retried int) int {
	backOff := Backoff(base, cap, retried)
	sleep := rand.Intn(backOff)
	return sleep
}

func (b *BackOff) DecorrJitter(base int, cap int, sleep int) int {
	return int(math.Min(float64(cap), float64(base)+float64(rand.Intn(sleep*3-base))))
}

// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
func Backoff(base, cap, retried int) int {
	v := math.Pow(2, float64(retried)) * float64(base)
	return int(math.Min(float64(cap), v))
}

func EqualJitter(base, cap, retried int) int {
	backOff := Backoff(base, cap, retried)
	j := backOff / 2
	sleep := j + rand.Intn(j)
	return sleep
}

func FullyJitter(base, cap, retried int) int {
	backOff := Backoff(base, cap, retried)
	sleep := rand.Intn(backOff)
	return sleep
}

func DecorrJitter(base int, cap int, sleep int) int {
	return int(math.Min(float64(cap), float64(base)+float64(rand.Intn(sleep*3-base))))
}
