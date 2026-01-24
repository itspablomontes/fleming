import { apiClient } from "@/lib/api-client";
import type { UploadResponse } from "../types";

export interface UploadPayload {
	file: File;
	eventType: string;
	description: string;
	provider?: string;
	date: string;
}

export const addEvent = (payload: UploadPayload): Promise<UploadResponse> => {
	const formData = new FormData();
	formData.append("file", payload.file);
	formData.append("eventType", payload.eventType);
	formData.append("description", payload.description);
	if (payload.provider) {
		formData.append("provider", payload.provider);
	}
	formData.append("date", payload.date);

	return apiClient("/api/timeline/events", {
		method: "POST",
		body: formData,
	});
};

export interface CorrectionPayload extends Omit<UploadPayload, "file"> {
	id: string;
	file?: File;
}

export const correctEvent = (payload: CorrectionPayload): Promise<UploadResponse> => {
	// For now, correction is a JSON POST to /events/:id/correction
	// If it needs a file, we might need multipart, but the handler suggests JSON
	return apiClient(`/api/timeline/events/${payload.id}/correction`, {
		method: "POST",
		body: JSON.stringify({
			eventType: payload.eventType,
			description: payload.description,
			provider: payload.provider,
			date: payload.date,
		}),
	});
};
