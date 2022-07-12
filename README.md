# gorm
Go written replacement for the Linux/Unix ***rm*** command.  Gorm implements the
[FreeDesktop.org](https://specifications.freedesktop.org/trash-spec/trashspec-latest.html)
Trash Specification v1.0 (Latest Release) along with some features I found
lacking in other cli/tui trash applications.

## Features
* Automatically `shreds` sensitive files and user definable files by extension
  and name, including.*
  * .key
  * .kbx
  * .gpg
  * id_rsa
  * id_ecdsa
  * id_ecdsa_sk
  * id_ed25519
  * id_ed25519_sk
* Fuzzy search by filename and path.*
* Local or original path restoration.
* Detail view including.
  * Directory contents
  * Type
  * Name
  * Full path
  * User, group name and ID
  * Permissions
  * Deletion date
  * Last modified date
  * Last status change date
  * File creation date

\* See TODO/BUGS and Unimplemented features list.

## Unimplemented features
* Prompt and auto generate new mount point trash files.
* Link handling
  * Prompt when trashing links (symbolic & hard).
  * Prompt when restoring links and confirm link is still valid.

## TODO / BUGS
* BUG: FilterState bug PATCHED
* TODO: add error handling to cmd package
* TODO: rewrite variables to constants where possible
* TODO: rename variables and funcs
* TODO: rewrite arg parsing to use the Kingpin package enabling default width
  and height flags.

## Status
gorm is in a alpha status it is useable for home directory trash and single user
systems.  I have been adding features almost daily but I do welcome anyone to
try it and give me there feedback.  Pull requests are welcome but I am still
merging code from the original application so if you submit a pull request it
will probably have to be changed as the other features are implemented.  If you
are interested in helping with the rewrite from C to Go feel free to contact me
and I can get you access to the other repositories so you can help with the
rewrite.

The tui interface is basically set but there is plenty of formatting and default
theme changes that need to be merged.  I have decided to drop the icon support
and all the icons have been replaced with single word descriptions *ie.* **File**,
**Directory** or **Link**.  If there is interest in reimplementing the icons, we
could see about adding a config option or flag to enable them.

vim: tw=80
