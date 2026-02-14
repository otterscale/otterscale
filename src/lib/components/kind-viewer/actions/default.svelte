<script lang="ts">
	import Ellipsis from '@lucide/svelte/icons/ellipsis';

	import Button from '$lib/components/ui/button/button.svelte';
</script>

<Button size="icon" variant="ghost" aria-label="Edit item" disabled>
	<Ellipsis size={16} aria-hidden="true" />
</Button>
