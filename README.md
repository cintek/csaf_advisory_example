# csaf_advisory_example

This is an example for how to use the [advisory model](https://github.com/cintek/csaf_distribution/blob/main/csaf/advisory.go) to change an advisory of the CSAF 2.0 standard. To be specific it changes the category of every branch to "legacy".

## How to use

Build the binary:
`go build`

Run the binary with one or more files as parameters:
`csaf_advisory_example [files]`
