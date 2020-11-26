# Version 1 of Dervaze JSON API


## `/v1/json/prefix/tr/?q=<word>`

Sends a list of Turkish words starting with `word` sorted by length

```
[ "word", "worda", "wordb", "wordabc"]
```

## `/v1/json/prefix/ot/?q=<word>`

Sends a list of Ottoman words starting with `word`

```
[ "word", "worda", "wordb", "wordabc"]
```

## `/v1/json/calc/abjad/?q=<word>`

Calculates abjad for the `word` given in unicode

```
{ "ottoman_unicode": "word",
"abjad": 1234 }
```

## `/v1/json/exact/tr/?q=<word>`

Returns records with Turkish Latin == `word`

```
[ { "ottoman_unicode": "abcd",
    "latin": "word",
    "abjad": 1234,
    "meanings": [
        {"source": "dervaze",
         "meaning": "meaning 1 of word"}]
    },
    { "ottoman_unicode": "abcddd",
    "latin": "word",
    "abjad": 1255,
    "meanings": [
        {"source": "kanar",
         "meaning": "meaning 1 of word"}]
    }]
    ```

## `/v1/json/exact/ot/?q=<word>`

Returns records with Ottoman == `word`

```
[ { "ottoman_unicode": "word",
    "latin": "abdc",
    "abjad": 1298,
    "meanings": [
        {"source": "dervaze",
         "meaning": "meaning 1 of word"}]
    },
    { "ottoman_unicode": "word",
    "latin": "anbdn",
    "abjad": 2191,
    "meanings": [
        {"source": "kanar",
         "meaning": "meaning 1 of word"}]
    }]
    ```

## `/v1/json/exact/abjad/?q=<number>`

Returns records with abjad == `number`

```
[ { "ottoman_unicode": "word",
    "latin": "abdc",
    "abjad": number,
    "meanings": [
        {"source": "dervaze",
         "meaning": "meaning 1 of word"}]
    },
    { "ottoman_unicode": "drwo",
    "latin": "anbdn",
    "abjad": number,
    "meanings": [
        {"source": "kanar",
         "meaning": "meaning 1 of word"}]
    }]
    ```


## `/v1/json/v2u/?q=<word>`

Converts `word` from visenc to unicode

```
{"ottoman_visenc": <word>,
"ottoman_unicode": <unicode>}
```

## `/v1/json/u2v/?q=<word>`

Converts `word` from unicode to visenc

```
{"ottoman_visenc": <visenc>,
"ottoman_unicode": <word>}
```

