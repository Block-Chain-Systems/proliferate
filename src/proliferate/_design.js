{
  "_id": "_design/ab9ea6009283e8a97e1a3f1e626b808b6fa61421",
  "_rev": "18-c346c94963de275797855f7b78451d9b",
  "language": "query",
  "views": {
    "serial-json-index": {
      "map": {
        "fields": {
          "serial": "desc"
        },
        "partial_filter_selector": {}
      },
      "reduce": "_count",
      "options": {
        "def": {
          "fields": [
            "serial"
          ]
        }
      }
    }
  }
}