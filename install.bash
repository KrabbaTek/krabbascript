
# In short, this script will install the KrabbaScript compiler from https://github.com/KrabbaTek/krabbascript
# It will do the following things:
# - run `go install https://github.com/KrabbaTek/krabbascript`
# - set up the KrabbaScript modules directory (~/.kscript/libs) if it doesn't already exist
# - prompt you to install Cranium (the KrabbaScript package manager — more information can be found here: https://github.com/KrabbaTek/Cranium-cli)

# Written by Reece Turner


# check if Golang is installed, since we cannot run `go install` without Golang
if ! command -v go >/dev/null 2>&1
then
    echo "You don't have Golang installed."
    echo "Sadly, I cannot install this for you, so you'll need to visit https://go.dev/doc/install and install Golang yourself."
    exit 1
fi

# install KrabbaScript if not already installed
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
    if ! grep -qs 'go/bin' ~/.zshrc; then
        echo "KrabbaScript was installed, but it has not yet been added to your PATH."
        read -p "Would you like to add ~/go/bin to your PATH? [Y/n] ~> " path_answer

        if [[ "{$path_answer,,}" == "y" || -z "$input" ]]
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
    read -p "You should skip this if you don't plan on using any external libraries, otherwise I recommend creating it. [Y/n] ~> " input

    if [[ "{$input,,}" == "y" || -z "$input" ]]
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
    echo $'\n\nYou don\'t have Cranium installed.'
    echo "Cranium is a package manager for KrabbaScript"
    echo "You can also install the Off the Grid version, which is a purely local package manager."
    
    echo "You can read more at https://github.com/KrabbaTek/Cranium-cli and https://github.com/KrabbaTek/Cranium-off-the-grid respectively."

    # this is the stupidest system I've designed in my entire life, but this is a pretty tight constraint.
    # lol that's a stupid excuse — this is stupid and theres no way of justifying it.
    read -p "Would you like to install Cranium? [P]lain version. | [O]ff the Grid version | [N]o, don't install Cranium ~> " installation

    if [[ "{$installation,,}" == "p" ]] # normal Cranium
    then
        echo "placeholder — plain"
    elif [[ "{$installation,,}" == "o" ]] # off the grid
    then
        echo "placeholder — off the grid"
    fi
fi


