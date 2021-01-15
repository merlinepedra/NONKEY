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

move static.go to package to exec "go run nonkey.go"  

add runmon for runfile, runstring with env 

update nonkey, run 1 line ,run file, repl with autostart.mon

## TODO

replace ';' with '\n' or '\r'


## bug to fix 

load data/stdlib.mon make error 

    //
    // Is the given array empty?
    //
    function array.empty?() {
    if ( len(self) == 0 ) {
        return true;
    }
    return false;
    }

    assert( "[].empty?()" );
    assert( "![1,2].empty?()" );
    assert( "![\"steve\",3].empty?()" );

    produce 
    
    Error calling object-method ERROR: Failed to invoke method: empty?
    Error calling `eval` : ERROR: Failed to invoke method: empty?
    Error calling `assert` : ERROR: Failed to invoke method: empty?