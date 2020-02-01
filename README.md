# datalock

Follow the rules on http://bit.ly/2UefZ1h

## Information for developers

* You can specify the format by adding `&_format=xml` or `&_format=json`; JSON is used as the default data format.
* You can shuffle response in methods that returns collections by adding `&_shuffle=0` or via Cookie `shuffle`; Shuffle can be based on any specified integer value or by current time if value set as `0`; Note that nested collections will be shuffled too.

### Season details

`GET /api/info_season?url=/serial-15825-Nelyudi-0-season.html`

* Type: `object`
* Shuffle support: `false`

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

* Type: `array`
* Shuffle support: `true`

Output format:

```json
[{
    "title": "...",
    "link": "..."
}]
```

### Available playlists

`GET /api/all_series?url=/serial-15825-Nelyudi-0-season.html`

> ***Note*** Method `/api/all_series` will be replaced by `/api/all_season_series` soon.

* Type: `array`
* Shuffle support: `true`

Output format:

```json
[{
    "name": "...",
    "playlist": [{
        "id": "...",
        "title": "...",
        "subtitle": "...",
        "file": "...",
        "galabel": "..."
    }]
}]
```

### Feeds

#### Updated series

`GET /api/updated_series`

* Type: `array`
* Shuffle support: `true`

Output format:

```json
[{
    "name": "...",
    "link": "...",
    "comment": "..."
}]
```

#### Popular series

`GET /api/popular_series`

* Type: `array`
* Shuffle support: `true`

Output format:

```json
[{
    "name": "...",
    "link": "...",
    "comment": "..."
}]
```


#### New series

`GET /api/new_series`

* Type: `array`
* Shuffle support: `true`

Output format:

```json
[{
    "name": "...",
    "link": "...",
    "comment": "..."
}]
```

### Health checks

* `GET /healthz` – Must returns `ok`;
* `GET /metrics` – Returns metrics for Prometheus.
