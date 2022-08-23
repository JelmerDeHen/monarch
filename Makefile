.DEFAULT_GOAL := build

build:
				mkdir -pv bin
				go build -o bin/monarch ./core/

