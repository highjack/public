import os
import sys
import re
import subprocess

if len(sys.argv) < 2:
    print("Usage: python pyscan.py <url-list.txt>")
    sys.exit(1)

urlList = sys.argv[1]
fp = open(urlList, "r")
urls= fp.readlines()
filenames = []
requirements_files = []
print("[+] Downloading packages....")
for url in urls:
    filename = os.path.basename(url)
    filename = filename.replace("\n","")
    filenames.append(filename)
    split_filename = filename.split("-")
    package_name = split_filename[0]
    version = split_filename[1]
    platform = split_filename[-1].replace(".whl","")
    if "." in platform:
        platform = platform.split(".")[0]
    python_version = ""
    python_version_pattern = r'cp(\d+)'
    python_version_matches = re.search(python_version_pattern, filename)
    if python_version_matches:
        python_version = python_version_matches.group(1)
    abi = ""
    abi_pattern = r'abi(\d+)'
    abi_matches = re.search(abi_pattern, filename)
    if abi_matches:
        abi = abi_matches.group(0)
    pip_download = "pip download --only-binary=:all: "
    if abi != "":
        pip_download = pip_download + "--abi \"{}\" ".format(abi)
    if platform != "any":
        pip_download = pip_download + "--platform {} ".format(platform)
       
    if python_version != "":
        pip_download = pip_download +  "--python-version={} ".format(python_version)
    python_package_fullname = "{}=={}".format(package_name, version)
    pip_download = pip_download + python_package_fullname +  " -d ./whl"
    print(pip_download)
    os.system(pip_download)
    requirements_files.append(python_package_fullname)
requirements_files = list(set(requirements_files))

try:
    os.system("mkdir ./vulns")
except:
    pass

for requirement in requirements_files:
    print("checking {} for vulns...".format(requirement))
    fp = open("requirements.txt","w")
    fp.write(requirement)
    result = subprocess.run(['pip-audit', '-r', 'requirements.txt'], capture_output=True, text=True)
    if "No known vulnerabilities found" not in result.stderr and "No known vulnerabilities found" not in result.stdout:
        with open("./vulns/{}.txt".format(requirement.replace("==","-")), "w") as f:
            f.write(result.stdout)   
