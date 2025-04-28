<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { goto } from '$app/navigation';
</script>

<Card.Root>
	<Card.Content>
		<h3 class="text-2xl tracking-tight">This page is currently under construction.</h3>
		<p class="text-sm text-muted-foreground">Please check back later.</p>
		<Button class="mt-4" href="/">Home</Button>
	</Card.Content>
</Card.Root>
