epoch-cat
=========

epoch-cat is a filter command that converts the UNIX time appearing in some log files to the human-readable format only for that part

## Installation

```console
$ go get github.com/b4b4r07/epoch-cat
```

or

```zsh
zplug "b4b4r07/epoch-cat", as:command, rename-to:ecat
```

## Usage

```console
$ ./epoch-cat -h
Usage of ./epoch-cat:
  -f string
        Specify the time format for convert
  -p string
        Add the prefix for searching to prevent mismatch when looking for UNIX time
  -q    Quote the date to be output
```

### Examples

```console
$ ./epoch-cat -q access.log | jq . | head
{
  "time": "2016-05-25T14:35:35Z",
  "uuid": "3bf99536760ae910dd77a0fb9cc493524dc5b808",
  "level": "WARNING",
  "message": {
    "message": "SQLSTATE[1234] Unknown database 'hyper.services',
    "_exception": {
        ...
        ...

$ echo 1476707810.8115 | ./epoch-cat
2016-10-17T12:36:50Z

$ echo 1476707810.8115 | TZ=JST ./epoch-cat
2016-10-17T21:36:50+09:00

$ echo 1476707810.8115 | ./epoch-cat -f "%a, %d %b %Y %H:%M:%S +0000"
Mon, 17 Oct 2016 12:36:50 +0000
```

About time format, see also: <https://docs.python.org/2/library/time.html#time.strftime>

## Author

b4b4r07

## License

MIT
