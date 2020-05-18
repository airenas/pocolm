-include Makefile.options
POCOLM_ROOT=$(shell cd ../..;pwd -P)
PATH:=$(PATH):$(POCOLM_ROOT)/scripts
##########################################################################################
$(dataDir)/initial:
	mkdir -p $@
$(dataDir)/working:
	mkdir -p $@	
##########################################################################################
$(dataDir)/initial/corpus.txt: | $(dataDir)/initial
	find $(corpusDir) -type f -name '*.txt' | grep 'Adm_teks' | xargs -d '\n' -n 1 cat > $@

$(dataDir)/working/corpus.clean.txt: $(dataDir)/initial/corpus.txt | $(dataDir)/working
	pv $^ | $(toolsDir)/clear.symbols > $@
$(dataDir)/working/corpus.fixquotes.txt: $(dataDir)/working/corpus.clean.txt
	pv $^ | $(toolsDir)/fix.quotes > $@
$(dataDir)/working/corpus.nopar.txt: $(dataDir)/working/corpus.fixquotes.txt
	pv $^ | $(toolsDir)/drop.parentheses > $@	
$(dataDir)/working/corpus.mapped.txt: $(dataDir)/working/corpus.nopar.txt
	pv $^ | $(toolsDir)/map.symbols > $@		
$(dataDir)/working/corpus.fixpunct.txt: $(dataDir)/working/corpus.mapped.txt
	pv $^ | $(toolsDir)/fix.punct > $@		
$(dataDir)/working/corpus.lc.txt: $(dataDir)/working/corpus.fixpunct.txt
	pv $^ | $(toolsDir)/to.lc > $@
$(dataDir)/working/corpus.ready.txt: $(dataDir)/working/corpus.lc.txt
	pv $^ | $(toolsDir)/drop.sentences > $@			


prepare-data: $(dataDir)/working/corpus.ready.txt
##########################################################################################
info:
	@echo "POCOLM_ROOT=$(POCOLM_ROOT)"
##########################################################################################
clean:
	rm -rf $(dataDir)
clean-last:
	rm -f $(dataDir)/working/corpus.lc.txt	
clean-lema:
	rm -f ~/.lema/cache/vocab
##########################################################################################
.PHONY: 
