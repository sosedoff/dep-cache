# dep-cache

Simple utility to cache dependencies to S3.

## Config

```json
{
  "s3": {
    "key": "... AWS access key ...",
    "secret": "... AWS secret key ...",
    "region": "us-east-1",
    "bucket": "AWS bucket name"
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

## Usage

```
$ dep-cache download -c ./config.json
$ bundle install --path .bundle --jobs=4
$ npm install
$ dep-cache upload -c ./config.json
```
