-include Makefile.options
POCOLM_ROOT=$(shell cd ../..;pwd -P)
PATH:=$(PATH):$(POCOLM_ROOT)/scripts
##########################################################################################
max_memory='--max-memory=10G'
min_counts=''
bypass_metaparam_optim_opt=
limit_unk_history_opt=
##########################################################################################
evalTxt=$(dataDir)/eval/$(corpusAdapt).eval.txt
devTxt=$(dataDir)/text/dev.txt
lmDir=$(dataDir)/lm
devDir=$(dataDir)/dev
finalDir=$(dataDir)/final
##########################################################################################
$(dataDir)/text:
	mkdir -p $@
$(dataDir)/dwn:
	mkdir -p $@	
$(dataDir)/vocab:
	mkdir -p $@		
$(dataDir)/eval:
	mkdir -p $@			
$(finalDir):
	mkdir -p $@	
$(devDir):
	mkdir -p $@		
##########################################################################################
$(dataDir)/text/$(corpus).txt: $(dataDir)/dwn/$(corpus).tar.gz | $(dataDir)/text
	tar -zxvf $(dataDir)/dwn/$(corpus).tar.gz -C $(dataDir)/text/
	touch $(dataDir)/text/$(corpus).txt

$(dataDir)/dwn/$(corpus).tar.gz: | $(dataDir)/dwn
	$(getCorpusCmd)
##########################################################################################
$(dataDir)/dwn/$(corpusAdapt).txt: | $(dataDir)/dwn
	$(getCorpusAdaptCmd)
$(dataDir)/text/$(corpusAdapt).txt $(evalTxt) $(devTxt): $(dataDir)/dwn/$(corpusAdapt).txt | $(dataDir)/text $(dataDir)/eval
# hold out one in every 50 lines as dev data and 50+1 line as evaluation data.
	cat $(dataDir)/dwn/$(corpusAdapt).txt | \
		awk \
		'{if(NR%50 == 0) { print >"$(devTxt)"; } \
		else {if(NR%50 == 1) { print >"$(evalTxt)"; } \
		else {print;}}}' >$(dataDir)/text/$(corpusAdapt).txt
##########################################################################################
$(dataDir)/dwn/$(vocab): | $(dataDir)/dwn
	$(getVocabCmd)	
$(dataDir)/vocab/$(vocab): $(dataDir)/dwn/$(vocab) | $(dataDir)/vocab
	cp $(dataDir)/dwn/$(vocab) $(dataDir)/vocab/$(vocab)
	wc -l $(dataDir)/vocab/$(vocab)	
##########################################################################################
$(dataDir)/.data.done: $(dataDir)/text/$(corpus).txt $(dataDir)/text/$(corpusAdapt).txt $(dataDir)/vocab/$(vocab)
	touch $(dataDir)/.data.done
prepare-data: $(dataDir)/.data.done	
##########################################################################################
$(dataDir)/lm/pocolm-$(order)gram/.done: $(dataDir)/.data.done
	train_lm.py --wordlist=$(dataDir)/vocab/$(vocab) --num-splits=5 --warm-start-ratio=10 $(max_memory) \
              --min-counts=$(min_counts) \
              --keep-int-data=true  $(bypass_metaparam_optim_opt) \
              $(limit_unk_history_opt) $(dataDir)/text $(order) $(dataDir)/lm/work $(dataDir)/lm/pocolm-$(order)gram
$(finalDir)/$(finalName)-pocolm-$(order)gram.arpa.gz: $(dataDir)/lm/pocolm-$(order)gram/.done | $(finalDir)
	format_arpa_lm.py $(max_memory) $(dataDir)/lm/pocolm-$(order)gram | gzip -c > $@
build-pocolm: $(finalDir)/$(finalName)-pocolm-$(order)gram.arpa.gz
##########################################################################################
$(finalDir)/$(finalName)-pocolm-$(order)gram-pocolm.ppl: $(dataDir)/lm/pocolm-$(order)gram/.done | $(finalDir)
	get_data_prob.py $(max_memory) $(evalTxt) $(dataDir)/lm/pocolm-$(order)gram 2>&1 | grep -F '[perplexity' > $@

