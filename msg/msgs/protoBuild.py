import subprocess
import os
import os.path

def build(SRC_DIR,DST_DIR,file,lang):
    if not os.path.exists(DST_DIR): os.makedirs(DST_DIR)
    cmd="protoc -I="+SRC_DIR+" --"+lang+"_out="+DST_DIR+" "+SRC_DIR+"/"+file+".proto"
    print cmd
    subprocess.call(cmd, shell=True)

build(".",".","msg","go")