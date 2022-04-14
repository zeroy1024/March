<template>
	<NCard title="磁盘使用率" style="height: 380px">
		<div id="chart-disk-usage"></div>
	</NCard>
</template>

<script setup lang="ts">
import { ListItem } from '@antv/component';
import { Bar } from '@antv/g2plot';
import { onMounted, PropType, ref, watchEffect } from 'vue';
import { IDiskUsage } from '../interfaces';
import utils from '../utils';

const props = defineProps({
	chartData: {
		type: Array as PropType<IDiskUsage[]>,
		required: true,
	},
});

const deviceInfo = ref<Record<string, any>>({});

let chart: Bar;

const data = [
	{
		country: 'Asia',
		year: '1750',
		value: 502,
	},
	{
		country: 'Asia',
		year: '1800',
		value: 635,
	},
	{
		country: 'Asia',
		year: '1850',
		value: 809,
	},
	{
		country: 'Asia',
		year: '1900',
		value: 947,
	},
	{
		country: 'Asia',
		year: '1950',
		value: 1402,
	},
	{
		country: 'Asia',
		year: '1999',
		value: 3634,
	},
	{
		country: 'Asia',
		year: '2050',
		value: 5268,
	},
	{
		country: 'Africa',
		year: '1750',
		value: 106,
	},
	{
		country: 'Africa',
		year: '1800',
		value: 107,
	},
	{
		country: 'Africa',
		year: '1850',
		value: 111,
	},
	{
		country: 'Africa',
		year: '1900',
		value: 133,
	},
	{
		country: 'Africa',
		year: '1950',
		value: 221,
	},
	{
		country: 'Africa',
		year: '1999',
		value: 767,
	},
	{
		country: 'Africa',
		year: '2050',
		value: 1766,
	},
	{
		country: 'Europe',
		year: '1750',
		value: 163,
	},
	{
		country: 'Europe',
		year: '1800',
		value: 203,
	},
	{
		country: 'Europe',
		year: '1850',
		value: 276,
	},
	{
		country: 'Europe',
		year: '1900',
		value: 408,
	},
	{
		country: 'Europe',
		year: '1950',
		value: 547,
	},
	{
		country: 'Europe',
		year: '1999',
		value: 729,
	},
	{
		country: 'Europe',
		year: '2050',
		value: 628,
	},
];

onMounted(() => {
	chart = new Bar('chart-disk-usage', {
		height: 280,
		data: props.chartData,
		xField: 'value',
		yField: 'device',
		seriesField: 'type',
		isPercent: true,
		isStack: true,
		animation: false,
		minBarWidth: 20,
		maxBarWidth: 40,
		xAxis: {
			tickCount: 5,
			label: {
				formatter: (text: string, item: ListItem) => {
					return item.value * 100 + '%';
				},
			},
		},
		tooltip: {
			showTitle: false,
			customContent: (title) => {
				if (title) {
					const device = title.split('\n')[0];
					const di = deviceInfo.value[device];
					return `
                    <div style="padding:12px;">
                        <div style="padding: 5px 0">
                            挂载点: ${di['mountPoint']}
                        </div>
                        <div style="padding: 5px 0">
                            文件系统: ${di['fsType']}
                        </div>
                        <div style="padding: 5px 0">
                            已用/可用: 
                            ${utils.formatBytes(di['usedBytes'])}
                            /${utils.formatBytes(di['freeBytes'])}
                        </div>
                    </div>
                    `;
				}
				return '';
			},
		},
		label: {
			position: 'middle',
			content: (item) => {
				return (item.value * 100).toFixed(2) + '%';
			},
			style: {
				fill: '#fff',
			},
		},
	});
	chart.render();
    
	watchEffect(() => {
		const chartData: Record<string, any>[] = [];
		props.chartData.map((diskUsage) => {
			deviceInfo.value[diskUsage.device] = {
				freeBytes: diskUsage.freeBytes,
				fsType: diskUsage.fsType,
				mountPoint: diskUsage.mountPoint,
				totalBytes: diskUsage.totalBytes,
				usedBytes: diskUsage.usedBytes,
				usedPercent: diskUsage.usedPercent,
			};

			const device =
				diskUsage.device + '\n' + utils.formatBytes(diskUsage.totalBytes);
			chartData.push({
				device: device,
				type: 'used',
				value: diskUsage.usedBytes,
			});
			chartData.push({
				device: device,
				type: 'free',
				value: diskUsage.freeBytes,
			});
		});
		chart.changeData(chartData);
	});
});
</script>

<style scoped></style>
