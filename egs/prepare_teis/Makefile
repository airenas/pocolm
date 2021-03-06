-include Makefile.options
POCOLM_ROOT=$(shell cd ../..;pwd -P)
PATH:=$(PATH):$(POCOLM_ROOT)/scripts
finalFile=$(dataDir)/teis-corpus-$(shell date +'%y-%m-%d_%H-%M').zip
##########################################################################################
$(dataDir)/initial:
	mkdir -p $@
$(dataDir)/working:
	mkdir -p $@
$(dataDir)/final:
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
	pv $^ | $(toolsDir)/clear.symbols > $@
$(dataDir)/working/corpus.fixquotes.txt: $(dataDir)/working/corpus.clean.txt
	pv $^ | $(toolsDir)/fix.quotes > $@
$(dataDir)/working/corpus.nopar.txt: $(dataDir)/working/corpus.fixquotes.txt
	pv $^ | $(toolsDir)/drop.parentheses > $@	
$(dataDir)/working/corpus.mapped.txt: $(dataDir)/working/corpus.nopar.txt
	pv $^ | $(toolsDir)/map.symbols -email -url -slash > $@		
$(dataDir)/working/corpus.fixpunct.txt: $(dataDir)/working/corpus.mapped.txt
	pv $^ | $(toolsDir)/fix.punct -abbr $(abbrFile) > $@		
$(dataDir)/working/corpus.lc.txt: $(dataDir)/working/corpus.fixpunct.txt
	pv $^ | $(toolsDir)/to.lc -abbr $(abbrFile) > $@
$(dataDir)/working/corpus.ready.txt: $(dataDir)/working/corpus.lc.txt
	pv $^ | $(toolsDir)/drop.sentences > $@	
$(dataDir)/working/corpus.ready.np.txt: $(dataDir)/working/corpus.ready.txt
	pv $^ | $(toolsDir)/drop.punct -abbr $(abbrFile) > $@
$(dataDir)/final/corpus.txt: $(dataDir)/working/corpus.ready.txt | $(dataDir)/final
	cp $^ $@		
	cp $^ $(dataDir)/final/dev.txt
$(dataDir)/final/corpus.np.txt: $(dataDir)/working/corpus.ready.np.txt | $(dataDir)/final
	cp $^ $@		

prepare-data: $(dataDir)/final/corpus.txt $(dataDir)/vocab/corpus.ord.counts \
	$(dataDir)/vocab/corpus.ord.vocab $(dataDir)/vocab/missing.counts
##########################################################################################
$(dataDir)/vocab/corpus.np.counts: $(dataDir)/final/corpus.np.txt | $(dataDir)
	get_word_counts.py $(dataDir)/final $(dataDir)/vocab
$(dataDir)/vocab/corpus.ord.counts: $(dataDir)/vocab/corpus.np.counts
	cat $^ | sort -k1,1nr -k2,2 > $@
$(dataDir)/vocab/corpus.ord.vocab: $(dataDir)/vocab/corpus.np.counts
	cat $^ | sort -k2,2 > $@	
$(dataDir)/vocab/missing.counts: $(dataDir)/vocab/corpus.ord.counts
	pv $^ | $(toolsDir)/not.in.voc -v $(mainVocab) > $@
##########################################################################################
$(dataDir)/final/.done: $(dataDir)/final/corpus.txt $(dataDir)/vocab/corpus.ord.counts \
	$(dataDir)/vocab/corpus.ord.vocab $(dataDir)/vocab/missing.counts
	cp $(dataDir)/vocab/corpus.ord.counts $(dataDir)/final/
	cp $(dataDir)/vocab/missing.counts $(dataDir)/final/
	rm -f $(dataDir)/final/dev.txt
	touch $@
$(finalFile): $(dataDir)/final/.done
	zip -r $@  $(dataDir)/final

prepare-final: $(finalFile)	
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
.DELETE_ON_ERROR: 
