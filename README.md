# datalock

Follow the rules on http://bit.ly/2rdoNTn

## Information for developers

### Season details

`GET /api/info_season?url=/serial-15825-Nelyudi-0-season.html`

Output format:

```json
{
    "title": "...",
    "id": 0,
    "serial": 0,
    "keywords": "...",
    "description": "..."
}
```

### Available seasons

`GET /api/all_seasons?url=/serial-15825-Nelyudi-0-season.html`

Output format:

```json
{
    "seasons": [{
        "title": "...",
        "link": "..."
    }]
}
```

### Available playlists

`GET /api/all_series?url=/serial-15825-Nelyudi-0-season.html`

Output format:

```json
[{
    "name": "...",
    "playlist": [{
        "comment": "...",
        "file": "...",
        "streamsend": "sec",
        "galabel": "0_0"
    }, {
        "comment": "...",
        "file": "...",
        "streamsend": "sec",
        "galabel": "0_0"
    }]
},
...
]
```
