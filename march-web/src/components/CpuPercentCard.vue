<template>
	<NCard title="CPU使用率">
		<template #header-extra>
			<div class="toolbar">
				<NRadioGroup
					name="date-range"
					v-model:value="state.dateRangeChecked"
					@update:value="dateRangeChange"
				>
					<NRadioButton key="realTime" label="实时" value="realTime" />
					<NRadioButton key="24h" label="近24小时" value="24h" />
					<NRadioButton key="7d" label="近7天" value="7d" />
					<NRadioButton key="30d" label="近30天" value="30d" />
				</NRadioGroup>
				<n-select
					v-model:value="state.timeSize"
					:options="state.timeSizeOption"
					@update:value="() => loadData()"
					style="margin-left: 10px; width: 100px"
				/>
			</div>
		</template>
		<div id="chart-cpu-percent-big" />
	</NCard>
</template>

<script setup lang="ts">
import { Datum, Line } from '@antv/g2plot';
import { onMounted, reactive, ref, watchEffect } from 'vue';
import http from '../http';

let chart: Line;
const state = reactive({
	data: [],
	dateStart: '',
	dateEnd: '',
	realTime: true,
	dateRangeChecked: 'realTime',

	timeSize: 5,
	timeSizeOption: [
		{
			label: '5秒',
			value: 5,
		},
		{
			label: '10秒',
			value: 10,
		},
		{
			label: '30秒',
			value: 30,
		},
		{
			label: '1分钟',
			value: 60,
		},
		{
			label: '5分钟',
			value: 300,
		},
		{
			label: '10分钟',
			value: 600,
		},
	],
});
const timer = ref();

const loadData = async () => {
	const result = await http.getCpuPercent(
		state.dateStart,
		state.dateEnd,
		state.timeSize.toString()
	);
	if (result.code === 200) {
		state.data = result.data;
	}
};

const dateRangeChange = (val: string) => {
	state.realTime = val === 'realTime';
	const nowTimestamp = Math.floor(new Date().getTime() / 1000);
	switch (val) {
		case 'realTime':
			state.dateStart = state.dateEnd = '';
			state.timeSize = 5;
			state.timeSizeOption = [
				{
					label: '5秒',
					value: 5,
				},
				{
					label: '10秒',
					value: 10,
				},
				{
					label: '30秒',
					value: 30,
				},
				{
					label: '1分钟',
					value: 60,
				},
				{
					label: '5分钟',
					value: 300,
				},
				{
					label: '10分钟',
					value: 600,
				},
			];
			break;
		case '24h':
			state.dateStart = (nowTimestamp - 60 * 60 * 24).toString();
			state.dateEnd = nowTimestamp.toString();

			state.timeSize = 60;
			state.timeSizeOption = [
				{
					label: '1分钟',
					value: 60,
				},
				{
					label: '5分钟',
					value: 300,
				},
				{
					label: '10分钟',
					value: 600,
				},
				{
					label: '30分钟',
					value: 1800,
				},
				{
					label: '1小时',
					value: 3600,
				},
			];
			break;
		case '7d':
			state.dateStart = (nowTimestamp - 60 * 60 * 24 * 7).toString();
			state.dateEnd = nowTimestamp.toString();

			state.timeSize = 600;
			state.timeSizeOption = [
				{
					label: '10分钟',
					value: 600,
				},
				{
					label: '30分钟',
					value: 1800,
				},
				{
					label: '1小时',
					value: 3600,
				},
				{
					label: '12小时',
					value: 43200,
				},
			];
			break;
		case '30d':
			state.dateStart = (nowTimestamp - 60 * 60 * 24 * 30).toString();
			state.dateEnd = nowTimestamp.toString();

			state.timeSize = 3600;
			state.timeSizeOption = [
				{
					label: '1小时',
					value: 3600,
				},
				{
					label: '12小时',
					value: 43200,
				},
			];
			break;
	}
	loadData();
};

watchEffect(() => {
	if (state.realTime) {
		timer.value = setInterval(() => {
			loadData();
		}, state.timeSize * 1000);
	} else {
		clearInterval(timer.value);
	}
});

onMounted(() => {
	loadData(); // 首次加载数据

	chart = new Line('chart-cpu-percent-big', {
		data: [],
		xField: 'datetime',
		yField: 'value',
		seriesField: 'name',
		xAxis: {
			tickCount: 6,
		},
		yAxis: {
			min: 0,
			max: 100,
			label: {
				formatter: (text: string) => {
					return text + '%';
				},
			},
		},
		tooltip: {
			formatter: (datum: Datum) => {
				return { name: datum.name, value: datum.value.toFixed(2) + '%' };
			},
		},
		animation: false,
	});
	chart.render();

	// 监听数据变化时更新表格数据
	watchEffect(() => {
		chart.changeData(state.data);
	});
});
</script>

<style scoped>
.toolbar {
	display: flex;
}
</style>
