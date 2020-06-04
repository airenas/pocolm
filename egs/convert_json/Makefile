##########################################################################################
##########################################################################################
# Build tools (json.to.files) before use 
# cd ../../tools && make clean build
##########################################################################################
# To convert jsons - set input zip file in Makefile.options and run 
# make all -j4
##########################################################################################
##########################################################################################
-include Makefile.options
finalFile=$(CURDIR)/$(dataDir)/ziniasklaida-$(shell date +'%y-%m-%d_%H-%M').zip
POCOLM_ROOT=$(shell cd ../..;pwd -P)
PATH:=$(PATH):$(POCOLM_ROOT)/tools/bin
##########################################################################################
txtFiles:=$(shell find $(dataDir)/initial -type f -name '*.json' | \
	sed 's:^$(dataDir)/initial/:$(dataDir)/working/:g' | \
	sed 's:.json:/.done:g')
##########################################################################################
$(dataDir)/initial:
	mkdir -p $@
$(dataDir)/working:
	mkdir -p $@
$(dataDir)/final:
	mkdir -p $@
##########################################################################################
$(dataDir)/initial/.done: | $(dataDir)/initial
	unzip $(dataFile) -d $(dataDir)/initial/ 
	touch $@
extract: $(dataDir)/initial/.done

$(dataDir)/working/%/.done: $(dataDir)/initial/%.json | $(dataDir)/working
	mkdir -p $(dir $@)
	json.to.files -d $(dir $@) $^

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
