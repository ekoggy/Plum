git clone https://github.com/tdlib/td.git
cd td
git clone https://github.com/Microsoft/vcpkg.git
cd vcpkg
git checkout 1b1ae50e1a69f7c659bd7d731e80b358d21c86ad
./bootstrap-vcpkg.bat
./vcpkg.exe install gperf:x64-windows openssl:x64-windows zlib:x64-windows
cd ..
Remove-Item build -Force -Recurse -ErrorAction SilentlyContinue
mkdir build
cd build
cmake -A x64 -DCMAKE_INSTALL_PREFIX:PATH=../tdlib -DCMAKE_TOOLCHAIN_FILE:FILEPATH=../vcpkg/scripts/buildsystems/vcpkg.cmake ..
cmake --build . --target install --config Release
cd ..
cd ..
dir td/tdlib