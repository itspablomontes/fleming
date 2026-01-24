/**
 * Get Graph Data API
 *
 * Fetches all events and edges for the authenticated patient's timeline graph.
 */

import type { GraphData } from "../types";

const API_BASE = "/api/timeline";

/**
 * Fetch graph data from the API.
 */
export async function getGraphData(): Promise<GraphData> {
	const response = await fetch(`${API_BASE}/graph`, {
		method: "GET",
		credentials: "include",
		headers: {
			"Content-Type": "application/json",
		},
	});

	if (!response.ok) {
		throw new Error(`Failed to fetch graph data: ${response.statusText}`);
	}

	return response.json();
}
