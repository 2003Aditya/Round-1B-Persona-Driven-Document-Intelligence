!/bin/bash

tmux new-session -d -s 1B

tmux send-keys "nvim ~/adobe/1B" C-m
tmux rename-window "Code"

tmux new-window -t 1B:2 -n "term"
tmux send-keys "nvim ~/adobe/1B -c 'terminal'" C-m

tmux attach -t 1B
