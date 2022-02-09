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
### With `go`

*> Via `go install`*
```shell
go install github.com/ariary/TrojanSourceFinder/cmd/tsfinder@latest
```
Make sure `$GOPATH` is in your `$PATH`

*> From source*
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

### With `curl`
*> From release*

```shell
curl -lO -L https://github.com/ariary/TrojanSourceFinder/releases/latest/download/tsfinder && chmod +x tsfinder
```

## Detect Trojan Source
*> Help the detection of Trojan source for manual code review or with CI/CD pipelines (Unicode bidirectional characaters)*

To detect Trojan source in file or directory *\<path\>*:
```shell
tsfinder [path]
```

### Detect only in text file
*> Source code files are likely text files. Withdraw them for scan could help to rule out false positives*

```shell
tsfinder -t [path]
```
Add `-v` help to see which file has been skipped by scan.

### Go further *(Homoglyph)*

Trojan Source is not new and isn't the only hazard. Another one is *"Homoglyph"*.(*[Kezako?](https://github.com/ariary/TrojanSourceFinder/blob/main/TrojanSource.md#homoglyph)*)

tsfinder help detecting them with `homoglyph` command:
```shell
tsfinder homoglyph [filename] [flags]
```

You could see if there is a sibling (ie word with same "skeleton") for the homographs found in `path` using the flag `--sibling`:
```shell
tsfinder homoglyph [filename] --sibling [path] 
```
*Functionality under development, mainly depending on other project*

## Visualize Trojan Source
*> Visualize how the code is really interpreted by machines/compiler*

*tsfinder* is deliberately not very verbose. By default, it will only output if Trojan Source code has been detected. To have more verbosity and **visualize the dangerous line add the flag `-v`**.

To better see where Trojan Sources were, you could enable colored output with `-c` flag (also useful with directory scan):
```shell
tsfinder -c -v <directory>
```

## Demo

![demo](https://github.com/ariary/TrojanSourceFinder/blob/main/img/tsfinder-demo-trojansource.gif)

### Homoglyph

![demo](https://github.com/ariary/TrojanSourceFinder/blob/main/img/tsfinder-demo-homoglyph.gif)

## Alternative

As mentioned by `@ioah86` [here](https://www.reddit.com/r/cybersecurity/comments/qlh5j9/my_take_on_trojan_source/), trojan source could also been detected w/ a one liner using grep.

The big difference is the output format and the exit status code (`tsfinder` exit with status code `0` if no Trojan source has been found, `1` otherwise; the opposite  for `grep`)

Also, this one-liner does not resolve the homoglyph issue

| Goal   |      `tsfinder`     |  `grep` one-liner |
|:----------:|:-------------|:------|
| Scan all files + show lines|  `tsfinder -v .` | `grep -arE $'(\u2066\|\u2067\|\u2068\|\u202A\|\u202B\|\u202D\|\u202E\|\u202C\|\u2069\|\u200E\|\u200F\|\u061C\|\u2066\|\u2067\|\u2068)'` |
| Scan only on human-readable files| `tsfinder -t .` | `grep -IrE $'(\u2066\|\u2067\|\u2068\|\u202A\|\u202B\|\u202D\|\u202E\|\u202C\|\u2069\|\u200E\|\u200F\|\u061C\|\u2066\|\u2067\|\u2068)'`|
|Exit with status code 1 if found|default|`[one-liner] && exit 1 \|\| exit 0`|
