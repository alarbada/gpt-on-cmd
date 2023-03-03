# gpt-on-cmd

Right now on MVP status, but the main idea is to have gpt on your own CMD,
thanks to the new `gpt-3.5-turbo` model that just came as of the time of
writing this.

# TODO

- [ ] Properly capture a new line. What if I press ctrl-D?
- [ ] What happens if I write while I'm receiving input from openai?
- [ ] Create some sort of menu so that one can change the default parameters of the chat (temperature and such)
- [ ] Deploy binaries on github
- [ ] Ugh, handle windows :(
    - Does the default implementation work on windows?
    - Put something on the README that says "Tested on linux only"
