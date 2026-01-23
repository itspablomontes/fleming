/**
 * Get Graph Data API
 *
 * Fetches all events and edges for the authenticated patient's timeline graph.
 */

// =============================================================================
// TODO: REMOVE MOCK DATA - Delete this import when connecting to real API
// =============================================================================
import { getMockGraphData } from "../mocks/graph-data";
import type { GraphData } from "../types";

const API_BASE = "/api/timeline";

// =============================================================================
// TODO: REMOVE MOCK DATA - Set to false to use real API
// When ready: change `USE_MOCK_DATA = false` or delete the mock logic entirely
// =============================================================================
const USE_MOCK_DATA = true;

/**
 * Fetch graph data from the API.
 * Currently returns mock data for development.
 */
export async function getGraphData(): Promise<GraphData> {
	// TODO: REMOVE MOCK DATA - Delete this block when connecting to real API
	if (USE_MOCK_DATA) {
		return getMockGraphData();
	}

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
