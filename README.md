## ParamFinder

ParamFinder crawls all input and textarea tags

## Installation
```
go install github.com/rix4uni/paramfinder@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/paramfinder/releases/download/v0.0.3/paramfinder-linux-amd64-0.0.3.tgz
tar -xvzf paramfinder-linux-amd64-0.0.3.tgz
rm -rf paramfinder-linux-amd64-0.0.3.tgz
mv paramfinder ~/go/bin/paramfinder
```
Or download [binary release](https://github.com/rix4uni/paramfinder/releases) for your platform.

## Compile from source
```
git clone --depth 1 https://github.com/rix4uni/paramfinder.git
cd paramfinder; go install
```

## Usage
```yaml
Usage of paramfinder:
      --concurrency int   number of concurrent goroutines (default 50)
      --output string     output file path
      --silent            silent mode.
      --timeout int       HTTP request timeout duration (in seconds) (default 30)
      --verbose           enable verbose mode
      --version           Print the version of the tool and exit.
```

**Note:** Insecure SSL connections are automatically enabled. The tool outputs only the transformed URL with all parameters set to `rix4uni`.

## Example usages

Single URL:
```yaml
echo "http://testphp.vulnweb.com/login.php" | paramfinder
```

Multiple URLs:
```yaml
cat urls.txt | paramfinder
```

urls.txt contains:
```yaml
http://testphp.vulnweb.com/login.php
http://testphp.vulnweb.com/guestbook.php
http://testphp.vulnweb.com/AJAX/index.php
```

Output:
```yaml
▶ cat urls.txt | paramfinder --silent
http://testphp.vulnweb.com/login.php?uname=rix4uni&pass=rix4uni&searchFor=rix4uni&goButton=rix4uni
http://testphp.vulnweb.com/guestbook.php?name=rix4uni&text=rix4uni&submit=rix4uni&searchFor=rix4uni&goButton=rix4uni
```

## Real world Example why this tool is usefull
```yaml
echo "https://domain.com/xyz/index.php" | paramfinder --silent
https://domain.com/xyz/index.php?view=rix4uni
```
## Found xss in `view` parameter
- https://domain.com/xyz/index.php?view=1'-confirm`K`-'=1
