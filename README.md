# nonkey
customized monkey interpreter 

from https://github.com/skx/monkey/

from https://interpreterbook.com/

## changed 

changed to typed enum by genenum (https://github.com/kasworld/genenum)

    tokentype
    objecttype
    precedence

change builtin function init 

    evaluator/builtinfunctions.go

change parse function init 

    parser/parser.go 
    prefixParseFns
    infixParseFns
    postfixParseFns

add repl from waig_code

## TODO

replace ';' with '\n' or '\r'


