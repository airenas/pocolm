-include Makefile.options
POCOLM_ROOT=$(shell cd ../..;pwd -P)
PATH:=$(PATH):$(POCOLM_ROOT)/scripts
##########################################################################################
$(dataDir)/initial:
	mkdir -p $@
$(dataDir)/working:
	mkdir -p $@	
##########################################################################################
$(dataDir)/initial/corpus.txt.txt: | $(dataDir)/initial
	find $(corpusDir) -type f -name '*.txt' | xargs -d '\n' -n 1 cat > $@
$(dataDir)/initial/corpus.doc.txt: | $(dataDir)/initial
	find $(corpusDir) -type f -name '*.doc' | xargs -d '\n' -n 1 catdoc > $@	
	find $(corpusDir) -type f -name '*.DOC' | xargs -d '\n' -n 1 catdoc >> $@

$(dataDir)/initial/corpus.txt: $(dataDir)/initial/corpus.txt.txt $(dataDir)/initial/corpus.doc.txt 
	cat $(dataDir)/initial/corpus.txt.txt $(dataDir)/initial/corpus.doc.txt > $@

$(dataDir)/working/corpus.clean.txt: $(dataDir)/initial/corpus.txt | $(dataDir)/working
	$(toolsDir)/clear.symbols $^ $@
$(dataDir)/working/corpus.fixquotes.txt: $(dataDir)/working/corpus.clean.txt
	$(toolsDir)/fix.quotes $^ $@	

prepare-data: $(dataDir)/working/corpus.fixquotes.txt
##########################################################################################
info:
	@echo "POCOLM_ROOT=$(POCOLM_ROOT)"
##########################################################################################
clean:
	rm -rf $(dataDir)
##########################################################################################
.PHONY: 
