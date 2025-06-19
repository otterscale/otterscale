<script lang="ts">
	import {
		Single as SingleAlert,
		Multiple as MultipleAlert
	} from '$lib/components/custom/alert/index';

	const alert: SingleAlert.AlertType = {
		level: 'ERROR',
		title: 'Need Help?',
		message:
			'Visit our documentation or join the community forum for detailed guides, troubleshooting tips, and to connect with other Svelte developers. Our support resources are available 24/7 to assist you with any questions or issues you may encounter.',
		action: () => console.log('Need Help?')
	};

	const alerts: SingleAlert.AlertType[] = [
		{
			level: 'DEBUG',
			title: 'Welcome to Svelte Lab!',
			message: 'Explore Svelte features and components.',
			action: () => console.log('Welcome to Svelte Lab!')
		},
		{
			level: 'INFORMATION',
			title: 'Svelte Lab is Ready',
			message:
				'This is a development environment for Svelte. You can test components, styles, and functionality here before deploying to production.',
			action: () => console.log('Svelte Lab is Ready')
		},
		{
			level: 'ERROR',
			title: 'Need Help?erfgjwnsdlcfvher;ifjsn/dlvijhsn.dlk/fjnws.delcfjwnsed.fkwnef.',
			message:
				'Visit our documentation or join the community forum for detailed guides, troubleshooting tips, and to connect with other Svelte developers. Our support resources are available 24/7 to assist you with any questions or issues you may encounter.',
			action: () => console.log('Need Help?')
		}
	];
	function variantGetter(level: string) {
		switch (level) {
			case 'INFORMATION':
				return 'information';
			case 'ERROR':
				return 'destructive';
			default:
				return 'default';
		}
	}
	let index: number = $state(0);
</script>

<SingleAlert.Root>
	<SingleAlert.Icon />
	<SingleAlert.Title>{alert.title}</SingleAlert.Title>
	<SingleAlert.Description>{alert.message}</SingleAlert.Description>
	<SingleAlert.Action onclick={alert.action}>Action</SingleAlert.Action>
</SingleAlert.Root>

<MultipleAlert.Root {alerts} bind:index duration={2000} {variantGetter}>
	<MultipleAlert.Icon />
	<MultipleAlert.Controller />
	<MultipleAlert.Title>{alerts[index].title}</MultipleAlert.Title>
	<MultipleAlert.Description>{alerts[index].message}</MultipleAlert.Description>
	<MultipleAlert.Action onclick={alerts[index].action}>Action</MultipleAlert.Action>
</MultipleAlert.Root>
