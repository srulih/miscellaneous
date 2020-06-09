program = line { line } <EOF>.
line = [ assignment | print | reset ] "\n" | ";".
assignment = type va ":=" expression.
type = "int" | "float"
print = "PRINT" expression.
reset = "RESET".
expression = term { addop term }.
term = power { mulop power }.
power = factor { powop factor }.
factor = "(" expression ")" | var | number.
addop = "+" | "-".
mulop = "*" | "/".
powop = "**" 
