# line2json

[NAME](#NAME)  
[SYNOPSIS](#SYNOPSIS)  
[DESCRIPTION](#DESCRIPTION)  
[OPTIONS](#OPTIONS)  
[EXAMPLE](#EXAMPLE)  
[SEE ALSO](#SEE%20ALSO)  
[AUTHOR](#AUTHOR)  

------------------------------------------------------------------------

## NAME <span id="NAME"></span>

line2json − convert text lines to JSON

## SYNOPSIS <span id="SYNOPSIS"></span>

**line2json** \[**−−keyRegex** *regex*\] \[**−−keyReplace**
*replacement*\] \[**−−object**\] \[**−−valueRegex** *regex*\]
\[**−−valueReplace** *replacement*\]

## DESCRIPTION <span id="DESCRIPTION"></span>

**line2json** is a text lines to JSON converter, originally designed to
work with *cmenu*(1) and *man*(1).

## OPTIONS <span id="OPTIONS"></span>

**−−keyRegex** *regex*

The regular expression to manipulate item if outputting an array or
manipulate key if outputting an object.

**−−keyReplace** *replacement*

The replacement string when key regular expression matches.

**−−object**

Whether to convert each line to a key−value object.

**−−valueRegex** *regex*

The regular expression to manipulate value if outputting an object.

**−−valueReplace** *replacement*

The replacement string when value regular expression matches.

## EXAMPLE <span id="EXAMPLE"></span>

Script for listing all man pages:

**\$ eval "\$(man −k . \| line2json −−object −−valueRegex '^(.+)
\\(.+)\\\[\[:space:\]\]\*−.\*\$' −valueReplace 'man \$2 \$1' \| cmenu
fzf)"**

## SEE ALSO <span id="SEE ALSO"></span>

***cmenu***(1), *man*(1), *fzf*(1)

## AUTHOR <span id="AUTHOR"></span>

Kris Andrie Ortega (andrieee44@gmail.com)

------------------------------------------------------------------------
