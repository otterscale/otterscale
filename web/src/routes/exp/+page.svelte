<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';
	import { Badge } from '$lib/components/ui/badge';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Button } from '$lib/components/ui/button';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ScopeService, type Scope } from '$gen/api/scope/v1/scope_pb';
	import {
		Essential_Type,
		EssentialService,
		type Essential
	} from '$gen/api/essential/v1/essential_pb';
	import { ApplicationService, type Application } from '$gen/api/application/v1/application_pb';
	import { getContext, onMount } from 'svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { page } from '$app/state';
	import {
		FacilityService,
		type Facility,
		type Facility_Status,
		type Facility_Unit
	} from '$gen/api/facility/v1/facility_pb';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');

	const scopeClient = createClient(ScopeService, transport);
	// const facilityClient = createClient(FacilityService, transport);
	scopeClient.listScopes({}).then((r) => {
		console.log(r.scopes);
	});

	const scopeOptions = writable<SingleSelect.OptionType[]>([]);
	let isScopesLoading = $state(true);
	async function fetchScopeOptions() {
		try {
			const response = await scopeClient.listScopes({});
			scopeOptions.set(
				response.scopes.map(
					(scope) =>
						({ value: scope.uuid, label: scope.name, icon: 'ph:cube' }) as SingleSelect.OptionType
				)
			);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isScopesLoading = false;
		}
	}

	// const facilityOptions = writable<SingleSelect.OptionType[]>([]);
	// let isFacilitiesLoading = $state(true);
	// async function fetchFacilityOptions() {
	// 	try {
	// 		const response = await facilityClient.listFacilities({
	// 			scopeUuid: selectedScope
	// 		});
	// 		facilityOptions.set(
	// 			response.facilities.map(
	// 				(facility) =>
	// 					({
	// 						value: facility.name,
	// 						label: facility.name,
	// 						icon: 'ph:cube'
	// 					}) as SingleSelect.OptionType
	// 			)
	// 		);
	// 	} catch (error) {
	// 		console.error('Error fetching:', error);
	// 	} finally {
	// 		isFacilitiesLoading = false;
	// 	}
	// }

	let selectedScope = $state('');
	// let selectedFacility = $state('');

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchScopeOptions();
			// await fetchFacilityOptions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		isMounted = true;
	});
</script>

{#if isMounted}
	<Picker.Root align="right">
		<Picker.Wrapper class="*:h-8">
			<Picker.Label>Scope</Picker.Label>
			<SingleSelect.Root options={scopeOptions} bind:value={selectedScope}>
				<SingleSelect.Trigger />
				<SingleSelect.Content>
					<SingleSelect.Options>
						<SingleSelect.Input />
						<SingleSelect.List>
							<SingleSelect.Empty>No results found.</SingleSelect.Empty>
							<SingleSelect.Group>
								{#each $scopeOptions as option}
									<SingleSelect.Item {option}>
										<Icon
											icon={option.icon ? option.icon : 'ph:empty'}
											class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
										/>
										{option.label}
										<SingleSelect.Check {option} />
									</SingleSelect.Item>
								{/each}
							</SingleSelect.Group>
						</SingleSelect.List>
					</SingleSelect.Options>
				</SingleSelect.Content>
			</SingleSelect.Root>
		</Picker.Wrapper>
	</Picker.Root>
{/if}
