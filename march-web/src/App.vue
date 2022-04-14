<template>
	<NConfigProvider :theme-overrides="{ common: { fontWeightStrong: '600' } }">
		<NMessageProvider>
			<NLayout>
				<NLayoutHeader class="header">
					<div class="logo">March Dashboard</div>
				</NLayoutHeader>
				<NLayoutContent class="content" content-style="padding: 24px 48px;">
					<Dashboard />
				</NLayoutContent>
				<NLayoutFooter class="footer">
					&copy; 2022 <a href="#">ZeroY</a>
				</NLayoutFooter>
			</NLayout>
		</NMessageProvider>
	</NConfigProvider>
</template>

<script setup lang="ts">
import Dashboard from './pages/Dashboard.vue';
</script>

<style scoped>
* {
	color: #333333;
	text-decoration: none;
}

.header {
	height: 64px;
	line-height: 64px;
	background: #001529;
	color: #ffffff;
	padding: 0 60px;
	position: fixed;
	top: 0;
	left: 0;
	z-index: 1001;
}

.header .logo {
	color: #ffffff;
	font-weight: 600;
	font-size: 18px;
}

.content {
	margin-top: 64px;
	background: #f0f2f5;
}

.footer {
	background: #f0f2f5;
	text-align: center;
}

.footer a {
	color: rgba(51, 51, 51, 0.7);
}
</style>
