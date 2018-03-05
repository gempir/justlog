#!/usr/bin/env python3

import datetime
import gzip
import os
import re
import shutil

mydate = datetime.datetime.now()
month = mydate.strftime('%B')
rootdir = '/var/twitch_logs'

for subdir, dirs, files in os.walk(rootdir):
    for file in files:
        if re.search(month, subdir, re.IGNORECASE):
            continue  # current month should be ignored
        if file == "channels":
            continue
        log = os.path.join(subdir, file)
        if re.search('.gz', file):
            continue  # already gzipped this
        with open(log, 'rb') as f_in, gzip.open(log + '.gz', 'wb') as f_out:
            shutil.copyfileobj(f_in, f_out)
        os.remove(log)