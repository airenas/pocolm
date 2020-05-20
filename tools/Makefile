tools=clear.symbols fix.quotes drop.parentheses map.symbols to.lc fix.punct \
	drop.sentences

#######################################################################
dist_dir=$(CURDIR)/bin
dist_tools_files=$(patsubst %, $(dist_dir)/%, $(tools))
#######################################################################
$(dist_dir): 
	mkdir -p $@
#######################################################################
$(dist_dir)/%: cmd/% | $(dist_dir)
	cd $^ && go build -o $@
#######################################################################
build: $(dist_tools_files)
#######################################################################
benchmark:
	go test ./.../... -bench=Benchmark
benchmark-cpu:
	go test ./.../$(gpkg)/... -bench=Benchmark -cpuprofile=.cpu.prof
view-cpu:
	go tool pprof -png .cpu.prof > .cpu.png
	google-chrome .cpu.png	
test:
	gotestsum
#######################################################################
clean:
	rm -rf $(dist_dir)
#######################################################################
