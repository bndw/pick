# Contributing to pick

We welcome contributions to pick of any kind including documentation, 
organization, tutorials, bug reports, issues, feature requests,
feature implementations, pull requests, etc.

## Table of Contents

* [Reporting Issues](#reporting-issues)
* [Submitting Patches](#submitting-patches)
* [Code Contribution Guidelines](#code-contribution-guidelines)

## Reporting Issues

If you believe you have found a bug in pick or its documentation, use
the Github [issue tracker](https://github.com/bndw/pick/issues) to report the problem to the pick maintainers.
When reporting the issue, please provide the version of pick in use (`pick version`) and your operating system.

## Submitting Patches

The pick project welcomes all contributors and contributions regardless of skill or experience level.
If you are interested in helping with the project, we will help you with your contribution.
Because we want to create the best possible product for our users and the best contribution experience for our developers,
we have a set of guidelines which ensure that all contributions are acceptable.
The guidelines are not intended as a filter or barrier to participation.

### Code Contribution Guidelines

To make the contribution process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes.  We encourage pull requests to allow for review and discussion of code changes.
* When you’re ready to create a pull request, be sure to:
    * Have test cases for the new code. If you have questions about how to do this, please ask in your pull request.
    * Add documentation if you are adding new features or changing functionality.  The docs site lives in `/docs`.
    * Squash your commits into a single commit. `git rebase -i`. It’s okay to force update your pull request with `git push -f`.
    * Make sure `make test` passes, and `go build` completes. [Travis CI](https://travis-ci.org/bndw/pick) (Linux and OS&nbsp;X) will catch most of this.
