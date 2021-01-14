
# del generated code 
# Get-ChildItem .\enum\ -Recurse -Include *_gen.go | Remove-Item

################################################################################
# generate enum
Write-Output "generate enums"
genenum -typename=TokenType -packagename=tokentype -basedir="." -vectortype=int

goimports -w tokentype

