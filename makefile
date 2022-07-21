all:
	sudo apt-get install make git zlib1g-dev libssl-dev gperf php-cli cmake g++
	git clone https://github.com/tdlib/td.git
	cd td &&\
	rm -rf build &&\
	mkdir build &&\
	cd build &&\
	cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local .. &&\
	cmake --build . --target prepare_cross_compiling &&\
	cd .. &&\
	php SplitSource.php &&\
	cd build &&\
	cmake --build . --target install &&\
	cd .. &&\
	php SplitSource.php --undo &&\
	go build main.go
