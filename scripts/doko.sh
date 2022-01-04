#!/bin/bash

dokoPath=/usr/local/bin
UNAME=$(uname)
ARCH=$(uname -m)

rmOldFiles() {
    if [ -f $dokoPath/doko ]; then
        sudo rm -rf $dokoPath/doko*
    fi
}

v=$(curl --silent "https://api.github.com/repos/abdfnx/doko/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

releases_api_url=https://github.com/abdfnx/doko/releases/download

successInstall() {
    echo "🙏 Thanks for installing Doko 🐋! If this is your first time using the CLI, be sure to run `doko --help` first."
}

mainCheck() {
    echo "Installing doko version $v"
    name=""

    if [ "$UNAME" == "Linux" ]; then
        if [ $ARCH = "x86_64" ]; then
            name="doko_linux_${v}_amd64"
        elif [ $ARCH = "i686" ]; then
            name="doko_linux_${v}_386"
        elif [ $ARCH = "i386" ]; then
            name="doko_linux_${v}_386"
        elif [ $ARCH = "arm64" ]; then
            name="doko_linux_${v}_arm64"
        elif [ $ARCH = "arm" ]; then
            name="doko_linux_${v}_arm"
        fi

        dokoURL=$releases_api_url/$v/$name.zip

        wget $dokoURL
        sudo chmod 755 $name.zip
        unzip $name.zip
        rm $name.zip

        # doko
        sudo mv $name/bin/doko $dokoPath

        rm -rf $name

    elif [ "$UNAME" == "Darwin" ]; then
        if [ $ARCH = "x86_64" ]; then
            name="doko_macos_${v}_amd64"
        elif [ $ARCH = "arm64" ]; then
            name="doko_macos_${v}_arm64"
        fi

        dokoURL=$releases_api_url/$v/$name.zip

        wget $dokoURL
        sudo chmod 755 $name.zip
        unzip $name.zip
        rm $name.zip

        # doko
        sudo mv $name/bin/doko $dokoPath

        rm -rf $name

    elif [ "$UNAME" == "FreeBSD" ]; then
        if [ $ARCH = "x86_64" ]; then
            name="doko_freebsd_${v}_amd64"
        elif [ $ARCH = "i386" ]; then
            name="doko_freebsd_${v}_386"
        elif [ $ARCH = "i686" ]; then
            name="doko_freebsd_${v}_386"
        elif [ $ARCH = "arm64" ]; then
            name="doko_freebsd_${v}_arm64"
        elif [ $ARCH = "arm" ]; then
            name="doko_freebsd_${v}_arm"
        fi

        dokoURL=$releases_api_url/$v/$name.zip

        wget $dokoURL
        sudo chmod 755 $name.zip
        unzip $name.zip
        rm $name.zip

        # doko
        sudo mv $name/bin/doko $dokoPath

        rm -rf $name
    fi

    # chmod
    sudo chmod 755 $dokoPath/doko
}

rmOldFiles
mainCheck

if [ -x "$(command -v doko)" ]; then
    successInstall
else
    echo "Download failed 😔"
    echo "Please try again."
fi
