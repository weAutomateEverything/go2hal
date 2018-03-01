#! /bin/bash
for pkg in $(go list ./... | grep -v vendor); do
     if ! go test -coverprofile=$(echo $pkg | tr / -).cover $pkg; then
        exit 1
     fi
done
echo "mode: set" > c.out
grep -h -v "^mode:" ./*.cover >> c.out
rm -f *.cover