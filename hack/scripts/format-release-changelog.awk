# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Transform a GoReleaser GitHub release body into docs-friendly changelog markdown.
# Usage: awk -v version=v1.0.0 -f format-release-changelog.awk changelog-raw.md

BEGIN {
	print "## " version
	print ""
	section_count = 0
}

/^## Changelog[[:space:]]*$/ {
	next
}

/^\* [0-9a-f]{40}[[:space:]]/ {
	line = $0
	sub(/^\* [0-9a-f]{40}[[:space:]]+/, "", line)
	sub(/^[a-z]+(\([^)]*\))?!?:[[:space:]]*/, "", line)
	if (length(line) > 0) {
		print "* " toupper(substr(line, 1, 1)) substr(line, 2)
	}
	next
}

/^### / {
	if (section_count++ > 0) {
		print ""
	}
	print
	next
}

{
	print
}
