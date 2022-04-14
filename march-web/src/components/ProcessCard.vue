<template>
	<NCard style="height: 380px" content-style="padding: 0;">
		<NDataTable
			size="small"
			:columns="colomns"
			:data="props.tableData"
			:style="{ height: `378px` }"
			flex-height
			:bordered="false"
			:bottom-bordered="false"
			:single-column="true"
			@update:sorter="handleSorterChange"
		/>
	</NCard>
</template>

<script setup lang="ts">
import { TableBaseColumn } from 'naive-ui/lib/data-table/src/interface';
import { PropType, ref } from 'vue';

const props = defineProps({
	tableData: {
		type: Array as PropType<Object[]>,
		required: true,
	},
});

const colomns = ref<TableBaseColumn[]>([
	{
		title: '进程ID',
		key: 'pid',
		align: 'center',
		sortOrder: false,
		sorter: 'default',
	},
	{
		title: '用户',
		key: 'runUser',
		align: 'center',
		sortOrder: false,
		sorter: 'default',
	},
	{
		title: 'CPU/%',
		key: 'cpuPercent',
		align: 'center',
		sortOrder: 'descend',
		sorter: 'default',
	},
	{
		title: '内存/%',
		key: 'memPercent',
		align: 'center',
		sortOrder: false,
		sorter: 'default',
	},
	{
		title: '命令',
		key: 'command',
		align: 'center',
		width: 300,
		render: (row) => {
			const command = row['command'] as string;
			if (command.length > 10) {
				return command.slice(0, 30) + '...';
			} else {
				return command.slice(0, 30);
			}
		},
	},
]);

const handleSorterChange = (sorter: Record<string, any>) => {
	colomns.value.map((column) => {
		if (column.sortOrder === undefined) return;
		if (!sorter) {
			column.sortOrder = false;
			return;
		}
		if (column.key == sorter.columnKey) {
			column.sortOrder = sorter.order;
		} else {
			column.sortOrder = false;
		}
	});
};
</script>

<style scoped></style>
