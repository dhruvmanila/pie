Output pyvenv-create.gif

Set Height 400
Set Width 1300
Set Padding 48

# Setup
Hide
Type "mkdir demo"
Enter
Type "cd demo"
Enter
Type "clear"
Enter
Show

# Demo
Type "pyvenv create" Sleep 0.5
Enter
# Wait until the environment is created and a bit more
Sleep 10

# Teardown
Hide
Type "pyvenv rm --yes"
Enter
Type "cd .."
Enter
Type "rm -rf demo"
Enter
Show
