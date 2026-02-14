import { Application_Type } from '$lib/api/application/v1/application_pb';

export const ApplicationTypeConfig = {
	[Application_Type.UNKNOWN]: { label: 'Unknown', icon: 'ph:question' },
	[Application_Type.DEPLOYMENT]: { label: 'Deployment', icon: 'ph:stack' },
	[Application_Type.STATEFUL_SET]: { label: 'StatefulSet', icon: 'ph:database' },
	[Application_Type.DAEMON_SET]: { label: 'DaemonSet', icon: 'ph:browsers' }
};
