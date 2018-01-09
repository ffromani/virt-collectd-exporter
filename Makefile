all: virt-collectd-exporter


virt-collectd-exporter:
	@echo "Preparing output directory (_output)"
	@[ ! -d _output ] && mkdir -p _output || true
	@echo "Building..."
	@cd cmd/virt-collectd-exporter && \
		go build . && \
		mv virt-collectd-exporter ../../_output && \
		cd ../..

docker:
	@echo "Building docker image..."
	@docker build .

clean:
	@echo "Cleaning up"
	@[ -d _output ] && rm -rf _output || true


.PHONY: clean
