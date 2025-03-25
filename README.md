# reflector-go

Mirrorlist updater inspired by [reflector](https://xyne.dev/projects/reflector/).
Only uses mirrors with 100% completion and sorts them by score.

This is mainly meant for personal use as I find reflector to be unnecessarily slow.

## Usage

```bash
Usage of reflector-go:
  -country string
        Comma separated list of countries (default "Finland,Sweden,Estonia,Norway,Denmark,Latvia,Lithuania,Poland,Germany,France")
  -limit int
        Maximum number of mirrorlists to use (default 20)
  -protocol string
        Comma separated list of protocols (default "https")
  -save string
        Save to file (default "/etc/pacman.d/mirrorlist")
```

## License

MIT License
