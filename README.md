# ParamFinder

ParamFinder crawl all input tags

## Installation
```
go install github.com/rix4uni/paramfinder@latest
```

## Usage
```
Usage of paramfinder:
  -c int
        number of concurrent goroutines (default 20)
  -k    allow insecure server connections when using SSL
  -o string
        output file path
  -timeout int
        HTTP request timeout duration (in seconds) (default 30)
  -v    enable verbose mode
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
URL: http://testphp.vulnweb.com/guestbook.php
<input type="hidden" name="name" value="anonymous user">
<input type="submit" name="submit" value="add message">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">

URL: http://testphp.vulnweb.com/login.php
<input name="uname" type="text" size="20" style="width:120px;">
<input name="pass" type="password" size="20" style="width:120px;">
<input type="submit" value="login" style="width:75px;">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">
```
