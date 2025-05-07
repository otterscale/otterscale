<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Card from '$lib/components/ui/card';
	import CreateCeph from '$lib/components/otterscale/management/general/create-ceph.svelte';
</script>

<Card.Root>
	<Card.Content class="flex aspect-square items-center justify-center">
		<div class="flex-col items-center justify-center space-y-2">
			<div class="flex items-center justify-center space-x-2">
				<Icon icon="simple-icons:ceph" class="size-24 p-2" style="color: #f0424d" />
				<div class="flex-col space-y-1">
					<p class="text-xl font-semibold leading-none tracking-tight">Ceph</p>
					<p class="text-xl font-semibold leading-none tracking-tight">Storage</p>
				</div>
			</div>
			<div class="flex items-center justify-center">
				<p class="text-sm text-muted-foreground">Deploy Your Distributed Storage Cluster</p>
			</div>
			<div class="flex items-center justify-center pt-6">
				<CreateCeph />
			</div>
		</div>
	</Card.Content>
</Card.Root>
