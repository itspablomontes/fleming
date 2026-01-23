import { apiClient } from "../../../lib/api-client";
import type { TimelineEvent } from "../types";

export interface TimelineResponse {
    events: TimelineEvent[];
}

export const getTimeline = (): Promise<TimelineResponse> => {
    return apiClient("/api/timeline");
};
