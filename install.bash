#!/usr/bin/env bash


# In short, this script will install the KrabbaScript compiler from https://github.com/KrabbaTek/krabbascript
# It will do the following things:
# - run `go install https://github.com/KrabbaTek/krabbascript`
# - set up the KrabbaScript modules directory (~/.kscript/libs) if it doesn't already exist
# - prompt you to install Cranium (the KrabbaScript package manager — more information can be found here: https://github.com/KrabbaTek/Cranium-cli)

# To get and run this script in a pretentious one-liner, run:
# curl -fsSL https://raw.githubusercontent.com/Centurion774477/krabbascript/refs/heads/dev/install.bash | bash

# Written by Reece Turner


# check if Golang is installed, since we cannot run `go install` without Golang
if ! command -v go >/dev/null 2>&1
then
    echo "You don't have Golang installed."
    echo "Sadly, I cannot install this for you, so you'll need to visit https://go.dev/doc/install and install Golang yourself."
    exit 1
fi

# install KrabbaScript if not already installed
# somehow this check doesn't work.
# I'm keeping it nevertheless because it could work on somebody else's machine
if ! command -v krabbascript >/dev/null 2>&1
then
    # I might want to pipe the output of this into /dev/null
    go install github.com/KrabbaTek/krabbascript@latest

    # if command failed
    if [ $? -ne 0 ]
    then
        echo "Something went wrong installing KrabbaScript."
        exit 1
    fi


    # add ~/go/bin to PATH unless it's already there
    # todo: fix this to check for other shells
    if ! grep -qs 'go/bin' ~/.zshrc; then
        echo "KrabbaScript was installed, but it has not yet been added to your PATH."
        read -p "Would you like to add ~/go/bin to your PATH? [Y/n] ~> " path_answer < /dev/tty


        # downcase the input
        path_answer=$(echo "$path_answer" | tr '[:upper:]' '[:lower:]')
        if [[ "$path_answer" == "y" || -z "$input" ]]
        then
            # find the correct config file
            for RC_FILE in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile"
            do
                # next unless this config file exists
                [[ -f "$RC_FILE" ]] || continue
                if ! grep -qs 'go/bin' "$RC_FILE"; then
                    echo 'export PATH="$PATH:$HOME/go/bin"' >> "$RC_FILE"
                fi
            done

            echo "Added go/bin to your path. Remember to restart your terminal for the changes to be applied."
        else
            echo "Understood, just remember to add KrabbaScript to your PATH. You can find the code in ~/go/bin"
        fi
    fi
fi


# create libs directory if it doesn't exist
if [ ! -d ~/".kscript/libs" ]
then
    echo $'\n\nWould you like to create the libs directory in ~/.kscript/libs?'
    echo "This will store any libraries installed from Cranium or anywhere else."
    read -p "You should skip this if you don't plan on using any external libraries, otherwise I recommend creating it. [Y/n] ~> " input < /dev/tty

    input=$(echo "$input" | tr '[:upper:]' '[:lower:]')
    if [[ "$input" == "y" || -z "$input" ]]
    then
        mkdir -p ~/".kscript/libs"

        # check if the command ran successfully
        if [ $? -ne 0 ]
        then
            echo "Sorry, but something went wrong creating ~/.kscript/libs."
        else
            echo "Successfully created libraries directory."
        fi
    fi
fi


# check if Cranium is installed, and don't prompt them if it is
if ! command -v cranium >/dev/null 2>&1
then
    echo "You might want to consider installing Cranium — the package manager for KrabbaScript."
    echo "More information — and an install script — for Cranium can be found here: https://github.com/KrabbaTek/Cranium-cli"
    echo "If you want a purely local version, then check out Cranium Off the Grid here: https://github.com/KrabbaTek/Cranium-off-the-grid"
fi

echo "Kiitos paljon, hei hei!"
