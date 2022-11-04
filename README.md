# go-config
Used to import general use libraries and corresponding env vars/flag setup

## Currently Supported
- Cache
  - Go Cache
  - Mem Cache
  - Cache Retry
    - will all caches to be retried if they fail
  - Tiered Cache
    - Used to try multiple caches before falling through to original caching function
- db
  - my sql/maria db
- logging
  - zap logger