$(finalDir)/%.ppl: $(finalDir)/%.arpa.gz | $(finalDir)
	$(srilmDir)/ngram -unk -order $(order) -lm $< -ppl $(evalTxt) > $@

$(finalDir)/%-srilm-$(order)gram.ppl: $(lmDir)/%-$(order)gram.arpa.gz | $(finalDir) 
	$(srilmDir)/ngram -unk -order $(order) -lm $< -ppl $(evalTxt) > $@

calc-ppl-pocolm: $(finalDir)/$(finalName)-pocolm-$(order)gram-pocolm.ppl $(finalDir)/$(finalName)-pocolm-$(order)gram.ppl 
calc-ppl-srilm: $(finalDir)/$(finalName)-srilm-$(order)gram.ppl $(finalDir)/$(corpus)-srilm-$(order)gram.ppl \
				$(finalDir)/$(corpusAdapt)-srilm-$(order)gram.ppl
show-ppl: 
	@echo -n "\n\n#####Results\n"
	@ls -1 $(finalDir)/*.ppl | xargs -n1 -I {} sh -c 'echo "";echo "{}"; cat {}'	
##########################################################################################
$(lmDir):
	mkdir -p $@	
$(lmDir)/%-$(order)gram.arpa.gz: $(dataDir)/text/%.txt $(dataDir)/vocab/$(vocab) | $(lmDir)
	$(srilmDir)/ngram-count -unk -order $(order) -vocab $(dataDir)/vocab/$(vocab) -kndiscount -text $< -lm $@
##########################################################################################
$(devDir)/%-$(order)gram.ppl: $(lmDir)/%-$(order)gram.arpa.gz | $(devDir)
	$(srilmDir)/ngram -debug 2 -order $(order) -unk -lm $< -ppl $(devTxt) > $@
##########################################################################################
$(devDir)/int.values.all: $(devDir)/$(corpus)-$(order)gram.ppl $(devDir)/$(corpusAdapt)-$(order)gram.ppl
	$(srilmDir)/compute-best-mix $(devDir)/$(corpus)-$(order)gram.ppl $(devDir)/$(corpusAdapt)-$(order)gram.ppl \
		> $(devDir)/int.values.all
$(devDir)/$(corpus).int.value: $(devDir)/int.values.all
	cat $(devDir)/int.values.all | grep "best lambda" | \
		sed 's/.*(//g'| sed 's/)//g' | awk '{print($$1)}' > $(devDir)/$(corpus).int.value
	echo "Corpus interpolate value: $$(cat $(devDir)/$(corpus).int.value)" 	
compute-interpolate-values: $(devDir)/$(corpus).int.value
	echo "Corpus interpolate value: $$(cat $(devDir)/$(corpus).int.value)" 	
##########################################################################################
$(finalDir)/$(finalName)-srilm-$(order)gram.arpa.gz: $(lmDir)/$(corpus)-$(order)gram.arpa.gz $(lmDir)/$(corpusAdapt)-$(order)gram.arpa.gz $(devDir)/$(corpus).int.value
	$(srilmDir)/ngram -order $(order) -unk \
  		-lm $(lmDir)/$(corpus)-$(order)gram.arpa.gz -lambda $$(cat $(devDir)/$(corpus).int.value) \
  		-mix-lm $(lmDir)/$(corpusAdapt)-$(order)gram.arpa.gz \
		-write-lm $(finalDir)/$(finalName)-srilm-$(order)gram.arpa.gz
build-srilm: $(finalDir)/$(finalName)-srilm-$(order)gram.arpa.gz
##########################################################################################
info:
	@echo "POCOLM_ROOT=$(POCOLM_ROOT)"
	@echo "ngram order=$(order)"
##########################################################################################
clean-final:
	rm -rf $(finalDir)
clean:
	rm -rf $(dataDir)/text
	rm -rf $(dataDir)/.data.done
	rm -rf $(dataDir)/vocab
	rm -rf $(dataDir)/eval
	rm -rf $(dataDir)/lm
	rm -rf $(dataDir)/dev
clean-all:
	rm -rf $(dataDir)
##########################################################################################
.PHONY: calc-ppl-srilm calc-ppl-pocolm

