# A Stupid Bot to scrape Github API for new good first issue and help wanted

* It uses git submodule as a source of truth 
* To make sure you don't spam others, turn your source of truth repo into **private**

## TODO:
- [ ] Figure a way to inject credential to `git pull`
- [ ] `Tracker` interfaces need to have `init`
- [ ] Propagate Tracker into `issues` package
- [ ] Deploy to local Kubernetes cluster
- [ ] Allow Update for source of truth from local in addition to remote