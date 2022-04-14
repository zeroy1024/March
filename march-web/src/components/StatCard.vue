<template>
	<NCard
		style="height: 100%"
		content-style="display: flex;padding: 0px;flex-direction: column;justify-content: space-between;"
	>
		<div class="content">
			<div class="title">{{ title }}</div>
			<div class="value">{{ value }}</div>
		</div>
		<div :id="chartKey" class="chart"></div>
	</NCard>
</template>

<script setup lang="ts">
import { TinyArea } from '@antv/g2plot';
import { onMounted, PropType, watchEffect } from 'vue';
import utils from '../utils';

const props = defineProps({
	chartKey: {
		type: String,
		required: true,
	},
	title: {
		type: String,
		required: true,
	},
	value: {
		type: String,
		required: true,
	},
	chartData: {
		type: Array as PropType<number[]>,
		required: true,
	},
	color: {
		type: String,
		required: true,
	},
});

let chart: TinyArea;

onMounted(() => {
	chart = new TinyArea(props.chartKey, {
		height: 70,
		padding: [5, -5, -1, -5],
		data: props.chartData,
		// smooth: true,
		yAxis: {
			min: 0,
			max: 100,
		},
		line: {
			size: 1.5,
			color: utils.hexToRgba(props.color, 1),
		},
		areaStyle: {
			fill: utils.hexToRgba(props.color, 0.5),
		},
		animation: false,
	});
	chart.render();

	watchEffect(() => {
		chart.changeData(props.chartData);
	});
});
</script>

<style scoped>
.content {
	padding: 24px;
}

.content .title {
	color: rgb(118, 124, 130);
}

.content .value {
	font-size: 24px;
}
</style>
