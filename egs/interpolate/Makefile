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
trainTxt=$(dataDir)/text/train.txt
lmDir=$(dataDir)/lm
devDir=$(dataDir)/dev
finalDir=$(dataDir)/final
vocab=vocab.txt
unkSym='<unk>'
##########################################################################################
$(dataDir)/text:
	mkdir -p $@
$(dataDir)/text.pocolm:
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
$(dataDir)/text/$(corpus).txt: $(corpusDir)/$(corpus).txt | $(dataDir)/text
	ln -s $^ $@
##########################################################################################
$(dataDir)/.text.done: $(corpusDir)/$(corpusAdapt).txt $(dataDir)/text/$(corpus).txt| $(dataDir)/text $(dataDir)/eval
# hold out one in every 20 lines as dev data and as evaluation data.
# use 5 % for dev and 5% for test
	pv $(corpusDir)/$(corpusAdapt).txt | \
		awk \
		'{if(NR%20 == 0) { print >"$(devTxt)"; } \
		else {if(NR%20 == 1) { print >"$(evalTxt)"; } \
		else {print;}}}' >$(trainTxt)
	cp $(corpusDir)/$(corpusAdapt).txt $(dataDir)/text/$(corpusAdapt).txt
	touch $@

##########################################################################################
$(dataDir)/.text.pocolm.done: $(dataDir)/.text.done | $(dataDir)/text.pocolm
	ln -s $(corpusDir)/$(corpus).txt  $(dataDir)/text.pocolm/$(corpus).txt
	ln -s $(dataDir)/text/train.txt  $(dataDir)/text.pocolm/train.txt
	ln -s $(dataDir)/text/dev.txt  $(dataDir)/text.pocolm/dev.txt
	touch $@

##########################################################################################
$(dataDir)/vocab/$(vocab): $(corpusDir)/$(corpus).txt $(corpusDir)/$(corpus1).txt $(corpusDir)/$(corpus2).txt | $(dataDir)/vocab
	pv $(corpusDir)/$(corpus).txt $(corpusDir)/$(corpus1).txt $(corpusDir)/$(corpus2).txt | \
		python local/prepare-vocab.py --add-unk $(unkSym) --add-eps '<eps>' > $@
	wc -l $@

$(dataDir)/vocab/vocab-all.txt: $(corpusDir)/$(corpus).txt $(corpusDir)/$(corpus1).txt $(corpusDir)/$(corpus2).txt | $(dataDir)/vocab
	pv $(corpusDir)/$(corpus).txt $(corpusDir)/$(corpus1).txt $(corpusDir)/$(corpus2).txt | \
		python local/prepare-vocab.py --add-unk $(unkSym) --add-eps '<eps>' --all-words> $@
	wc -l $@		
##########################################################################################
$(dataDir)/.data.done: $(dataDir)/vocab/$(vocab) $(dataDir)/vocab/vocab-all.txt \
					$(dataDir)/.text.done $(dataDir)/.text.pocolm.done
	touch $(dataDir)/.data.done
prepare-data: $(dataDir)/.data.done	
##########################################################################################
$(dataDir)/lm/pocolm-$(order)gram/.done: $(dataDir)/.data.done
	train_lm.py --wordlist=$(dataDir)/vocab/$(vocab) --num-splits=5 --warm-start-ratio=10 $(max_memory) \
              --min-counts=$(min_counts) \
              --keep-int-data=true  $(bypass_metaparam_optim_opt) \
              $(limit_unk_history_opt) $(dataDir)/text.pocolm $(order) $(dataDir)/lm/work $(dataDir)/lm/pocolm-$(order)gram
$(finalDir)/$(finalName)-pocolm-$(order)gram.arpa.gz: $(dataDir)/lm/pocolm-$(order)gram/.done | $(finalDir)
	format_arpa_lm.py $(max_memory) $(dataDir)/lm/pocolm-$(order)gram | gzip -c > $@
build-pocolm: $(finalDir)/$(finalName)-pocolm-$(order)gram.arpa.gz
##########################################################################################
$(finalDir)/$(finalName)-pocolm-$(order)gram-pocolm.ppl: $(dataDir)/lm/pocolm-$(order)gram/.done | $(finalDir)
	get_data_prob.py $(max_memory) $(evalTxt) $(dataDir)/lm/pocolm-$(order)gram 2>&1 | grep -F '[perplexity' > $@

