<script lang="ts" module>
	import * as Table from '$lib/components/custom/table';
</script>

<Table.Root>
	<Table.Caption>Caption</Table.Caption>
	<Table.Header>
		<Table.Row>
			<Table.Head>HEAD</Table.Head>
			<Table.Head>
				HEAD
				<Table.SubHead>SUB-HEAD</Table.SubHead>
			</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each Array(3) as _}
			<Table.Row>
				<Table.Cell>Cell</Table.Cell>
				<Table.Cell>
					Cell
					<Table.SubCell>SubCell</Table.SubCell>
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
	<Table.Footer>
		<Table.Row>
			<Table.Cell>footer</Table.Cell>
			<Table.Cell>
				footer
				<Table.SubCell>sub-footer</Table.SubCell>
			</Table.Cell>
		</Table.Row>
	</Table.Footer>
</Table.Root>
