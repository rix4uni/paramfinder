# ParamFinder

ParamFinder crawl all input tags

## Installation
```
go install github.com/rix4uni/paramfinder@latest
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
http://testphp.vulnweb.com/guestbook.php
<input type="hidden" name="name" value="anonymous user">
<input type="submit" name="submit" value="add message">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">

http://testphp.vulnweb.com/AJAX/index.php

http://testphp.vulnweb.com/login.php
<input name="uname" type="text" size="20" style="width:120px;">
<input name="pass" type="password" size="20" style="width:120px;">
<input type="submit" value="login" style="width:75px;">
<input name="searchFor" type="text" size="10">
<input name="goButton" type="submit" value="go">
```
