<script lang="ts">
	import { Multiple as Alert } from '$lib/components/custom/alert';

	const alerts: Alert.AlertType[] = [
		{
			variant: 'default',
			title: 'Welcome to Svelte Lab!',
			message: 'Explore Svelte features and components.',
			action: () => console.log('Welcome to Svelte Lab!')
		},
		{
			variant: 'information',
			title: 'Svelte Lab is Ready',
			message:
				'This is a development environment for Svelte. You can test components, styles, and functionality here before deploying to production.',
			action: () => console.log('Svelte Lab is Ready')
		},
		{
			variant: 'destructive',
			title:
				'Need Help? We Are Here To Assist You With All Your Technical Questions And Development Challenges In The World Of Svelte Programming',
			message:
				'Visit our documentation or join the community forum for detailed guides, troubleshooting tips, and to connect with other Svelte developers. Our support resources are available 24/7 to assist you with any questions or issues you may encounter.',
			action: () => console.log('Need Help?')
		}
	];
	let index: number = $state(0);
</script>

<Alert.Root {alerts} bind:index duration={2000}>
	<Alert.Icon />
	<Alert.Controller />
	<Alert.Title>{alerts[index].title}</Alert.Title>
	<Alert.Description>{alerts[index].message}</Alert.Description>
	<Alert.Action onclick={alerts[index].action}>Action</Alert.Action>
</Alert.Root>
