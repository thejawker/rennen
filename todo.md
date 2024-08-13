# dos that you must do, ok?
> in order for this to be a proper v1 release, the following must be done, at the   
> very very _very_    
> _v e r y_   
> *very* least:

## most important
- [x] clear notification bell right away when viewing it
- [x] ability to restart an individual task (e.g. by pressing 'r')
- [x] ability to gracefully exit the program (e.g. by pressing 'q')
- [x] ability to gracefully stop one task (e.g. by pressing 'x')
- [x] fix tab bar width overflow issue
- [x] fix the ugly asni line clearing that messes up the content (looking at you yarn)
- [ ] add a predefined set of commands you can easily kick off from the overview page

## onboarding --- on_boooooring_ jk jk
- [x] add an init command that creates a `ren.json` file
- [x] suggest init if no ren.json

## meh dont really care
- [ ] load stuff from the package.json ????? not sure tho
- [ ] add a super minimal view that kills all the borders and spacing and just shows the shiz
- [x] enable ansi colors from the processes 
- [x] in process: info line at the bottom right of the screen (e.g. "Press 'q' to quit" or "Press 'r' to restart")
- [x] overview should show
  - [x] a selector to trigger shortcuts/commands eg like open browser db etc
  - [x] table of all running processes and commands and their last output
  - [x] a hint line
- [x] disable logging by default
- [ ] scrollable content
- [ ] search
- [x] ability to clear the screen (e.g. by pressing 'c')