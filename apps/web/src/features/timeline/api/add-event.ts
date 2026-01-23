import { apiClient } from "../../../lib/api-client";
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
