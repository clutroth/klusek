# klusek

It grabs keyboard until `Ctrl+C` is pressed and released.
It works with X11.

## build app

```
go build
```

## install app locally

```
go install
```

## Running on i3
No idea why, `bindsym $mod+c exec ~/go/bin/klusek` launches program, that doesn't grab keyboard. 
Instead, I'm running `klusek` under `xfce4-terminal` moving this terminal to unused workspace.
Ugly hack 'for now'.

Configuration `~/.config/i3/config`.
```
bindsym $mod+c exec "/usr/bin/xfce4-terminal -e ~/go/bin/klusek -T klusek"
for_window [class="^Xfce4-terminal$" title="^klusek$"], move container to workspace 9
```
