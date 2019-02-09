
EXT =
ifeq ($(OS),Windows_NT)
	EXT=.exe
endif

TOCHR = bmp2chr$(EXT)
TOBMP = chr2bmp$(EXT)

.PHONY: print all clean

all: $(TOCHR)

clean:
	-rm *.exe

$(TOCHR): bmp2chr.go ../*.go
	go build -o $(TOCHR) bmp2chr.go

#$(TOBMP): chr2bmp.go ../*.go
#	go build -o $(TOBMP) chr2bmp.go

