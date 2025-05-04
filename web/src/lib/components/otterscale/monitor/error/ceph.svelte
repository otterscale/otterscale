<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Card from '$lib/components/ui/card/index.js';
	import CreateCeph from '$lib/components/otterscale/management/general/create-ceph.svelte';
</script>

<!-- 
//TODO: add hover
Ceph provides reliable and scalable storage for your cluster. Install it to enable persistent storage features.
-->

<Card.Root class="h-min w-[350px]">
	<Card.Header class="flex flex-col items-center justify-center">
		<Icon icon="simple-icons:ceph" class="h-24 w-24 p-2" style="color: #f0424d" />
		<Card.Title>Ceph Storage</Card.Title>
		<Card.Description>Deploy Distributed Cluster</Card.Description>
	</Card.Header>
	<Card.Footer class="flex items-center justify-center">
		<CreateCeph />
	</Card.Footer>
</Card.Root>
