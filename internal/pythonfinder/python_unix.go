package pythonfinder

import "regexp"

var pythonFileRegex = regexp.MustCompile(`^python(\d(\.\d\d?)?)?$`)
