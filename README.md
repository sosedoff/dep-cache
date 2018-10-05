# dep-cache

Simple utility to cache dependencies to S3.

## Config

```json
{
  "s3": {
    "key": "$CACHE_AWS_KEY",
    "secret": "$CACHE_AWS_SECRET",
    "region": "$CACHE_S3_REGION",
    "bucket": "$CACHE_S3_BUCKET"
  },
  "cache": [
    {
      "manifest": "Gemfile.lock",
      "path": ".bundle",
      "prefix": "bundler"
    },
    {
      "manifest": "package-lock.json",
      "path": "node_modules",
      "prefix": "npm"
    }
  ]
}
```

Values of $CACHE_AWS_KEY/$CACHE_AWS_SECRET will be replaced with environment variables.

## Usage

```
$ dep-cache download -c ./config.json
$ bundle install --path .bundle --jobs=4
$ npm install
$ dep-cache upload -c ./config.json
```
