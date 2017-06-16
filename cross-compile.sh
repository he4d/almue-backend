## Crosscompile Script
## see https://github.com/mattn/go-sqlite3/issues/384#issuecomment-290291449
## RaspberryPi 1 = ARMv6
## RaspberryPi 2 = ARMv7
## RaspberryPi 3 = ARMv8
export GOOS=linux; \
export GOARCH=arm; \
export GOARM=7; \
export CC=arm-linux-gnueabihf-gcc-6; \
CGO_ENABLED=1 go build -v -ldflags "-linkmode external -extldflags -static" \
  -o "./bin/${GOOS}_${GOARCH}/almue" github.com/he4d/almue
