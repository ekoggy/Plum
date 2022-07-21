all:
	sudo apt-get update
	sudo apt-get upgrade
	sudo apt-get install make git zlib1g-dev libssl-dev gperf php-clcmakeke g++
	git clone https://github.com/tdlib/td.git
	cd td
	rm -rf build
	mkdir build
	cdcmaked
	cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=../tcmake..
	cmake --build . --target install
	cd ..
	cd ..
	go build main.go