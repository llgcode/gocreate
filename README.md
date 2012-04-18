Installation
============
    go get bitbucket.org/llg/gocreate

Usage
=====

Command line utility that create files from templates.

example usage:
    
    $ gocreate
    $ gocreate -help
    $ gocreate file -help
    $ gocreate file mymodule
	$ gocreate package mypackage

You can change template directory by defining an environment variable $GOTEMPLATE where live your templates. (by default $GOPATH/src/bitbucket.org/llg/gocreate/templates)

You may want to change the template files 'copyright' and 'author' if you don't want to have my name that appear in all your source files.


Add a template
==============

Before starting creating your template see defaults [templates].
    
    execute:
      $ gocreate template mytemplate
    modify the template: see the doc of [template package] to help using template vars
    copy the template folder to your templates folder (see gocreate -help to see where is your templates folder):
      $ mv mytemplate $GOPATH/src/bitbucket.org/llg/gocreate/templates
    or
      $ mv mytemplate $GOTEMPLATE 
    if you have defined $GOTEMPLATE environment variable

[templates]: https://bitbucket.org/llg/gocreate/src/tip/templates
[template package]: golang.org/pkg/text/template