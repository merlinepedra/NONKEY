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

    interpreter/evaluator/builtin_init.go

change parse function init 

    interpreter/parser/parser.go 

    prefixParseFns
    infixParseFns
    postfixParseFns

add repl from waig_code

add runmon for runfile, runstring with env 

update nonkey, run 1 line ,run file, repl

add autoload arg

add global data package builtinfunctions, pragmas

change some map to slice for performance

    interpreter/parser/parser.go

    prefixParseFns
    infixParseFns
    postfixParseFns

    enum/tokentype/attrib.go

    Token2Precedences

## TODO

replace ';' with '\n' or '\r'

del method call (tokentype PERIOD .) implement incomplete 

    If the implementation is hard to explain, it's a bad idea.

del redundant function define (tokentype DEFINE_FUNCTION function) 

    There should be one--and preferably only one--obvious way to do it.

del redundant ternary operator( tokentype QUESTION  "? :"  )

del indentifier composite char ( ? . %  ) 

identifier start char must letter and _ , not digit, following letter,_ and digit

## bug to fix 

method call act oddly