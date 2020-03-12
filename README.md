# proliferate
Proliferate Blockchain Framework

## CouchDB Requirements

### Index
```json
{
   "index": {
      "fields": [
         "serial"
      ]
   },
   "name": "serial-json-index",
   "type": "json"
}
```

### maxSerial Map
```javascript
function(doc){
  emit("serial", doc.serial);
}
```

### maxSerial Reduce
```javascript
function (key, values, rereduce) {
  return Math.max.apply({}, values);
}
```
