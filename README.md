Installation
============
    go get bitbucket.org/llg/gocreate

Usage
=====

Command line utility that create files from templates.

example usage:
    
    $ gocreate
    $ gocreate -help
    $ gocreate main -help
    $ gocreate main

You can change template directory by defining an environment variable $GOTEMPLATE where live your templates. (by default $GOPATH/src/bitbucket.org/llg/gocreate/templates)


Add a template
==============

Before starting creating your template see defaults template:
 https://bitbucket.org/llg/gocreate/src/tip/templates
    
    execute:
      $ gocreate template -name mytemplate
    modify the template
    copy the template folder to your templates folder (see gocreate -help to see where is your templates folder):
      $ mv mytemplate $GOPATH/src/bitbucket.org/llg/gocreate/templates
    or
      $ mv mytemplate $GOTEMPLATE 
    if you have defined $GOTEMPLATE environment variable

