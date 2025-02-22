<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
</script>

<Sheet.Root>
	<Sheet.Trigger>
		<Button variant="outline" size="icon" class="bg-header">
			<Icon icon="ph:notification" class="h-5 w-5" />
			<!-- <Icon icon="ph:notification-fill" class="h-5 w-5" /> -->
		</Button>
	</Sheet.Trigger>
	<Sheet.Content>
		<Sheet.Header>
			<Sheet.Title>Are you sure absolutely sure?</Sheet.Title>
			<Sheet.Description>
				This action cannot be undone. This will permanently delete your account and remove your data
				from our servers.
			</Sheet.Description>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>
