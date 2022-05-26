package main

type DBExecutor interface {
	Write()
	Get(string, chan<- Data)
	Backup(chan<- Data)
}
