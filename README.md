# toSarif

Utility for converting code quality JSON output to SARIF format.

Basically it seems like GitHub is all-in on SARIF as the format for code quality
data but the tools out there for getting your output from one tool or another to
SARIF are still all over the place. So I guess [what's one more]?

## Usage

Download a release for your system and architecture from the [releases page] and
you're off:

```bash
tosarif -h

Usage: tosarif [-h] [-f value] [-i value] [-o value] [parameters ...]
 -f, --format=value
                    Source format. Options: pyright, php-cs-fixer
 -h, --help         Display this help message
 -i, --input=value  Input file
 -o, --output=value
                    Output file
```

[what's one more]: https://xkcd.com/927/
[releases page]: https://github.com/bensaufley/toSarif/releases
