const valuesMapList: Record<string, { [key: string]: string }> = {
	minio: {
		'service.type': 'NodePort',
		'service.nodePorts.api': '30001',
		'service.nodePorts.console': '30002'
	},
	nginx: {
		'service.type': 'NodePort',
		'service.nodePorts.http': '31001',
		'service.nodePorts.https': '31002'
	},
	grafana: {
		'service.type': 'NodePort',
		'service.nodePorts.grafana': '32001'
	},
	'code-server-go': {
		'codeServer.password': 'password'
	},
	'code-server-python': {
		'codeServer.password': 'password'
	}
};

export { valuesMapList };
