Create files from templates
===========================

Command line utility that create files from templates.

example usage:
	gocreate -help
	gocreate main -help
	gocreate main

Change template directory:
	- define an environment variable $GOTEMPLATE where live your templates

Add a template:
	- execute:
		gocreate template -name mytemplate
	- modify the template
	- copy the template folder to your templates folder (see gocreate -help to see where is your templates folder)
