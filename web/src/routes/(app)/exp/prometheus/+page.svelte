<script lang="ts">
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { group } from 'd3-array';
	import { Arc, AreaChart, Group, LinearGradient, PieChart, Text } from 'layerchart';
	import { curveCatmullRom } from 'd3-shape';
	import { longData } from '../chart/gen';

	const dataByYear = group(longData, (d) => d.year);
	const data = dataByYear.get(2019) ?? [];

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	const prom = new PrometheusDriver({
		endpoint: 'http://10.102.197.157:32002',
		baseURL: '/api/v1' // default value
	});

	const start = new Date().getTime() - 1 * 60 * 60 * 1000; // past 1h
	const end = new Date();
	const step = 14; // 6 * 60 * 60: 1 point every 6 hours ??

	let iops = [] as SampleValue[];
	prom
		.rangeQuery('sum(irate(ceph_osd_op_w{}[5m]))', start, end, step)
		.then((res) => {
			if (res.result) {
				iops = res.result[0].values;
			}
		})
		.catch(console.error);

	let clusterCapacity = 0.0;
	prom
		.instantQuery(
			'(ceph_cluster_total_bytes{}-ceph_cluster_total_used_bytes{})/ceph_cluster_total_bytes{}'
		)
		.then((res) => {
			if (res.result) {
				clusterCapacity = res.result[0].value.value;
			}
		})
		.catch(console.error);
</script>

<div class="h-[300px] rounded border p-4">
	<AreaChart
		data={iops}
		x="time"
		y="value"
		props={{ area: { curve: curveCatmullRom } }}
		{renderContext}
		{debug}
	/>
</div>
<div class="h-[300px] resize overflow-auto rounded border p-4">
	<PieChart {renderContext} {debug}>
		<svelte:fragment slot="marks">
			<LinearGradient class="from-secondary to-primary" let:gradient>
				<Group y={20}>
					<Arc
						value={clusterCapacity * 100}
						domain={[0, 100]}
						outerRadius={80}
						innerRadius={-15}
						cornerRadius={10}
						padAngle={0.02}
						range={[-120, 120]}
						fill={gradient}
						track={{ class: 'fill-none stroke-surface-content/10' }}
						let:value
					>
						<Text
							value={value.toFixed(1) + '%'}
							textAnchor="middle"
							verticalAnchor="middle"
							class="text-4xl tabular-nums"
						/>
					</Arc>
				</Group>
			</LinearGradient>
		</svelte:fragment>
	</PieChart>
</div>
