import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const API_URL = process.env.API_URL || 'http://localhost:3001';

export interface ServerType {
	id: string;
	name: string;
	description: string;
	versions: string[];
}

export const load: PageServerLoad = async ({ locals }) => {
	// TODO: Fetch from API when endpoint is available
	// For now, hardcode the server types as specified
	const serverTypes: ServerType[] = [
		{
			id: 'vanilla',
			name: 'Vanilla',
			description: 'Official Minecraft server',
			versions: ['1.21.4', '1.21.3', '1.21.1', '1.20.6', '1.20.4']
		},
		{
			id: 'paper',
			name: 'Paper',
			description: 'High-performance Spigot fork with optimizations',
			versions: ['1.21.4', '1.21.3', '1.21.1', '1.20.6', '1.20.4']
		},
		{
			id: 'fabric',
			name: 'Fabric',
			description: 'Lightweight mod loader for Minecraft',
			versions: ['1.21.4', '1.21.3', '1.21.1', '1.20.6', '1.20.4']
		},
		{
			id: 'folia',
			name: 'Folia',
			description: 'Paper fork for large servers with regionized threading',
			versions: ['1.21.4', '1.21.3', '1.21.1']
		},
		{
			id: 'forge-installer',
			name: 'Forge Installer',
			description: 'Minecraft Forge mod loader (requires installation)',
			versions: ['1.21.4', '1.21.3', '1.20.6', '1.20.4', '1.20.1']
		},
		{
			id: 'neoforge-installer',
			name: 'NeoForge Installer',
			description: 'Modern Forge fork for modded Minecraft',
			versions: ['1.21.4', '1.21.3', '1.21.1', '1.20.6']
		},
		{
			id: 'purpur',
			name: 'Purpur',
			description: 'Paper fork with additional gameplay features',
			versions: ['1.21.4', '1.21.3', '1.21.1', '1.20.6', '1.20.4']
		}
	];

	return {
		user: locals.user,
		serverTypes
	};
};

export const actions: Actions = {
	create: async ({ request, cookies }) => {
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		const sessionCookie = cookies.get('session_id');
		const data = await request.formData();

		const name = data.get('name') as string;
		const serverType = data.get('serverType') as string;
		const version = data.get('version') as string;
		const port = data.get('port') as string;
		const minMemory = data.get('minMemory') as string;
		const maxMemory = data.get('maxMemory') as string;

		// Validation
		if (!name || name.trim().length === 0) {
			return fail(400, {
				error: 'Server name is required',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (name.length > 100) {
			return fail(400, {
				error: 'Server name must be 100 characters or less',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (!serverType) {
			return fail(400, {
				error: 'Server type is required',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (!version) {
			return fail(400, {
				error: 'Server version is required',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		// Validate port if provided
		if (port) {
			const portNum = parseInt(port, 10);
			if (isNaN(portNum) || portNum < 1024 || portNum > 65535) {
				return fail(400, {
					error: 'Port must be between 1024 and 65535',
					values: { name, serverType, version, port, minMemory, maxMemory }
				});
			}
		}

		// Validate memory
		const minMemNum = parseInt(minMemory || '2', 10);
		const maxMemNum = parseInt(maxMemory || '4', 10);

		if (isNaN(minMemNum) || minMemNum < 1) {
			return fail(400, {
				error: 'Minimum memory must be at least 1 GB',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (isNaN(maxMemNum) || maxMemNum < minMemNum) {
			return fail(400, {
				error: 'Maximum memory must be greater than or equal to minimum memory',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		// TODO: Call backend API to create server when endpoint is available
		// For now, just simulate success and redirect
		console.log('Creating server:', {
			name: name.trim(),
			serverType,
			version,
			port: port ? parseInt(port, 10) : null,
			minMemory: minMemNum,
			maxMemory: maxMemNum
		});

		try {
			// TODO: Replace with actual API call
			// const response = await fetch(`${API_URL}/api/v1/servers`, {
			//     method: 'POST',
			//     headers: {
			//         'Content-Type': 'application/json',
			//         'Cookie': `session_id=${sessionCookie}`
			//     },
			//     body: JSON.stringify({
			//         name: name.trim(),
			//         type: serverType,
			//         version,
			//         port: port ? parseInt(port, 10) : null,
			//         minMemory: minMemNum,
			//         maxMemory: maxMemNum
			//     })
			// });
			//
			// if (!response.ok) {
			//     const errorData = await response.json().catch(() => ({ error: 'Failed to create server' }));
			//     return fail(response.status, {
			//         error: errorData.error || 'Failed to create server',
			//         values: { name, serverType, version, port, minMemory, maxMemory }
			//     });
			// }

			// Successful creation - redirect to dashboard
			throw redirect(302, '/dashboard');
		} catch (error) {
			if (isRedirect(error)) {
				throw error;
			}

			console.error('Server creation error:', error);
			return fail(500, {
				error: 'An error occurred while creating the server',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}
	}
};
