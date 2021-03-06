# translate_tool

# Description
Text translate tool

# How to use
Build after Set GOPATH=Current catalog

Trans is a text translate tool that can help you to extract all chinese from
file or directory. it can analyzes lua script, unity prefab and table file.
If there is more demand, you can easily add more file support. The first time
you run the programtrans, Automatically generate "config.ini" and "ignore.conf"
file. you can modify these files according to your requirements
```
Usage:
  trans [command]

Available Commands:
  getstring   Extract chinese characters
  transcoding file transcoding
  translate   Translation file or directory
  update      Update translation to dictionary
  version     View version

Flags:
  -h, --help   help for trans

Use "trans [command] --help" for more information about a command.
```

## SubCommand
getstring: Extract Chinese characters from a file or directory and save it to a text file
```
Usage:
    trans getstring [flags]

Flags:
    -d, --db string    Translation data dictionary (default "dictionary.txt")
    -s, --src string   The extracted file or directory path
```

update:	Update translation to dictionary
```
Usage:
  trans update [flags]

Flags:
  -d, --db string       Translation data dictionary (default "dictionary.txt")
  -u, --update string   The new translation data (default "chinese.txt")
```

translate: Translation using dictionary file or directory. If the output does not exist will be created automatically
```
 Usage:
    trans translate [flags]

 Flags:
    -d, --db string       Translation data dictionary (default "dictionary.txt")
    -o, --output string   The output file or directory path translated
    -r, --routine int     Goroutine number. This is a test parameters (default 1)
    -s, --src string      Translated file or directory path
```

transcoding: file transcoding
```
file transcoding Support utf8, gbk, hz-gb2312, gb18030, big5, euc-jp, iso-2022-jp, shift_jis, euc-kr.
Notice: This tool can only transcoding, and can not be translated.

Usage:
  trans transcoding [flags]

Flags:
  -d, --decoding string   encoding of input file
  -e, --encoding string   encoding of output file
  -i, --input string      Input file or directory
  -o, --output string     Output file or directory
```

License
-------------

The MIT License (MIT)

Copyright (c) 2016 liubo5

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
