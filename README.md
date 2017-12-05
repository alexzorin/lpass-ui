# lpass-ui

## Deprecation Notice
This project is no longer being worked on, as I have [moved onto a different solution](https://github.com/alexzorin/i3-lastpass).

![gif preview](http://i.imgur.com/YkMK8Fz.gif)

This is a very thin feature-bare QML over the top of [lpass/lastpass-cli](https://github.com/lastpass/lastpass-cli/) for use with Linux (though it may work with other platforms, there's no Linux specific code).

The motivation for this was to have a quick way to access LastPass vault passwords on Linux that did not involve using the browser extensions (which, imo, at this stage are too buggy and unsafe to use). `lpass` is great but it is a pain to have to type stuff into a terminal every time I need to login somewhere.

Because this is just a wrapper around `lpass`, all problems of authentication and talking to the LastPass API are avoided.

The intended way the program is executed goes something like:

- You have a global hotkey to start the program
- You type in your search term (which is directly passed to `lpass`)
- The results are rendered in the table as you type 
- Press <kbd>tab</kbd> to navigate to the result you want
- Press <kbd>return/enter</kbd> to copy the password for that result to the system clipboard
- The program immediately exits (press <kbd>Esc</kbd> at any time to bail)

## Installation
(Assuming you already have `lpass` setup).

Again, only ever ran on Linux, other platforms probably work fine.

### Deps
You need Qt. [Look here](https://github.com/go-qml/qml#requirements-on-ubuntu).

### Application

```
go get -u github.com/alexzorin/lpass-ui
```

The repo is a bit fat because I had to vendor a branch of go-qml that works with Go 1.8.

### Hotkey
Setup a global hotkey (this is out of scope for the project, sorry). I use i3wm so it is a simple:

```
bindsym $mod+x exec $GOPATH/bin/lpass-ui
```

aaand you are done.

## Caveats
- If you have credentials that are "protected"/require re-authentication, when you search for them, it will cause your system key agent to re-prompt you for your LastPass master password. Sorry, but this is a problem of LastPass itself + `lpass`. Feel free to press cancel when it prompts you.
- No secure notes (maybe for the future, but definitely not files)

## Licence
MIT
