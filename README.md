# distill

A utility collection.

# Examples

Get the size of files and directories at a certain root directory:

```shell
mkdir -p ~/Documents/foo/bar
mkdir -p ~/Documents/foo/baz
echo 'something' > ~/Documents/foo/foo.txt
echo 'something else' > ~/Documents/foo/bar/bar.txt
echo 'something else entirely' > ~/Documents/foo/baz/baz.txt
go run . fs dirsize -startingdir ~/Documents/foo
Starting directory: /Users/jacob/Documents/foo

File                                Size in bytes
foo.txt                             10

Directory                           Size in bytes
bar                                 15
baz                                 24
```
