#!/bin/bash

SESSION_NAME="go_chat_dev"
tmux new-session -s $SESSION_NAME -d
#tmux split-window -h -d
#tmux split-window -v -d

PRIMARY_WINDOW_NAME="exec_code"
tmux rename-window -t $SESSION_NAME:0 $PRIMARY_WINDOW_NAME

send_str() {
	STRING=$1
	tmux send-keys -t $SESSION_NAME "$STRING "
}
send_key() {
	KEYS=$1
	tmux send-keys -t $SESSION_NAME $KEYS 
}

# set up editor
send_str "code . || codium ."
send_key ENTER
send_key C-l


# configure windows
send_str "go run main.go client 127.0.0.1 9001 alice:bob"
#tmux select-pane -t 1
tmux split-window -h 
send_str "go run main.go server 127.0.0.1 9001"
tmux select-pane -l
tmux split-window -v 
send_str "go run main.go client 127.0.0.1 9001 bob:alice"
tmux select-pane -R
# connect to session
tmux attach -t $SESSION_NAME

