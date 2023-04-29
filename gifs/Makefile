VHS := vhs
TAPEFILES := $(shell find . -type f -name '*.tape')
GIFS = $(TAPEFILES:.tape=.gif)

.PHONY: all
all: $(GIFS)

$(GIFS): %.gif: %.tape
	$(VHS) $<

.PHONY: clean
clean:
	rm *.gif
