<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.time_till_full();
	const CHART_DESCRIPTION = m.six_hour_average();

	const query = $derived(
		`
	(
		ceph_pool_max_avail{
			job=~".+",
			juju_application=~".*",
			juju_model=~".*",
			juju_model_uuid=~"${scope.uuid}",
			juju_unit=~".*"
		} 
		/ 
		deriv(ceph_pool_stored{
			job=~".+",
			juju_application=~".*",
			juju_model=~".*",
			juju_model_uuid=~"${scope.uuid}",
			juju_unit=~".*"
		}[6h])
	) 
	* 
	on(pool_id) group_left(instance, name) ceph_pool_metadata{
		job=~".+",
		juju_application=~".*",
		juju_model=~".*",
		juju_model_uuid=~"${scope.uuid}",
		juju_unit=~".*",
		name=~".mgr"
	} > 0
	`
	);

	// Format time duration
	function formatTimeTillFull(days: number): string {
		if (!isFinite(days) || days <= 0) {
			return '∞ years';
		}

		if (days < 1) {
			const hours = Math.round(days * 24);
			return `${hours} hour${hours !== 1 ? 's' : ''}`;
		} else if (days < 30) {
			const roundedDays = Math.round(days);
			return `${roundedDays} day${roundedDays !== 1 ? 's' : ''}`;
		} else if (days < 365) {
			const months = Math.round(days / 30);
			return `${months} month${months !== 1 ? 's' : ''}`;
		} else {
			const years = Math.round(days / 365);
			return `${years} year${years !== 1 ? 's' : ''}`;
		}
	}
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	<Card.Root class="gap-2">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">
			{#if response.result.length > 0 && response.result[0].value}
				{@const days = parseFloat(response.result[0].value.value)}
				{formatTimeTillFull(days)}
			{:else}
				∞ years
			{/if}
		</Card.Content>
	</Card.Root>
{:catch error}
	<Card.Root class="gap-2">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">LOADING ERROR</Card.Content>
	</Card.Root>
{/await}
