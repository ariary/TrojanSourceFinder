<h1 align="center">TrojanSourceFinder</h1>
<h4 align="center">TrojanSourceFinder helps developers detect "Trojan Source" vulnerability in source code.</h4>
<p align="center">
  Trojan Source vulnerability allows an attacker to make malicious code appear innocent.
  In general, the attacker tries to lure by passing his code off as a comment (visually). It is a serious threat because it concerns many languages. Projects with multiple "untrusted" sources could be concerned
  <br><br>
  <strong>
    <a href="https://github.com/ariary/TrojanSourceFinder#detect-trojan-source">Detect evil 🔎</a>
    ·
    <a href="https://github.com/ariary/TrojanSourceFinder#visualize-trojan-source">Track evil 👀</a>
    ·
    <a href="https://github.com/ariary/TrojanSourceFinder/blob/main/TrojanSource.md">Trojan Source ❓</a>
  </strong>
</p>

## Install

```shell
go install github.com/ariary/TrojanSourceFinder@latest
```

Or
```shell
git clone https://github.com/ariary/TrojanSourceFinder
cd TrojanSourceFinder
make before.build
make build.tsfinder
```

If the command `make build.tsfinder` failed, try:
```shell
env GOOS=target-OS GOARCH=target-architecture
go build -o tsfinder cmd/main.go
```

## Detect Trojan Source
*> Help the detection of Trojan source for manual code review or with CI/CD pipelines*

To detect Trojan source in file *\<filename\>*:
```shell
tsfinder <filename>
```

You could also scan a whole directory using `-r` flag:
```shell
tsfinder -r <directory>
```

## Visualize Trojan Source
*> Visualize how the code is really interpreted by machines/compiler*

*tsfinder* is deliberately not very verbose. By default, it will only output if Trojan Source code has been detected. To have more verbosity and **visualize the dangerous line add the flag `-v`**

## Demo
![demo](https://github.com/ariary/TrojanSourceFinder/blob/main/img/tsfinder.gif)
