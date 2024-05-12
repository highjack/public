#!/bin/bash
echo [+] Downloading $PACKAGE $VERSION
mono /usr/bin/nuget.exe install $PACKAGE -Version $VERSION -OutputDirectory /packages -Source https://www.powershellgallery.com/api/v2/
echo [+] Grepping for potentially dangerous functionality
rm /results/* -r 2>/dev/null
rm -r ./nupkg 2>/dev/null
mkdir ./nupkg
find . -wholename "./packages/*.nupkg" -exec cp {} ./nupkg/ \;
zip -r $PACKAGE-$VERSION.zip ./nupkg
cp $PACKAGE-$VERSION.zip /results
grep -iEHnr --include='*.ps1' --include='*.psm1' -f /regexes /packages >> /results/$PACKAGE-$VERSION-report.txt
echo Done!
