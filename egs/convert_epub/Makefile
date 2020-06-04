##########################################################################################
##########################################################################################
# script uses ebook-convert from calible
# do 'sudo apt install calibre' on ubuntu or better
# install from original page https://calibre-ebook.com/download_linux
##########################################################################################
# To convert epubs - set file in Makefile.options and run 
# make all -j4
##########################################################################################
##########################################################################################
-include Makefile.options
finalFile=$(CURDIR)/$(dataDir)/epub-$(shell date +'%y-%m-%d_%H-%M').zip
changeSymbols=\|(\|)\|&\|'\|;
##########################################################################################
txtFiles:=$(shell find $(dataDir)/initial -type f -name '*.epub' | \
	sed "s/ $(changeSymbols)/_/g" | \
	sed 's:^$(dataDir)/initial/:$(dataDir)/working/:g' | \
	sed 's/.epub/.txt/g')
##########################################################################################
$(dataDir)/initial:
	mkdir -p $@
$(dataDir)/working:
	mkdir -p $@
$(dataDir)/final:
	mkdir -p $@
##########################################################################################
$(dataDir)/initial/.done: | $(dataDir)/initial
	tar -zxvf $(dataFile) -C $(dataDir)/initial/ --transform "s/ $(changeSymbols)/_/g"
	touch $@
extract: $(dataDir)/initial/.done

$(dataDir)/working/%.txt: $(dataDir)/initial/%.epub | $(dataDir)/initial
	mkdir -p $(dir $@)
	ebook-convert $^ $@

$(dataDir)/working/.done: $(txtFiles) | $(dataDir)/working
	touch $@
##########################################################################################
convert: $(dataDir)/working/.done
##########################################################################################
$(finalFile): $(dataDir)/working/.done
	cd $(dataDir)/working; zip -r $@ *

zip: $(finalFile)	
	echo "Final file is: $(finalFile)"
##########################################################################################
all:
	$(MAKE) extract
	$(MAKE) zip
##########################################################################################
clean:
	rm -rf $(dataDir)
##########################################################################################
.PHONY:
.DELETE_ON_ERROR: 
