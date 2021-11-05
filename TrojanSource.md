# Trojan Source Vulnerability *Kezako* ?

## From the trick to the vulnerability
* The trick:

  **~>** **Unicode** := an encoding system that allows computers to exchange information regardless of the language used

  **~>** The **Unicode Bidirectional Algorithm** is used to define if we are writing Left-to-Right or Right-to-Left. 
  
  **~>**  For example: The `RLO` character (U+202e in unicode) is designed to support languages that are written right to left, such as Arabic and Hebrew. [Copy/paste](https://unicode-explorer.com/c/202E) it in your editor and you will see that you are now writing Right to Left
 

* Why it could be harmful?

  **~>** In fact it is not a new discovery, as these bugs/tricks have been known for like 20 years but it has to be widely and publicly reported. It is why a recent  research (2021) give the vulnerability a name, *Trojan Source*.
  
  **~>** An attacker could use this characters to produce source code whose tokens are logically encoded in a different order from the one  in  which  they  are  display. It is particularly, in a system based on human code review/validation, so basically all open source project. Indeed, a malicious contributor can add malicious code lines that will be seen as comments by reviewers, and thus divert theirs attentions.

## How to detect it in a different way
* Github [print a warning](https://github.co/hiddenchars) when file contains unicode bidirection algorithms characters. *(does Gitlab detect it?)*
* [nickboucher/bidi-viewer](https://github.com/nickboucher/bidi-viewer): View and Analyze trojan source in a React-based webapp
* If your text editors have a “strip all comments” mode, enable it.

## Sources

* [nickboucher/trojan-source repo](https://github.com/nickboucher/trojan-source): Expose vulnerability and many many examples
* [Whitepaper](https://trojansource.codes/trojan-source.pdf): present the vulnerability
