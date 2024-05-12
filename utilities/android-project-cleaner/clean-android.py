#!/usr/bin/python3
import os
import shutil 

bad_exts = [".sh",".bat",".bin",".class",".jar",".apk"]
def rm(path):
    print("removing " + path)
    if os.path.isfile(path):
        try:
          os.remove(path)
        except:
          pass
    else:
        try:
          shutil.rmtree(path)
        except:
          pass

def check_ext(filename):
    for ext in bad_exts:
        if filename.endswith(ext):
            return True
    return False


current_directory = os.getcwd()
all_files = os.listdir(current_directory)
zip_files = [file for file in all_files if file.endswith('.zip')]
for zip_file in zip_files:
    os.system("unzip " + zip_file)
    extracted_folder = "./"+zip_file.replace(".zip","")
    rm(extracted_folder+"/.gradle")
    rm(extracted_folder+"/build")
    
    for dirpath, dirs, files in os.walk(extracted_folder):
        for file in files:
            if check_ext(file):
                rm(dirpath +  "/" + file)

