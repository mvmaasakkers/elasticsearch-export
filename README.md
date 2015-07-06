# elasticsearch-export
Elasticsearch Import/Export tool

This tool can be used to (partially) migrate elasticsearch indexes from one cluster to another. I started it because we needed an easy way to do query-based migrations between data-centers.

Features:

- Command line tool
- Recreate mapping on target (if none already specified)
- Source data can be based on query

## Usage examples

````
// Copy index1 from machine 1 to machine 2
elasticexport -sh localhost -si indexa -dh otherhost -di indexa
````

## Command line options

Key|Example|Description
-|-|-
sh|localhost|Source host
sp|9200|Port of source ES host
si|index|Source index
dh|otherhost|Destination host
dp|9200|Port of destination ES host
di|index|Destination index
q|_all:*|Querystring query to select data
ba|50|Bulk amount to define amount of data indexed per request.

# Changelog

## 1.0.0

- Initial version with basic import/export functionality