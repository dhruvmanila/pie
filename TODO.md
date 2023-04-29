# TODO

- Upgrade Python version for a given virtual environment. Upgrade all
  virtualenvs to a given Python version. This might be difficult to achieve.

- Ability to list all the available Python versions the tool can detect.
  - Output it when the provided version does not exists?
  - A flag to list it out? (`pyvenv --execs` | `pyvenv --pythons`)
  - Maybe output as per section where each section is a single provider
  - What about symlinks? Do they need to be resolved before displaying?
    If not, then there will be entries pointing to the same executable albeit
    with a different path. If we do resolve, then there will be duplicate
    entries.

- A flag to provide provider(s) to prioritize them when finding the Python
  versions. If multiple providers are given, the order in which they're given
  is the priority.
  ***Very low priority, not required***

- Is it possible to create a subshell with the environment activated similar to
  `pipenv` in golang? If so, allow that with an `activate` command.

  Then, fuzzy find all the managed venvs and allowing to create a subshell for
  the selected environment, `cd` into the project and activate it. Pressing
  <kbd>CTRL-D</kbd> will exit the subshell and the user will be back to the
  original shell.

  _See: https://github.com/gtalarico/pipenv-pipes_

- Ability to link a different project to a virtual environment. By default,
  the project used to create the virtualenv will be linked, but the user can
  change this. Maybe a `link` subcommand?
  ***Low priority, probably not required***

## New providers

- Windows registry

## NOTES

### Upgrade Python versions

- _Created a venv with Python 3.10.2_
- Remove all the base executables from the bin directory: `pip`, `pip<major>`,
  `pip<major>.<minor>`, `python`, `python<major>`, `python<major>.<minor>`
- Run the command `python -m venv venv --upgrade` with Python 3.10.5

Here, as it's just a patch version change, nothing else needs to be done. If
the major/minor version changes (`3.10` -> `3.11`), then the
`lib/python3.10/site-packages` needs to be removed backing up all the
dependencies which were installed and then reinstalling them with the new
pip.

***Do we even want to allow a major/minor version bump in a virtual
environment?***
