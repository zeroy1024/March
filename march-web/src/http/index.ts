import { IResponse } from '../interfaces';

// const baseURL = 'http://127.0.0.1:52080/api/v1';
const baseURL = '/api/v1';

const request = (url: string, options?: RequestInit) => {
	return new Promise<IResponse>((resolve, reject) => {
		fetch(baseURL + url, {
			headers: {},
			...options,
		})
			.then((res) => res.json())
			.then((res) => resolve(res))
			.catch((res) => reject(res));
	});
};

export default {
	getCurrent: () => {
		return request('/current/', {
			method: 'GET',
		});
	},

	getHostInfo: () => {
		return request('/host/', {
			method: 'GET',
		});
	},

	getCpuInfo: () => {
		return request('/cpu/', {
			method: 'GET',
		});
	},

	getCpuPercent: (dateStart: string, dateEnd: string, timeSize: string) => {
		return request(
			'/cpu/percent?dateStart=' +
				dateStart +
				'&dateEnd=' +
				dateEnd +
				'&timeSize=' +
				timeSize,
			{
				method: 'GET',
			}
		);
	},

	getMemUsage: (dateStart: string, dateEnd: string, timeSize: string) => {
		return request(
			'/mem/percent?dateStart=' +
				dateStart +
				'&dateEnd=' +
				dateEnd +
				'&timeSize=' +
				timeSize,
			{
				method: 'GET',
			}
		);
	},

	getDiskUsage: () => {
		return request('/disk/usage', {
			method: 'GET',
		});
	},

	getDiskIO: (dateStart: string, dateEnd: string, timeSize: string) => {
		return request(
			'/disk/io?dateStart=' +
				dateStart +
				'&dateEnd=' +
				dateEnd +
				'&timeSize=' +
				timeSize,
			{
				method: 'GET',
			}
		);
	},

	getNetworkLoad: (dateStart: string, dateEnd: string, timeSize: string) => {
		return request(
			'/network/load?dateStart=' +
				dateStart +
				'&dateEnd=' +
				dateEnd +
				'&timeSize=' +
				timeSize,
			{
				method: 'GET',
			}
		);
	},

	getProcessInfo: () => {
		return request('/process', {
			method: 'GET',
		});
	},
};
