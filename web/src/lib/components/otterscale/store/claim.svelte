<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog';

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title class="text-center text-xl">Claim</AlertDialog.Title>
		</AlertDialog.Header>
		<p class="hyphens-auto text-sm">
			The Helm charts recommended on this page are provided by their respective maintainers and are
			not developed, maintained, or directly affiliated with this project.
		</p>
		<p class="text-sm text-muted-foreground">
			For more information, please refer to the official sources.
		</p>
		<AlertDialog.Footer>
			<AlertDialog.Action class="text-xs" onclick={close}>Agree</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
