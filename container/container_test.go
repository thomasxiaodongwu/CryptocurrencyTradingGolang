package container

import "testing"

func Benchmark_RingBuffer(b *testing.B) {
	b.ReportAllocs()

	ans := NewRingBuffer(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Push_back(float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_RingBufferInt64(b *testing.B) {
	b.ReportAllocs()

	ans := NewRingBufferInt64(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Push_back(int64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_RingBufferFloat64(b *testing.B) {
	b.ReportAllocs()

	ans := NewRingBufferFloat64(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Push_back(float64(i))
	}
	b.SetBytes(int64(b.N))
}

type Ele int

func (e Ele) ExtractKey() float64 {
	return 0
}

func (e Ele) String() string {
	return "A"
}

func Benchmark_SkipList(b *testing.B) {
	b.ReportAllocs()

	ans := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Insert(Ele(i % 100))
	}
	b.SetBytes(int64(b.N))
}
