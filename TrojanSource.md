# Trojan Source Vulnerability *Kezako* ?

## From the trick to the vulnerability
* The trick:

  **~>** **Unicode** := an encoding system that allows computers to exchange information regardless of the language used

  **~>** The **Unicode Bidirectional Algorithm** is used to define if we are writing Left-to-Right or Right-to-Left. 
  
  **~>**  For example: The `RLO` character (U+202e in unicode) is designed to support languages that are written right to left, such as Arabic and Hebrew. [Copy/paste](https://unicode-explorer.com/c/202E) it in your editor and you will see that you are now writing Right to Left
 

* Why it could be harmful?

  **~>** In fact it is not a new discovery, as these bugs/tricks have been known for like 20 years but it has to be widely and publicly reported. It is why a recent  research (2021) give the vulnerability a name, *Trojan Source*.
  
  **~>** An attacker could use this characters to produce source code whose tokens are logically encoded in a different order from the one  in  which  they  are  display. It is particularly dangerous, in a system based on human code review/validation, so basically all open source project. Indeed, a malicious contributor can add malicious code lines that will be seen as comments by reviewers, and thus divert theirs attentions.See [example](https://github.com/ariary/TrojanSourceFinder/blob/main/tests/comment-out.cpp) (from nickboucher repository)


## How to differently detect it?
* Github [print a warning](https://github.co/hiddenchars) when file contains unicode bidirection algorithms characters. *(does Gitlab detect it?)*
* [nickboucher/bidi-viewer](https://github.com/nickboucher/bidi-viewer): View and Analyze trojan source in a React-based webapp
* If your text editors have a “strip all comments” mode, enable it.

## Homoglyph

In the same category of risk we have homoglyph:

* A homoglyph is one of two or more graphemes, characters, or glyphs with shapes that appear identical or very similar.

  **~>** In fact we are speaking about homograph (word composed of homoglyph)

  **~>** For example `агіагу` is a homogram of `ariary`

  **~>** For an attacker who wants malicious code and lure the reviewers, homoglyph are a good weapon. For example he could write a function with a name which is an homogram of an already existing one and then call the homoglyph. In that way, we could think he is calling the original one.
See [example](https://github.com/ariary/TrojanSourceFinder/blob/main/tests/homoglyphe-function.go) (from nickboucher repository)

* Another example is the *IDN homograph attack*:

  **~>** The internationalized domain name (IDN) homograph attack is a technique used by malicious people to impersonate a website. With the internationalization of domain names, it is now possible to register an address with homoglyphs like ɡіτһυЬЬ.ϲоⅿ (and therefore create a similar or different website)
### Example
In python:
``` python
а = True # in a naive world, everybody is admin 

# some stuff 

a = False # finally the world is mean, we don't want admin anymore 

if а:
    print("You are an admin!")
else:
    print("You are not admin.")
```

The output:
```
You are an admin!
```

This trick is effecient to make believe the value of a variable changes. Especially for languages where assignment and initialization steps are identical. In this example we fake changing the value of `a`


## Sources

* [nickboucher/trojan-source repo](https://github.com/nickboucher/trojan-source): Expose vulnerability and many many examples
* [Whitepaper](https://trojansource.codes/trojan-source.pdf): present the vulnerability
* [article](https://certitude.consulting/blog/en/invisible-backdoor/): explanation & illustration of the vuln with javascript
* [dcode](https://www.dcode.fr/generateur-homoglyphes-homographes): make your own  homoglyph
