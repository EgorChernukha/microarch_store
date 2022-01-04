

# List of command package names,
#  each one builds from Go package at 'cmd/$NAME' to executable at 'bin/$NAME'
APP_CMD_NAMES = \
	user \
	auth \
	order \
	billing \
	notification

# Contains common make targets, including 'build', 'test' and 'check'
include make/rules.mk
