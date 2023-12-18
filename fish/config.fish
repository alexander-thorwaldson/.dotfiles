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

starship init fish | source
