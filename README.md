# ParamFinder

ParamFinder crawl all input tags

## Installation
```
go install github.com/rix4uni/paramfinder@latest
```

## Usage
```
Usage of paramfinder:
  -ao string
        File to append the output instead of overwriting.
  -c int
        number of concurrent goroutines (default 50)
  -insecure
        allow insecure server connections when using SSL
  -o string
        output file path
  -silent
        silent mode.
  -timeout int
        HTTP request timeout duration (in seconds) (default 10)
  -turl
        transform URL with extracted parameters
  -verbose
        enable verbose mode
  -version
        Print the version of the tool and exit.
```

## Example usages

Single URL:
```
echo "https://testphp.vulnweb.com" | paramfinder
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

output:
```
▶ cat urls.txt | go run paramfinder.go -silent -turl
URL: http://testphp.vulnweb.com/guestbook.php
<input type="hidden" name="name" value="anonymous user">
<textarea name="text" rows="5" wrap="VIRTUAL" style="width:500px;">
<input type="submit" name="submit" value="add message">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">

TRANSFORM_URL: http://testphp.vulnweb.com/guestbook.php?name=cqzsupn&text=dspeicx&submit=ldqmjfc&searchFor=hbtgwmk&goButton=eigygqa
URL: http://testphp.vulnweb.com/login.php
<input name="uname" type="text" size="20" style="width:120px;">
<input name="pass" type="password" size="20" style="width:120px;">
<input type="submit" value="login" style="width:75px;">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">

TRANSFORM_URL: http://testphp.vulnweb.com/login.php?uname=imhagwh&pass=heijicp&searchFor=yncbvdr&goButton=rzqqczf
```
