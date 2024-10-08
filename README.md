## ParamFinder

ParamFinder crawls all input and textarea tags

## Installation
```
go install github.com/rix4uni/paramfinder@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/paramfinder/releases/download/v0.0.2/paramfinder-linux-amd64-0.0.2.tgz
tar -xvzf paramfinder-linux-amd64-0.0.2.tgz
rm -rf paramfinder-linux-amd64-0.0.2.tgz
mv paramfinder ~/go/bin/paramfinder
```
Or download [binary release](https://github.com/rix4uni/paramfinder/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/paramfinder.git
cd paramfinder; go install
```

## Usage
```console
Usage of paramfinder:
  -a, --append string     File to append the output instead of overwriting.
  -c, --concurrency int   number of concurrent goroutines (default 50)
  -i, --insecure          allow insecure server connections when using SSL
  -n, --no-turl           Do not print transform URL with extracted parameters
      --only-hidden       print only hidden input tags
  -o, --output string     output file path
  -s, --silent            silent mode.
  -t, --timeout int       HTTP request timeout duration (in seconds) (default 10)
  -v, --verbose           enable verbose mode
  -V, --version           Print the version of the tool and exit.
```

## Example usages

Single URL:
```
echo "http://testphp.vulnweb.com/login.php" | paramfinder
```

Multiple URLs:
```
cat urls.txt | paramfinder
```

urls.txt contains:
```
http://testphp.vulnweb.com/login.php
http://testphp.vulnweb.com/guestbook.php
http://testphp.vulnweb.com/AJAX/index.php
```

Output:
```
▶ cat urls.txt | go run paramfinder.go --silent
URL: http://testphp.vulnweb.com/login.php
<input name="uname" type="text" size="20" style="width:120px;">
<input name="pass" type="password" size="20" style="width:120px;">
<input type="submit" value="login" style="width:75px;">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">

TRANSFORM_URL: http://testphp.vulnweb.com/login.php?uname=xeelxcp&pass=xvmlmom&searchFor=rnzfvyo&goButton=topgczf
URL: http://testphp.vulnweb.com/guestbook.php
<input type="hidden" name="name" value="anonymous user">
<textarea name="text" rows="5" wrap="VIRTUAL" style="width:500px;">
<input type="submit" name="submit" value="add message">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">

TRANSFORM_URL: http://testphp.vulnweb.com/guestbook.php?name=dqfgcnd&text=uaxtznk&submit=rzbzvep&searchFor=zsifjbr&goButton=zcduicb
```

## Real world Example why this tool is usefull
```
echo "https://https://domain.com/xyz/index.php" | go run paramfinder.go --silent --only-hidden
URL: https://https://domain.com/xyz/index.php
<input type="hidden" name="view" value="">

TRANSFORM_URL: https://https://domain.com/xyz/index.php?view=nqkfbwf
```
## Found xss in `view` parameter
- https://https://domain.com/xyz/index.php?view=1'-confirm`K`-'=1