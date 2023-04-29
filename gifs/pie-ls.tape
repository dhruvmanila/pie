Output pie-ls.gif

Set Height 550
Set Width 1300
Set Padding 48
Set FontSize 30

# Setup
Hide
Type "mkdir {foo,bar,baz}"
Enter
Type "cd foo"
Enter
Type "pie create"
Enter
Sleep 5
Type "cd ../bar"
Enter
Type "pie create"
Enter
Sleep 5
Type "cd ../baz"
Enter
Type "pie create"
Enter
Sleep 5
Type 'source "$(pie --venv)"/bin/activate'
Enter
Type "clear"
Enter
Show

# Demo
Type "pie ls" Sleep 0.5
Enter
Sleep 2
Type "pie ls --verbose" Sleep 0.5
Enter
Sleep 5

# Teardown
Hide
Type "pie rm --yes"
Enter
Sleep 1
Type "cd ../bar"
Enter
Type "pie rm --yes"
Enter
Sleep 1
Type "cd ../foo"
Enter
Type "pie rm --yes"
Enter
Sleep 1
Type "cd .."
Enter
Type "rm -rf foo bar baz"
Enter
Show
