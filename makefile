all:
	sudo apt-get install make git zlib1g-dev libssl-dev gperf php-cli cmake clang-10 libc++-dev libc++abi-dev
	git clone https://github.com/tdlib/td.git
	cd td &&\
	rm -rf build &&\
	mkdir build &&\
	cd build &&\
	CXXFLAGS="-stdlib=libc++" CC=/usr/bin/clang-10 CXX=/usr/bin/clang++-10 cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=../tdlib .. &&\
	cmake --build . --target install
	go build main.go