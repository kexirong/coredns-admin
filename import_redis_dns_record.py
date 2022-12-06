import json
import urllib.request as request
from urllib.error import HTTPError
req1 = request.Request("http://192.168.2.223/api/v1/records",
                       headers={"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIiLCJ1c2VybmFtZSI6ImFkbWluIiwibmFtZSI6IiIsInBob25lIjoiIn0.wOoWObSeBogotdFvImgzNywGIQAAIy0ZpzSBIh5p_uY"},
                       method="GET"
                       )
with request.urlopen(req1) as f:
    data = json.loads(f.read()).get('data')
    for i in data:
        if i["type"] == "AAA":
            i["type"] = "AAAA"
        del i['key']

        req2 = request.Request("http://127.0.0.1:8088/api/v1/redis/record",
                               data=json.dumps(i).encode(), headers={'Content-Type': 'application/json'}, method="POST")
        try:
            with request.urlopen(req2) as f:
                pass
        except HTTPError as e:
            print(i)
            print(e.fp.read())

            break
