# TODO

- Upgrade Python version for a given virtual environment.
  - Upgrade all virtualenvs to a given Python version

- Ability to list all the available Python versions the tool can detect.
  - Output it when the provided version does not exists?
  - A flag to list it out? (`pyvenv --execs` | `pyvenv --pythons`)
  - Maybe output as per section where each section is a single provider
- A flag to provide provider(s) to prioritize them when finding the Python
  versions. If multiple providers are given, the order in which they're given
  is the priority.

- Is it possible to create a subshell with the environment activated similar to
  `pipenv` in golang? If so, allow that with an `activate` command.
- Then, fuzzy find all the managed venvs and allowing to create a subshell for
  the selected environment, `cd` into the project and activate it. Pressing
  <kbd>CTRL-D</kbd> will exit the subshell and the user will be back to the
  original shell.
Ref: https://github.com/gtalarico/pipenv-pipes

- Environment variable `PYVENV_HOME` pointing to the data directory where all
  the environments are stored, similar to `PIPENV_HOME`.

- Ability to link a different project to a virtual environment. By default,
  the project used to create the virtualenv will be linked, but the user can
  change this. Maybe a `--link` option?

## New providers

- Windows registry
- MacOS (Python.Framework)
- Homebrew?
- `asdf` provider - https://github.com/asdf-vm/asdf
