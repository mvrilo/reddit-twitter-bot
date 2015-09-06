# reddit-twitter-bot

A reddit bot that looks periodically for new or hot posts in a subreddit and post them on twtiter.

# Usage

```
Usage of reddit-twitter-bot:
  -access_secret string
        Twitter access secret
  -access_token string
        Twitter access key
  -hot
        Lookup for hot posts instead of new
  -key string
        Twitter API consumer key
  -secret string
        Twitter API consumer secret
  -subreddit string
        Subreddit to watch
  -time int
        Time in seconds to fetch for posts (default 30)
```

Exmple:

```
reddit-twitter-bot -hot \
                   -subreddit=programming \
                   -time=30 \
                   -key=$KEY \
                   -secret=$SECRET \
                   -access_token=$ACCESS_TOKEN \
                   -access_secret==$ACCESS_SECRET
```

# License

MIT

# Author

Murilo Santana <<mvrilo@gmail.com>>
