function fish_greeting
    echo "
Wake up, "(git config user.name)"..."
end

# SSH agent setup - reuse existing or start new
if test -z "$SSH_AUTH_SOCK"
    set -l SSH_AGENT_FILE ~/.ssh/ssh-agent-info
    if test -f $SSH_AGENT_FILE
        source $SSH_AGENT_FILE >/dev/null
        ssh-add -l &>/dev/null
        if test $status -ne 0
            eval (ssh-agent -c | tee $SSH_AGENT_FILE)
        end
    else
        eval (ssh-agent -c | tee $SSH_AGENT_FILE)
    end
end

alias ll "eza -l --icons --git"

# Jack operator — load token if initialized
if test -f ~/.jack/operator/token
    set -gx JACK_MSG_TOKEN (cat ~/.jack/operator/token)
    set -gx JACK_HOMESERVER "http://localhost:6167"
end

starship init fish | source
