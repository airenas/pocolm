-include Makefile.options
POCOLM_ROOT=$(shell cd ../..;pwd -P)
PATH:=$(PATH):$(POCOLM_ROOT)/scripts
##########################################################################################
$(dataDir)/initial:
	mkdir -p $@
##########################################################################################
$(dataDir)/initial/corpus.txt.txt: | $(dataDir)/initial
	find $(corpusDir) -type f -name '*.txt' | xargs -d '\n' -n 1 cat > $@
$(dataDir)/initial/corpus.doc.txt: | $(dataDir)/initial
	find $(corpusDir) -type f -name '*.doc' | xargs -d '\n' -n 1 catdoc > $@	
	find $(corpusDir) -type f -name '*.DOC' | xargs -d '\n' -n 1 catdoc >> $@

$(dataDir)/initial/corpus.txt: $(dataDir)/initial/corpus.txt.txt $(dataDir)/initial/corpus.doc.txt 
	cat $(dataDir)/initial/corpus.txt.txt $(dataDir)/initial/corpus.doc.txt > $@

prepare-data: $(dataDir)/initial/corpus.txt
##########################################################################################
info:
	@echo "POCOLM_ROOT=$(POCOLM_ROOT)"
##########################################################################################
clean:
	rm -rf $(dataDir)
##########################################################################################
.PHONY: 
