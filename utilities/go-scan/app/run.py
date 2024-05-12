import os

url = os.getenv("url")
pkg_ver = os.getenv("pkg_ver")
url = url.replace("https://","")
url = url.replace("http://","")
url = url.rstrip("/")

if url != None and pkg_ver != None and len(url) > 0 and len(pkg_ver) > 0:
    gomod = """
    module test/test

    go 1.19

    require  {} {}
    """.format(url,pkg_ver)
    f = open("go.mod","a")
    f.write(gomod)
    f.close()
    os.system("go mod download {}@{}".format(url,pkg_ver))
    os.system("go get -u -v all {}@{}".format(url,pkg_ver))
    os.system("mkdir ./output")
    os.system("go list -json -m all | ./nancy sleuth -o csv > ./output/nancy-scan.csv")
    os.system("cp -r /go/pkg/mod ./output")
    pkg_name = url.split("/")[-1]
    os.system("zip -r /complete/{}.zip -r ./output".format(pkg_name))
else:
    print("./goscan.sh <url> <version>")
