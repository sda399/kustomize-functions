build:
	for i in functions/* ; do \
	   	cd ./$$i && make build; \
	done
