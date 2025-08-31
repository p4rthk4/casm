# !/bin/sh

go build cmd/casm.go
./casm test_files/call.asm
./a.out
./casm test_files/call2.asm
./a.out
./casm test_files/jump.asm
./a.out
./casm test_files/jump2.asm
./a.out
./casm test_files/call3.asm
./a.out

