<script lang="ts">
	import Icon from '@iconify/svelte';
	import CreateCeph from '$lib/components/otterscale/management/general/create-ceph.svelte';
</script>

<div
	class="flex animate-pulse flex-col items-center justify-center p-24 hover:animate-none hover:cursor-pointer"
>
	<Icon icon="simple-icons:ceph" class="h-72 w-72" />
	<h1 class="text-lg">Ceph storage is not installed</h1>
	<p class="w-72 text-xs text-muted-foreground">
		Ceph provides reliable and scalable storage for your cluster. Install it to enable persistent
		storage features.
	</p>
	<CreateCeph />
</div>
