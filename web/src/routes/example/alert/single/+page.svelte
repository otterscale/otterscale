<script lang="ts">
	import { Single as Alert } from '$lib/components/custom/alert';

	const alert: Alert.AlertType = {
		variant: 'destructive',
		title: 'Need Help?',
		message:
			'Visit our documentation or join the community forum for detailed guides, troubleshooting tips, and to connect with other Svelte developers. Our support resources are available 24/7 to assist you with any questions or issues you may encounter.',
		action: () => console.log('Need Help?')
	};
</script>

<Alert.Root {alert}>
	<Alert.Icon />
	<Alert.Title>{alert.title}</Alert.Title>
	<Alert.Description>{alert.message}</Alert.Description>
	<Alert.Action onclick={alert.action}>Action</Alert.Action>
</Alert.Root>
