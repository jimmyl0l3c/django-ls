# 0.1.0

- [x] Detecting django's project root and path to settings
- [x] Automatic introspection when opening project
- [x] Buffer sync (document open, close, change)
- [x] Detect when completions are for `values` or `filter` method
- [x] Detect a model of the queryset
- [x] Return completions for the models field lookups in `values`/`filter` methods
- [x] Only show completions for specific methods (like `values`/`filter`)
- [ ] Polish the logic a bit (thread synchronization)

# 0.2.0

- [ ] Change the completion value based on method (e.g. `values` expects strings, `filter` expects kwargs)
- [ ] Write tests for parser and safemap
- [ ] LS should have debug argument
- [ ] Show additional field lookups (e.g. `in`, `isnull`, ...)
- [ ] Improve completion parser (queryset/model detection)

# 0.3.0

- [ ] Show lookups of related fields up to a specified depth
- [ ] Add support for `select_related`, `prefetch_related`
- [ ] Add more py venv discovery strategies

# 0.4.0

- [ ] Support `F` and `Q` objects
- [ ] Re-run introspection when models / installed apps change
- [ ] Add caching

# 0.5.0

- [ ] Query diagnostics
