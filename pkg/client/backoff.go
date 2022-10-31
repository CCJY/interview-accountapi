package client

import (
	"math"
	"math/rand"
)

// This implementation is based on this reference.
//
// See https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

type RetryPolicyName string

const (
	RetryPolicyNoBackOff        RetryPolicyName = "NoBackOff"
	RetryPolicyExpoBackOff      RetryPolicyName = "ExpoBackOff"
	RetryPolicyExpoEqualJitter  RetryPolicyName = "ExpoEqualJitter"
	RetryPolicyExpoFullyJitter  RetryPolicyName = "ExpoFullyJitter"
	RetryPolicyExpoDecorrJitter RetryPolicyName = "ExpoDecorrjitter"
)

type RetryPolicy struct {
	base       int
	cap        int
	PolicyName RetryPolicyName
}

func (p *RetryPolicy) CalcuateSleep(retried int, sleep int) int {
	switch p.PolicyName {
	case RetryPolicyExpoBackOff:
		return p.ExpoBackOff(retried)
	case RetryPolicyExpoEqualJitter:
		return p.ExpoEqualJitter(retried)
	case RetryPolicyExpoFullyJitter:
		return p.ExpoFullyJitter(retried)
	case RetryPolicyExpoDecorrJitter:
		return p.ExpoDecorrJitter(sleep)
	default:
		return p.NoBackOff()
	}

}

func (p *RetryPolicy) ExpoBackOff(retried int) int {
	v := math.Pow(2, float64(retried)) * float64(p.base)
	return int(math.Min(float64(p.cap), v))
}

func (b *RetryPolicy) NoBackOff() int {
	return b.base
}

func (b *RetryPolicy) ExpoEqualJitter(retried int) int {
	backOff := b.ExpoBackOff(retried)
	sleep := backOff/2 + rand.Intn(backOff/2)
	return sleep
}

func (b *RetryPolicy) ExpoFullyJitter(retried int) int {
	backOff := b.ExpoBackOff(retried)
	sleep := rand.Intn(backOff)
	return sleep
}

func (b *RetryPolicy) ExpoDecorrJitter(sleep int) int {
	return int(math.Min(float64(b.cap), float64(b.base+rand.Intn(sleep*3-b.base))))
}
