### Simple Demo

1) Start by taking a look at the program in `main.go`. All it does is read a variable from the runtime environment and print its value:

```
$ cat main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Getenv("SUPER_SECRET_VARIABLE"))
}
```


2) Create a padl project to manage secrets for the program:

```
$ padl project create --name padl-demo --description "throwaway demo project"
project padl-demo initialized successfully!
```

3) Take a look at the padlfile, see that it has the project's shared key ID, as well the creator user's key ID:

```
$ cat .padlfile.yaml
data:
  project_id: padl-demo
  variables: {}
  keys:
  - e44148a0d66af9564d3f439b2cda1b81
  shared_key: b50656b2f58496b282972b43af613407
```


4) Add a secret to the padlfile with the name as expected in the program:

```
$ padl file secret set --name SUPER_SECRET_VARIABLE --secret "this is going to be encrypted shortly"
padlfile updated!
```

5) See that your padlfile now has an encrypted secret in it:

```
$ cat .padlfile.yaml
data:
  project_id: padl-demo
  variables:
    SUPER_SECRET_VARIABLE: |
      -----BEGIN PADL ENCRYPTED SECRET-----
      ZTQ0MTQ4YTBkNjZhZjk1NjRkM2Y0MzliMmNkYTFiODEoZWR6T2IvbGNobTdZSDZP
      dFNWQnNlaUpobnpwRUlPNnE1ZWwwTkJDdE5Bc2d5VzhRa1BYT3U5NXFxNS9ibnZQ
      YUt0SVYwMm5HZUc2Q0ZzeDUvaDJjWkZLRHZ3dUZ2QVFjeXpaM2JzSE41ZVdvRkVk
      QXlwd1V4RzRpaUJLQWsvanQ0cFVXQVhXdk5hQlhSU3ptZDIrckNPLzlvb2xmaVI4
      cHBOVVc5T2tGV2VmZ09QNjdIb0JWU09INWxhT0oxbDJtZ2Qzc2F2b0hwdzZPMmZ1
      TGEzM2lBbDljdW8rb0pid0FUa2kyeStEY0ZQa1dwd3dST2FFWG05MU9kWlZtakYr
      MXdrV2lLcHEyZmEySTdjdUQ0VDkrSjQ0VE9XeHdtajFOZC9kdWRYVFRmVkp6MEE1
      c3phZnJMUVRJTi9yUSt0SWFDaWN2eE9EOHBZN0ZXS0crK3VlcEUzUlB0cTVjdGJT
      eVl2NG1jSUU2WDR0dXJzWjJjM3Y2SFhvaGovZXErdjRjZlNrbnNJMHdZTGxZbEpR
      SHUrcS9JbUN1V2NVRUpWcjlrWURxb3crQWcwUGhCVzQwN3l2aEE1VmtNUnFyMHZz
      VW9GOHBSSVZ5dnk3YStLQmpOd3l5d3A5L0NpVkdwNUljNmtMT2tZanhkNGZlNi9t
      ejZyUGt0RC9QUVk2SFpPUGcyMjJZZnEzdkZnd0N0WVlpeGs1d0ZBN2I1eFZIbDhY
      ckJzcXJ6TUtUcFJMUHpURUNaRWNLSmJDVUpUZ1NHWXVJWGQvL2dXYWF4Z1ArSTJK
      VnZDRXF1UVU2ZTZKQVJkeGkwd0JBYlBsVXZ2cXRwemdmd0hoWnlIOWxET3JVM1Mr
      ODBTOEFYTFQ5VW1PRDZyVk5vTWg0clB5b0ZnQU1ITy8yQ296R1VSZWVObVk9KQpi
      NTA2NTZiMmY1ODQ5NmIyODI5NzJiNDNhZjYxMzQwNyhMbDVCYVBBdFgzYU4xTXFC
      RGVZK3k4bWlGWnRuY2xYYlJDKy92eVJId3ZZY2VkUEhpbHdaRnNqNVBWYkN5YnVk
      NmhqMVlIcWJSMGVzbXVqSXBzUVp1b0VXNEp4YVpYcFZrREdEYUphb1M2K3ViM1FP
      cmswVFRLSEpJS1dTdnBVeEl5VHR1TEFnc1R6ZnJtSVp2M0tWVUV6dzV2T1ZYNk1S
      SW1BamNxbU9XcHc0bGZVei8zY0wvTUVsQUdQTmdCNzV3NUZnK2FMS2c5QzFuNnow
      U2NSc05qYUF1MVk1NmQ5OHJTZFJnSjMzbkd1TFJNaDZxYitITzJNSEQraUgxVEVF
      d0FpUFBZYUZqRTh4OE9KVzBHYUlvbk1sMTZWei9uQkk2NytkNVcyQmhjQzdqalpQ
      YzAxS3pDbnhJVG9sbm91Wld3UzdaaXBkNnB1a0FsOW1PVThOU0E9PSk=
      -----END PADL ENCRYPTED SECRET-----
  keys:
  - e44148a0d66af9564d3f439b2cda1b81
  shared_key: b50656b2f58496b282972b43af613407
```

6) Run the program with the `padl run` wrapper which decrypts padlfile secrets and loads them onto the environment of successive commands:

```
$ padl run go run main.go
this is going to be encrypted shortly
``` 

#### Success!