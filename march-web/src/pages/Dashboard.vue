<template>
	<NGrid :x-gap="12" :y-gap="12" :cols="24">
		<NGi :span="12">
			<NCard style="min-height: 80px">
				<NDescriptions :column="4" title="主机信息">
					<NDescriptionsItem :span="2">
						<template #label>
							<span class="nd-title">
								<NIcon :component="Server" :size="20" />
								<span class="text">主机名:</span>
							</span>
						</template>
						{{ state.hostInfo.hostname }}
					</NDescriptionsItem>
					<NDescriptionsItem :span="1">
						<template #label>
							<span class="nd-title">
								<NIcon :component="InfoCircle" :size="20" />
								<span class="text">系统:</span>
							</span>
						</template>
						{{ state.hostInfo.platform }} {{ state.hostInfo.os }}
					</NDescriptionsItem>
					<NDescriptionsItem :span="1">
						<template #label>
							<span class="nd-title">
								<NIcon :component="Linux" :size="20" />
								<span class="text">内核:</span>
							</span>
						</template>
						{{ state.hostInfo.kernelVersion }}
						{{ state.hostInfo.kernelArch }}
					</NDescriptionsItem>
					<NDescriptionsItem :span="2">
						<template #label>
							<span class="nd-title">
								<NIcon :component="Cpu" :size="20" />
								<span class="text">处理器:</span>
							</span>
						</template>
						{{ state.cpuInfo.modelName }}
					</NDescriptionsItem>
					<NDescriptionsItem :span="1">
						<template #label>
							<span class="nd-title">
								<NIcon :component="Memory" :size="20" />
								<span class="text">内存:</span>
							</span>
						</template>
						{{ state.hostInfo.maxMem }}
					</NDescriptionsItem>
					<NDescriptionsItem :span="1">
						<template #label>
							<span class="nd-title">
								<NIcon :component="Clock" :size="20" />
								<span class="text">启动时间:</span>
							</span>
						</template>
						<NTime :time="state.hostInfo.bootTime" unix />
					</NDescriptionsItem>
				</NDescriptions>
			</NCard>
		</NGi>
		<NGi :span="6">
			<StatCard
				chart-key="chart-cpu-percent"
				title="CPU"
				:value="current.cpu + '%'"
				:chart-data="current.cpuList"
				color="#2ed573"
			/>
		</NGi>
		<NGi :span="6">
			<StatCard
				chart-key="chart-mem-percent"
				title="Memory"
				:value="current.mem + '%'"
				:chart-data="current.memList"
				color="#ffa502"
			/>
		</NGi>
		<NGi :span="12">
			<DiskUsageCard :chart-data="state.diskUsage" />
		</NGi>
		<NGi :span="12">
			<ProcessCard :table-data="state.processInfo" />
		</NGi>
		<NGi :span="24">
			<CpuPercentCard />
		</NGi>
		<NGi :span="24">
			<MemPercentCard />
		</NGi>
		<NGi :span="24">
			<DiskIOCard />
		</NGi>
		<NGi :span="24">
			<NetworkLoadCard />
		</NGi>
	</NGrid>
</template>

<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { Cpu, Server, InfoCircle, Clock, Alarm } from '@vicons/tabler';
import { Memory, Linux } from '@vicons/fa';
import { useMessage } from 'naive-ui';

import http from '../http';
import { ICpuInfo, IHostInfo, IDiskUsage, IProcessInfo } from '../interfaces';

import StatCard from '../components/StatCard.vue';
import DiskUsageCard from '../components/DiskUsageCard.vue';
import ProcessCard from '../components/ProcessCard.vue';
import CpuPercentCard from '../components/CpuPercentCard.vue';
import MemPercentCard from '../components/MemPercentCard.vue';
import DiskIOCard from '../components/DiskIOCard.vue';
import NetworkLoadCard from '../components/NetworkLoadCard.vue';

const message = useMessage();

const current = reactive({
	cpu: 0.0,
	mem: 0.0,
	cpuList: [0] as number[],
	memList: [0] as number[],
});

const state = reactive({
	hostInfo: {
		bootTime: 0,
		cpuCore: 0,
		cpuThread: 0,
		hostId: '',
		hostname: '',
		kernelArch: '',
		kernelVersion: '',
		maxMem: '',
		os: '',
		platform: '',
		platformFamily: '',
		platformVersion: '',
		uptime: 0,
	} as IHostInfo,
	cpuInfo: {
		modelName: '',
		arch: '',
		frequency: 0,
		cacheSize: 0,
		flags: [],
	} as ICpuInfo,
	diskUsage: [] as IDiskUsage[],
	processInfo: [] as IProcessInfo[],
});

const loadCurrent = async () => {
	const result = await http.getCurrent();
	if (result.code === 200) {
		// cpu
		current.cpu = parseFloat(parseFloat(result.data['cpuPercent']).toFixed(2));
		if (current.cpuList.length >= 60) {
			current.cpuList.splice(0, 1);
		}
		current.cpuList.push(current.cpu);

		//mem
		current.mem = parseFloat(parseFloat(result.data['memPercent']).toFixed(2));
		if (current.memList.length >= 60) {
			current.memList.splice(0, 1);
		}
		current.memList.push(current.mem);
	}
};

const loadHostInfo = async () => {
	const result = await http.getHostInfo();
	if (result.code === 200) {
		state.hostInfo = result.data;
	} else {
		message.error(result.message);
	}
};

const loadCpuInfo = async () => {
	const result = await http.getCpuInfo();
	if (result.code === 200) {
		state.cpuInfo.modelName = result.data[0]['modelName'];
		state.cpuInfo.arch = result.data[0]['model'];
		state.cpuInfo.frequency = result.data[0]['mhz'];
		state.cpuInfo.cacheSize = result.data[0]['cacheSize'];
		state.cpuInfo.flags = result.data[0]['flage'];
	}
};

const loadDiskUsage = async () => {
	const result = await http.getDiskUsage();
	if (result.code === 200) {
		state.diskUsage = result.data;
	}
};

const loadProcessInfo = async () => {
	const result = await http.getProcessInfo();
	if (result.code === 200) {
		state.processInfo = result.data;
	}
};

onMounted(() => {
	loadHostInfo();
	loadCpuInfo();
	loadDiskUsage();

	setInterval(() => {
		loadCurrent();
		loadProcessInfo();
	}, 1000);
});
</script>

<style scoped>
.nd-title {
	display: flex;
	align-items: center;
}

.nd-title .text {
	padding-left: 5px;
}
</style>
