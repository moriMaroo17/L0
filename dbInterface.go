package main

type DBExecutor interface {
	Write()
	Backup()
}
