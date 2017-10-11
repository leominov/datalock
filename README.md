# datalock

Follow the rules on http://bit.ly/2rdoNTn

## Information for developers

### Season details

`GET /api/info_season?url=/serial-15825-Nelyudi-0-season.html`

Output format:

```json
{
    "title": "...",
    "id": 15825,
    "serial": 9756,
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
        "link": "/serial-15825-Nelyudi-0-season.html"
    }]
}
```

### Available playlists

`GET /api/all_series?url=/serial-15825-Nelyudi-0-season.html`

Output format:

```json
[{
    "name": "Hamster",
    "playlist": [{
        "comment": "...",
        "file": "...",
        "streamsend": "sec",
        "galabel": "15825_482304"
    }, {
        "comment": "...",
        "file": "...",
        "streamsend": "sec",
        "galabel": "15825_483231"
    }]
},
...
]
```
