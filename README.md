# dep-cache

Parallel dependency caching utility for use in CI environments.

## Overview

If your project uses multiple dependency management systems like Rubygems and
NPM, most likely your CI worker spends too much time installing all dependencies
on every run. Dep cache tool helps you to speed up install time by storing the
dependencies in Amazon S3 and reusing them as necessary. In parallel.

## Config

The example configuration file:

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

Instead of putting your real credentials into `key` and `secret` fields you can 
set them with `$CACHE_AWS_KEY` and `$CACHE_AWS_SECRET` respectively. These values
will be replaced automatically if your system has those set as environment variables.

For Amazon EC2 users with IAM roles configured: you dont have to set `key`, `secret` 
and `region`. Dep cache will automatically pick up configuration from the EC2 environment.
The only required option is `bucket`.

## Usage

Take a look at the example project structure (real files omitted):

```
.
..
bundle
Gemfile
Gemfile.lock
node_modules
package-lock.json
package.json
```

Use the example configuration file and save it as `.dep-cache.json`.

Next, perform download, install and upload of dependencies:


```bash
# This will download the archives if they exist in S3
$ dep-cache download 

# Install dependencies. These would run super fast if the cache exists
$ bundle install --path .bundle --jobs=4
$ npm install

# Upload dependencies bundles to S3. Upload is skipped if there are no changes.
$ dep-cache upload
```

Alternatively, you can set config path with:

```bash
$ dep-cache -c ./configs/dep-cache.json
```

In case if you need to invalidate the existing caches you can use the following command:

```bash
$ dep-cache reset
```

## License

MIT License
