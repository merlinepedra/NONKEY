
# del generated code 
# Get-ChildItem .\enum\ -Recurse -Include *_gen.go | Remove-Item

################################################################################
# generate enum
Write-Output "generate enums"
genenum -typename=TokenType -packagename=tokentype -basedir=enum 
genenum -typename=ObjectType -packagename=objecttype -basedir=enum 
genenum -typename=Precedence -packagename=precedence -basedir=enum 

goimports -w enum