$(finalDir)/%.ppl: $(finalDir)/%.arpa.gz | $(finalDir)
	$(srilmDir)/ngram -unk -map-unk $(unkSym) -order $(order) -lm $< -ppl $(evalTxt) > $@

$(finalDir)/%-srilm-$(order)gram.ppl: $(lmDir)/%-$(order)gram.arpa.gz | $(finalDir) 
	$(srilmDir)/ngram -unk -map-unk $(unkSym) -order $(order) -lm $< -ppl $(evalTxt) > $@

calc-ppl-pocolm: $(finalDir)/$(finalName)-pocolm-$(order)gram-pocolm.ppl $(finalDir)/$(finalName)-pocolm-$(order)gram.ppl 
calc-ppl-srilm: $(finalDir)/$(finalName)-srilm-$(order)gram.ppl $(finalDir)/$(corpus)-srilm-$(order)gram.ppl \
				$(finalDir)/$(corpusAdapt)-srilm-$(order)gram.ppl $(finalDir)/train-srilm-$(order)gram.ppl
show-ppl: 
	@echo -n "\n\n#####Results\n"
	@ls -1 $(finalDir)/*.ppl | xargs -n1 -I {} sh -c 'echo "";echo "{}"; cat {}'	
##########################################################################################
$(lmDir):
	mkdir -p $@	
$(lmDir)/%-$(order)gram.arpa.gz: $(dataDir)/.data.done | $(lmDir)
	$(srilmDir)/ngram-count -unk -map-unk $(unkSym) -order $(order) -vocab $(dataDir)/vocab/$(vocab) -kndiscount \
		-text $(dataDir)/text/$*.txt -lm $@
##########################################################################################
$(devDir)/%-$(order)gram.ppl: $(lmDir)/%-$(order)gram.arpa.gz | $(devDir)
	$(srilmDir)/ngram -debug 2 -order $(order) -unk -map-unk $(unkSym) -lm $< -ppl $(devTxt) > $@
##########################################################################################
$(devDir)/int.values.all: $(devDir)/$(corpus)-$(order)gram.ppl $(devDir)/train-$(order)gram.ppl
	$(srilmDir)/compute-best-mix $(devDir)/$(corpus)-$(order)gram.ppl $(devDir)/train-$(order)gram.ppl \
		> $(devDir)/int.values.all
$(devDir)/$(corpus).int.value: $(devDir)/int.values.all
	cat $(devDir)/int.values.all | grep "best lambda" | \
		sed 's/.*(//g'| sed 's/)//g' | awk '{print($$1)}' > $(devDir)/$(corpus).int.value
	echo "Corpus interpolate value: $$(cat $(devDir)/$(corpus).int.value)" 	
compute-interpolate-values: $(devDir)/$(corpus).int.value
	echo "Corpus interpolate value: $$(cat $(devDir)/$(corpus).int.value)" 	
##########################################################################################
$(finalDir)/$(finalName)-srilm-$(order)gram.arpa.gz: $(lmDir)/$(corpus)-$(order)gram.arpa.gz \
		$(lmDir)/$(corpusAdapt)-$(order)gram.arpa.gz $(devDir)/$(corpus).int.value | $(finalDir)
	$(srilmDir)/ngram -order $(order) -unk -map-unk $(unkSym) \
  		-lm $(lmDir)/$(corpus)-$(order)gram.arpa.gz -lambda $$(cat $(devDir)/$(corpus).int.value) \
  		-mix-lm $(lmDir)/$(corpusAdapt)-$(order)gram.arpa.gz \
		-write-lm $(finalDir)/$(finalName)-srilm-$(order)gram.arpa.gz
build-srilm: $(finalDir)/$(finalName)-srilm-$(order)gram.arpa.gz

build: build-srilm calc-ppl-srilm
##########################################################################################
info:
	@echo "POCOLM_ROOT=$(POCOLM_ROOT)"
	@echo "ngram order=$(order)"
##########################################################################################
put-arpa: 
	$(putCmd) $(finalDir)/*.arpa.gz $(putDir)

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

