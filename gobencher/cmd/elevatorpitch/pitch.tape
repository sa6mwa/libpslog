Output elevatorpitch.gif
Set Width 756
Set Height 560
Set FontSize 11
Set Margin 4
Set Padding 10
Set MarginFill "#674EFF"
Set BorderRadius 12

Set Shell "bash"
Env PS1 "\[\e[0;34m\]❯\[\e[0m\] "

Set TypingSpeed 0.05

Type "./c-pslog-vs-go-pslog"
Enter
Wait+Line /^❯\s*$/
Sleep 10s
