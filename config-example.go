{
	"node": {
		"memoryRecordLimit": 40
	},
	"network": {
		"consensusAlgorithm": "raft",
		"role": 2,
		"maxPeers": 10,
		"discoveryURL": ["http://127.0.0.1"],
		"comments": {
			"consensusAlgorithm": "Currently only 'raft' is available (string)",
			"role": "Default role of node on network. '0':dynamic, '1':follow, '2':lead (integer)",
			"maxPeers": "Maximum number of concurrent peer connections (integer)",
			"discoverURL": "Array of URLs to use when discovering peers"
		}
	},
	"logging": {
		"enabled": true,
		"level": 4,
		"console": true,
		"file": false,
		"fileLocation": "/var/log/proliferate",
		"comments": {
			"enabled": "Enables/disables all logging. Available values (boolean)",
			"level": "Log level, higher numbers produce more logs. (integer) Values: '0': Fatal, '1': Error, '2': Warning, '3': Notice, '4': Verbose, '5': Noisy",
			"console": "Emit logs to console (boolean)'",
			"file": "Log output file location (string)"
		}
	},
	"couchDB": {
		"enabled": true,
		"host": "127.0.0.1",
		"port": "5984",
		"protocol": "http",
		"database": "proliferate",
		"comments": {
			"enabled": "Enables/disables CouchDB for block storage",
			"host": "CouchDB host domain or IP",
			"port": "CouchDB http port",
			"protocol": "HTTP protocol 'http' or 'https' (string)" 
		}
	}
}
