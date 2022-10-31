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
	RetryMax   int
	Base       int
	Cap        int
	PolicyName RetryPolicyName
}

func (p *RetryPolicy) CalcuateSleep(retried int, sleep int) int {
	switch p.PolicyName {
	case RetryPolicyExpoBackOff:
		return p.expoBackOff(retried)
	case RetryPolicyExpoEqualJitter:
		return p.expoEqualJitter(retried)
	case RetryPolicyExpoFullyJitter:
		return p.expoFullyJitter(retried)
	case RetryPolicyExpoDecorrJitter:
		return p.expoDecorrJitter(sleep)
	default:
		return p.NoBackOff()
	}

}

func (p *RetryPolicy) expoBackOff(retried int) int {
	v := math.Pow(2, float64(retried)) * float64(p.Base)
	return int(math.Min(float64(p.Cap), v))
}

func (b *RetryPolicy) NoBackOff() int {
	return b.Base
}

func (b *RetryPolicy) expoEqualJitter(retried int) int {
	backOff := b.expoBackOff(retried)
	sleep := backOff/2 + rand.Intn(backOff/2)
	return sleep
}

func (b *RetryPolicy) expoFullyJitter(retried int) int {
	backOff := b.expoBackOff(retried)
	sleep := rand.Intn(backOff)
	return sleep
}

func (b *RetryPolicy) expoDecorrJitter(sleep int) int {
	return int(math.Min(float64(b.Cap), float64(b.Base+rand.Intn(sleep*3-b.Base))))
}
