/**
 * Link Events API
 *
 * Creates a relationship between two timeline events.
 */

import type { LinkResponse, RelationshipType } from "../types";

const API_BASE = "/api/timeline";

interface LinkEventsParams {
	fromEventId: string;
	toEventId: string;
	relationshipType: RelationshipType;
}

/**
 * Create a link (edge) between two events.
 */
export async function linkEvents(
	params: LinkEventsParams,
): Promise<LinkResponse> {
	const response = await fetch(
		`${API_BASE}/events/${params.fromEventId}/link`,
		{
			method: "POST",
			credentials: "include",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				toEventId: params.toEventId,
				relationshipType: params.relationshipType,
			}),
		},
	);

	if (!response.ok) {
		throw new Error(`Failed to link events: ${response.statusText}`);
	}

	return response.json();
}

/**
 * Remove a link between events.
 */
export async function unlinkEvents(
	edgeId: string,
): Promise<{ success: boolean }> {
	const response = await fetch(`${API_BASE}/edges/${edgeId}`, {
		method: "DELETE",
		credentials: "include",
	});

	if (!response.ok) {
		throw new Error(`Failed to unlink events: ${response.statusText}`);
	}

	return response.json();
}
