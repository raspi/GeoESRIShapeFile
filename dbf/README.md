# dBase database file (.dbf)
REQUIRED

# Fields
Please note that dBase has *zero clue*:
- about character encoding, so you must know in which format the entries are
- which fields are NULL, `-1`, `0`, '', ' ' or other *nullable* types
- integers: 
  - is it 32 or 64 bit
  - is it signed/unsigned 
  - is it little/big endian
- in which format possible date and time format is, it might be: 
  - `dd.mm.yyyy hh:mm:ss`
  - `dd.mm.yyyy`
  - `yyyy-mm-dd`
  - etc..

**So you must know what converter(s) to use for each field!**

See [defaultconverters.go](defaultconverters.go) for default converters.