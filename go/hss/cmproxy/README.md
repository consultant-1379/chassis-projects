
The environment variable ADP_CM_URI is needed to be configured whit format:
http://<hostname>:<port>/cm/api/v1/configurations/<configuration name>

example:
ADP_CM_URI=http://172.17.0.7:5003/cm/api/v1/configurations/smfregistration

before that, the schema and configuration need to be created firstly,

example:

curl -D - -o - -sS -X POST -H "Content-Type: application/json" -d '{"name":"smfreg","title":"smfreg schema","jsonSchema":{
  "type": "object",
  "properties": {
    "smfreg": {
      "type": "object",
      "properties": {
        "service": {
          "type": "object",
		  "required":["port"],
		  "properties":{
			"port": {"type":"string"}},
		  "additionalProperties": true
		 },
		 "common": {
          "type": "object",
		  "required":["cmport","pmport"],
		  "properties":{
			"cmport": {"type":"string"},
			"pmport":{"type":"string"},
			"logport":{"type":"string"},
			"dbEndpoint":{"type":"string"}},
		  "additionalProperties": false
		 }},
      "additionalProperties": true
    }
  },
  "additionalProperties": false
}}' http://172.17.0.7:5003/cm/api/v1/schemas

curl -D - -o - -sS -X POST -H "Content-Type: application/json" -d '{"name":"smfreg","title":"smfreg config","data":{
  "smfreg": {
    "service":{
		"port":"9001"
		}, 
    "common":{
		"cmport":"9080",
		"pmport":"9100",
		"logport":"9090",
		"dbEndpoint":"udm-dbproxy:9001"
		}		
  }
}}' http://172.17.0.7:5003/cm/api/v1/configurations