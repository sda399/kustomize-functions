build:
	find functions/ -type d -mindepth 1 -maxdepth 1 -exec sh -c 'cd {} && make build' \;
