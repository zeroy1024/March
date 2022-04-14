export default {
	formatDatetime: (timestamp: number) => {
		const date = new Date(timestamp);
		return (
			date.getFullYear() +
			'-' +
			(date.getMonth() + 1) +
			'-' +
			date.getDate() +
			' ' +
			date.getHours() +
			':' +
			date.getMinutes() +
			':' +
			date.getSeconds()
		);
	},
	formatBytes: (bytes: number) => {
		if (bytes > 1024 && bytes <= 1024 * 1024) {
			return (bytes / 1024).toFixed(2) + 'KB';
		} else if (bytes > 1024 * 1024 && bytes < 1024 * 1024 * 1024) {
			return (bytes / 1024 / 1024).toFixed(2) + 'MB';
		} else if (
			bytes > 1024 * 1024 * 1024 &&
			bytes < 1024 * 1024 * 1024 * 1024
		) {
			return (bytes / 1024 / 1024 / 1024).toFixed(2) + 'GB';
		} else {
			return bytes + 'B';
		}
	},
	hexToRgba: (hex: string, opacity: number) => {
		return (
			'rgba(' +
			parseInt('0x' + hex.slice(1, 3)) +
			',' +
			parseInt('0x' + hex.slice(3, 5)) +
			',' +
			parseInt('0x' + hex.slice(5, 7)) +
			',' +
			opacity +
			')'
		);
	},
};
