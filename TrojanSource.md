# Trojan Source Vulnerability *Kezako* ?

## From the trick to the vulnerability
* The trick:

  **~>** **Unicode** := an encoding system that allows computers to exchange information regardless of the language used

  **~>** The **Unicode Bidirectional Algorithm** is used to define if we are writing Left-to-Right or Right-to-Left. for example: The `RLO` character (U+202e in unicode) is designed to support languages that are written right to left, such as Arabic and Hebrew
 

* Why it could be harmful?

  **~>** In fact it is not a new discovery, as these bugs/tricks have been known for like 20 years but it has to be widely and publicly reported. It is why the research give the vulnerability a name.
  
  **~>** An attacker could use this characters to produce source code whose tokens are logically encoded in a different order from the one  in  which  they  are  display
...✏️

## How to detect it in a different way
* Github [print a warning](https://github.co/hiddenchars) when file contains unicode bidirection algorithms characters. *(does Gitlab detect it?)*
* [nickboucher/bidi-viewer](https://github.com/nickboucher/bidi-viewer): View and Analyze trojan source in a React-based webapp
* If your text editors have a “strip all comments” mode, enable it.

## Sources

* [nickboucher/trojan-source repo](https://github.com/nickboucher/trojan-source): Expose vulnerability and many many examples
* [Whitepaper](https://trojansource.codes/trojan-source.pdf): present the vulnerability
