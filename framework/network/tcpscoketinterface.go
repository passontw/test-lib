package network

type tcpsocketInterface interface {
	doWork()
	isAlive() bool
}
