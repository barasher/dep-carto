# Dependencies cartographier

![Build Status](https://github.com/barasher/picdexer/workflows/Dep-Carto/badge.svg)

## What is `dep-carto`

`dep-carto` describe dependencies between components.

It is composed of :
- a server that stores and exposes registered dependencies (push mode)
- a dependency crawler that parses configuration files and registers dependencies to the server

## `dep-carto` use-cases

Let's consider this architecture (acme company).

![architecture](doc/architecture.jpg)

Which components depends on Elasticsearch ? Which components will be impacted if Elasticsearch goes down ?

Once dependencies crawled, `dep-carto` can build this kind representation :

![representation](doc/representation.jpg)

## Getting `dep-carto`

### Binary release

### Docker

## Configuration
