import { apiClient } from "@/lib/api-client";

export const deleteEvent = (id: string): Promise<{ success: boolean }> => {
    return apiClient(`/api/timeline/events/${id}`, {
        method: "DELETE",
    });
};
