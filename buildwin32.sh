#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo $DIR
cd ~
rm -rf ~/buildmineframewin32
GOOS=windows GOARCH=386 revel build github.com/MineCraftWebFrame/mineframe ~/buildmineframewin32 prod
rm -rf buildmineframewin32/src/github.com/MineCraftWebFrame/mineframe/react-project
cd ~/buildmineframewin32
printf "\r\npause\r\n" >> run.bat
BUILDZIP=mineframe_win32_`date "+%Y-%m-%d_%H-%M"`
rm ../$BUILDZIP.zip
zip -r "../$BUILDZIP.zip" * -x \*.zip
mv ../$BUILDZIP.zip $DIR
