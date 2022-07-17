# gorm
Go written replacement for the Linux/Unix ***rm*** command.  Gorm implements the
[FreeDesktop.org](https://specifications.freedesktop.org/trash-spec/trashspec-latest.html)
Trash Specification v1.0 (Latest Release) along with some features I found
lacking in other cli/tui trash applications.

## Commands
* **-put** moves a file to that mount points trash file.
* **-list** cli list in unicode table format of all trash files with key stats.
* ***No flag*** starts gorm in tui mode allowing detail view of trashed files,
  restoring files or emptying trash folder.
* **-empty** empties trash folder.
* **-debug** enables debug log.

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

## Installation
gorm currently supports local builds only, after the first release a binary and
Arch AUR package will be available.  To build the binary use the following
steps.

Clone the repo.
`git clone https://github.com/croyleje/gorm.git`

Build the package.
`go build -o gorm main.go`

Trash folders are created at the root of each mount point only when needed.  The
only exception to this is the user home directory trash files which are stored
in `~/.local/share/Trash/`.

## Unimplemented features
* View multiple mount point trash files.

## TODO / BUGS
* TODO: rewrite arg parsing to use the Kingpin package enabling default width
  and height flags.
* BUG: interface conversion: list.Item is nil, not ui.item causing Panic

## Status
gorm is now in a beta testing phase most of the core features are implemented
there are still a few yet to be merged but it is now useable on multi user
systems and provides all the core features ie. trashing, restoration (working
directory and original path) and detail view of trashed items.

vim: tw=80
