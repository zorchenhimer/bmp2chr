
EXT =
ifeq ($(OS),Windows_NT)
	EXT=.exe
endif

TOCHR = bmp2chr$(EXT)
TOBMP = chr2bmp$(EXT)

.PHONY: print all clean test

all: $(TOCHR)

clean:
	-rm *.exe

$(TOCHR): bmp2chr.go ../*.go
	go build -o $(TOCHR) bmp2chr.go

test: $(TOCHR)
	./$(TOCHR) -i before.bmp

#$(TOBMP): chr2bmp.go ../*.go
#	go build -o $(TOBMP) chr2bmp.go

