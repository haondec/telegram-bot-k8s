import subprocess
import json

batcmd="./telegram-dev"
# result = subprocess.check_output(batcmd, shell=True)
log = open("telegram-dev.log", "a")
subprocess.Popen(batcmd, stdout=log, stderr=log, shell=True)
