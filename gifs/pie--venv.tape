Output pie--venv.gif

Set Height 500
Set Width 1300
Set Padding 48
Set FontSize 32

# Setup
Hide
Type "mkdir demo"
Enter
Type "cd demo"
Enter
Type "pie create"
Enter
Sleep 5
Type "clear"
Enter
Show

# Demo
Type "pie --venv" Sleep 0.5
Enter
Sleep 1
Type 'source "$(pie --venv)"/bin/activate'
Enter
Sleep 1
Type 'echo "$VIRTUAL_ENV"'
Enter
Sleep 3

# Teardown
Hide
Type "pie rm --yes"
Enter
Sleep 1
Type "cd .."
Enter
Type "rm -rf demo"
Enter
Show
