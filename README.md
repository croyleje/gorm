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
  * Prompt when restoring links.
  * Add link icon to tui display.

## TODO / BUGS
* BUG: FilterState bug PATCHED
* TODO: add error handling to cmd package
* TODO: rewrite variables to constants where possible
* TODO: rename variables and funcs

vim: tw=80
